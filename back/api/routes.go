package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/anotherhadi/eleakxir/leak"
	"github.com/gin-gonic/gin"
)

type SearchResult struct {
	Results  interface{}
	FilePath string
	Error    error
	Done     bool // <-- Cette propriété n'est pas utilisée, on peut la retirer.
}

func sanitizeQuery(query string) string {
	query = strings.TrimSpace(query)
	query = strings.ReplaceAll(query, "“", `"`)
	query = strings.ReplaceAll(query, "”", `"`)
	query = strings.ReplaceAll(query, "‘", `'`)
	query = strings.ReplaceAll(query, "’", `'`)
	query = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, query)
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
	return query
}

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
		if len(columns) == 0 || (len(columns) == 1 && columns[0] == "") {
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "Invalid query: no columns provided."})
			return
		}

		if api.Dataleaks.TotalDataleaks == 0 {
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "No parquet files configured."})
			return
		}

		sendSSEEvent(c.Writer, flusher, "start", gin.H{"percentage": 0})

		var wg sync.WaitGroup
		resultsChan := make(chan SearchResult) // Canal non bufferisé, c'est bien
		progressChan := make(chan float64)     // Canal non bufferisé, c'est bien

		// Lancer les goroutines de recherche
		for _, file := range api.Dataleaks.Dataleaks {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				res, err := api.Dataleaks.Search(filePath, columns, query)
				if err != nil {
					resultsChan <- SearchResult{Error: err, FilePath: filePath}
					return
				}
				resultsChan <- SearchResult{Results: res, FilePath: filePath}
			}(file.Path)
		}

		// Goroutine qui attend que toutes les recherches soient faites, PUIS ferme les canaux.
		// C'est la clé de la correction.
		go func() {
			wg.Wait() // Attendre que TOUTES les goroutines de recherche aient appelé Done()
			close(resultsChan)
			close(progressChan) // Fermer le canal de progression ici aussi
		}()

		// Goroutine pour gérer l'envoi de la progression globale
		// Cette goroutine lit depuis progressChan et envoie les événements SSE
		go func() {
			for p := range progressChan {
				sendSSEEvent(c.Writer, flusher, "progress", gin.H{
					"percentage": p,
				})
			}
		}()

		// Boucle principale pour recevoir les résultats et les erreurs des goroutines
		processedFilesCount := 0
		for searchRes := range resultsChan { // Cette boucle se terminera quand resultsChan sera fermé
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
			// Assurez-vous que progressChan n'est pas fermé avant cet envoi.
			// La goroutine qui ferme les canaux attend `wg.Wait()`, donc cela devrait être sûr.
			progressChan <- float64(processedFilesCount) / float64(api.Dataleaks.TotalDataleaks) * 100
		}

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
