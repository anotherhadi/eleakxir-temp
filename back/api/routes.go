package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
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

type sseMessage struct {
	event string
	data  gin.H
}

func (api *API) SetupRoutes() {
	api.Router.GET("/dataleaks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"dataleaks": api.Dataleaks})
	})

	api.Router.GET("/search", func(c *gin.Context) {
		// Prepare headers for SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Context to cancel if client disconnects
		ctx := c.Request.Context()

		// Channel for safely sending events to the client
		sseChan := make(chan sseMessage)
		defer close(sseChan)

		// Writer goroutine
		go func() {
			for {
				select {
				case msg, ok := <-sseChan:
					if !ok {
						// Channel closed, goroutine should exit
						return
					}
					// Check if client has disconnected before attempting to write
					select {
					case <-ctx.Done():
						// Client disconnected, do not write
						log.Printf("Client disconnected, abandoning SSE write for event %s", msg.event)
						return
					default:
						// Client still connected, proceed to write
						sendSSEEvent(c.Writer, msg.event, msg.data)
						flusher.Flush()
					}
				case <-ctx.Done():
					// Context cancelled (client disconnected), goroutine should exit
					return
				}
			}
		}()

		query := leak.ParseQuery(sanitizeQuery(c.Query("q")))
		if len(query.Terms) == 0 {
			sseChan <- sseMessage{"error", gin.H{"message": "Invalid query: no search terms provided."}}
			return
		}

		columns := strings.Split(c.Query("columns"), ",")
		if len(columns) == 0 || (len(columns) == 1 && columns[0] == "") {
			sseChan <- sseMessage{"error", gin.H{"message": "Invalid query: no columns provided."}}
			return
		}

		if api.Dataleaks.TotalDataleaks == 0 {
			sseChan <- sseMessage{"error", gin.H{"message": "No parquet files configured."}}
			return
		}

		sseChan <- sseMessage{"start", gin.H{"percentage": 0}}

		// Launch concurrent workers for each dataleak
		done := make(chan struct{})

		go func() {
			defer close(done)
			for i, file := range api.Dataleaks.Dataleaks {
				progress := float64(i) / float64(api.Dataleaks.TotalDataleaks) * 100
				sseChan <- sseMessage{"progress", gin.H{"percentage": progress}}

				results, err := api.Dataleaks.Search(file.Path, columns, query)
				if err != nil {
					sseChan <- sseMessage{"file_error", gin.H{
						"file_path": file.Path,
						"message":   fmt.Sprintf("Error processing file: %s", err.Error()),
					}}
					continue
				}

				if len(results) > 0 {
					sseChan <- sseMessage{"new_results", gin.H{"results": results}}
				}
			}
		}()

		<-done
		sseChan <- sseMessage{"complete", gin.H{}}
	})
}

func sendSSEEvent(w http.ResponseWriter, eventType string, data gin.H) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling SSE data for event %s: %v", eventType, err)
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, jsonData)
}
