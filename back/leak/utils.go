package leak

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func ParseQuery(jsonstr string) Query {
	var query Query
	if err := json.Unmarshal([]byte(jsonstr), &query); err != nil {
		log.Error("Failed to parse query JSON", "error", err)
		return Query{}
	}
	return query
}

func listParquetPaths(baseDir string) ([]string, error) {
	var paths []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".parquet") {
			return err
		}
		paths = append(paths, strings.TrimPrefix(strings.TrimPrefix(path, baseDir), "/"))
		return nil
	})
	return paths, err
}

func getFileSize(path string) uint64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return uint64(info.Size() / (1024 * 1024)) // MB
}

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		log.Infof("Created directory: %s", path)
	}
	return nil
}
