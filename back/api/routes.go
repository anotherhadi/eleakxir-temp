package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync" // <-- Importez le package sync pour WaitGroup
	"unicode"

	"github.com/anotherhadi/eleakxir/leak"
	"github.com/gin-gonic/gin"
)

// Définition de la structure pour les résultats et les erreurs des goroutines
type SearchResult struct {
	Results  interface{} // Utilisez interface{} car gin.H est map[string]interface{}
	FilePath string
	Error    error
	Done     bool // Indique si ce résultat marque la fin du traitement pour un fichier
}

func sanitizeQuery(query string) string {
	query = strings.TrimSpace(query)

	// Replace fancy quotes with standard ones
	query = strings.ReplaceAll(query, "“", `"`)
	query = strings.ReplaceAll(query, "”", `"`)
	query = strings.ReplaceAll(query, "‘", `'`)
	query = strings.ReplaceAll(query, "’", `'`)

	// Remove non-printable and control characters
	query = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, query)

	// Collapse multiple spaces
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	return query
}

// TODO: implement full-text search

func (api *API) SetupRoutes() {
	api.Router.GET("/dataleaks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"dataleaks": api.Dataleaks})
	})

	api.Router.GET("/search", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Flush()

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		query := leak.ParseQuery(sanitizeQuery(c.Query("q")))
		if len(query.Terms) == 0 {
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "Invalid query: no search terms provided."})
			return
		}

		columns := strings.Split(c.Query("columns"), ",")
		if len(columns) == 0 || (len(columns) == 1 && columns[0] == "") { // Handle empty string after split
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "Invalid query: no columns provided."})
			return
		}

		if api.Dataleaks.TotalDataleaks == 0 {
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "No parquet files configured."})
			return
		}

		sendSSEEvent(c.Writer, flusher, "start", gin.H{"percentage": 0})

		// --- Début des modifications pour la concurrence ---
		var wg sync.WaitGroup // Utilisé pour attendre que toutes les goroutines se terminent
		// Canal pour envoyer les résultats et les erreurs des goroutines au thread principal
		resultsChan := make(chan SearchResult)
		// Canal pour envoyer les mises à jour de progression (optionnel, mais utile)
		progressChan := make(chan float64)

		// Lancer les goroutines pour la recherche de fichiers
		for _, file := range api.Dataleaks.Dataleaks {
			wg.Add(1) // Incrémente le compteur de goroutines à attendre
			go func(filePath string) {
				defer wg.Done() // Décrémente le compteur quand la goroutine se termine
				// Effectuer la recherche
				res, err := api.Dataleaks.Search(filePath, columns, query)
				if err != nil {
					// Envoyer l'erreur via le canal
					resultsChan <- SearchResult{Error: err, FilePath: filePath}
					return
				}
				// Envoyer les résultats via le canal
				resultsChan <- SearchResult{Results: res, FilePath: filePath}
			}(file.Path) // Passer le chemin du fichier à la goroutine
		}

		// Goroutine pour gérer l'envoi de la progression globale
		// Cette goroutine est importante pour éviter les problèmes de concurrence
		// lors de la mise à jour de la variable `processedFilesCount`
		go func() {
			for p := range progressChan {
				sendSSEEvent(c.Writer, flusher, "progress", gin.H{
					"percentage": p,
				})
			}
		}()

		// Goroutine pour fermer les canaux une fois que toutes les goroutines de recherche sont terminées
		go func() {
			wg.Wait()           // Attendre que toutes les goroutines de recherche soient terminées
			close(resultsChan)  // Fermer le canal des résultats
			close(progressChan) // Fermer le canal de progression
		}()

		// Boucle principale pour recevoir les résultats et les erreurs des goroutines
		processedFilesCount := 0
		for searchRes := range resultsChan {
			if searchRes.Error != nil {
				sendSSEEvent(c.Writer, flusher, "file_error", gin.H{
					"file_path": searchRes.FilePath,
					"message":   fmt.Sprintf("Error processing file: %s", searchRes.Error.Error()),
				})
			} else {
				sendSSEEvent(c.Writer, flusher, "new_results", gin.H{
					"results": searchRes.Results,
				})
			}
			processedFilesCount++
			// Envoyer la progression via le canal de progression
			progressChan <- float64(processedFilesCount) / float64(api.Dataleaks.TotalDataleaks) * 100
		}
		// --- Fin des modifications pour la concurrence ---

		sendSSEEvent(c.Writer, flusher, "complete", gin.H{})
	})
}

func sendSSEEvent(w http.ResponseWriter, flusher http.Flusher, eventType string, data gin.H) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling SSE data for event %s: %v", eventType, err)
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, jsonData)
	flusher.Flush()
}
