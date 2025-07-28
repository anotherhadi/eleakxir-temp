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
		if len(columns) == 0 {
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "Invalid query: no columns provided."})
			return
		}

		if api.Dataleaks.TotalDataleaks == 0 {
			sendSSEEvent(c.Writer, flusher, "error", gin.H{"message": "No parquet files configured."})
			return
		}

		type sseMsg struct {
			event string
			data  gin.H
		}

		sseChan := make(chan sseMsg, 100)
		var wg sync.WaitGroup

		// Writer goroutine
		go func() {
			for msg := range sseChan {
				sendSSEEvent(c.Writer, flusher, msg.event, msg.data)
			}
		}()

		sseChan <- sseMsg{"start", gin.H{"percentage": 0}}

		for i, file := range api.Dataleaks.Dataleaks {
			wg.Add(1)

			go func(i int, file leak.Dataleak) {
				defer wg.Done()

				progress := float64(i) / float64(api.Dataleaks.TotalDataleaks) * 100
				sseChan <- sseMsg{"progress", gin.H{"percentage": progress}}

				results, err := api.Dataleaks.Search(file.Path, columns, query)
				if err != nil {
					sseChan <- sseMsg{"file_error", gin.H{
						"file_path": file.Path,
						"message":   fmt.Sprintf("Error processing file: %s", err.Error()),
					}}
					return
				}

				if len(results) > 0 {
					sseChan <- sseMsg{"new_results", gin.H{"results": results}}
				}
			}(i, file)
		}

		// Fin de traitement
		go func() {
			wg.Wait()
			sseChan <- sseMsg{"complete", gin.H{}}
			close(sseChan)
		}()
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
