package leak

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

const CACHE_FILE = "dataleaks_cache.json"

// TODO: check os.FileInfo.ModTime() to see if the file has changed since last cache update

func (d *Dataleaks) getCache() error {
	if err := createDirectoryIfNotExists(d.CacheDirectory); err != nil {
		return err
	}

	cachePath := filepath.Join(d.CacheDirectory, CACHE_FILE)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		log.Warn("Failed to read dataleaks cache file", "error", err)
	} else if err := json.Unmarshal(data, &d.Dataleaks); err != nil {
		return errors.New("failed to unmarshal dataleaks cache file")
	}

	existing := make(map[string]bool)
	for _, d := range d.Dataleaks {
		existing[d.Path] = true
	}

	paths, err := listParquetPaths(d.DataleaksDirectory)
	if err != nil {
		return err
	}
	log.Infof("Found a total of %d dataleaks", len(paths))

	var (
		mu       sync.Mutex
		newLeaks []Dataleak
		wg       sync.WaitGroup
	)

	for _, path := range paths {
		if existing[path] {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			path := filepath.Join(d.DataleaksDirectory, path)
			columns, err := d.GetParquetColumns(path)
			if err != nil {
				log.Error("Failed to read CSV header", "error", err, "path", path)
				return
			}

			length, err := d.GetParquetLength(path)
			if err != nil {
				log.Error("Failed to get row count", "error", err, "path", path)
				return
			}
			size := getFileSize(path)

			leak := Dataleak{
				Path:    strings.TrimPrefix(strings.TrimPrefix(path, d.DataleaksDirectory), "/"),
				Columns: columns,
				Length:  length,
				Size:    size,
			}
			leak.Name = formatParquetName(leak.Path)

			log.Infof("Adding new dataleak: %s", path)

			mu.Lock()
			d.Dataleaks = append(d.Dataleaks, leak)
			newLeaks = append(newLeaks, leak)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if len(newLeaks) > 0 {
		data, err := json.MarshalIndent(d.Dataleaks, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling cache: %w", err)
		}
		if err := os.WriteFile(cachePath, data, 0644); err != nil {
			return fmt.Errorf("error writing cache: %w", err)
		}
		log.Info("Dataleaks cache updated successfully", "count", len(d.Dataleaks))
	}

	// STATS
	d.TotalDataleaks = uint64(len(d.Dataleaks))
	d.TotalRows = 0
	d.TotalSize = 0
	for _, leak := range d.Dataleaks {
		d.TotalRows += leak.Length
		d.TotalSize += leak.Size
	}

	return nil
}

func formatParquetName(path string) string {
	str := strings.TrimSuffix(path, ".parquet")
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, "-", " ")
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.ReplaceAll(str, "/", " ")

	// capitalize first letter of each word
	words := strings.Fields(str)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	str = strings.Join(words, " ")

	return str
}
