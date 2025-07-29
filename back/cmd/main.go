package main

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/anotherhadi/eleakxir/api"
	"github.com/anotherhadi/eleakxir/leak"
)

func main() {
	if os.Getenv("LEAK_DIRECTORY") == "" {
		panic("LEAK_DIRECTORY environment variable is not set")
	}
	leakDir := filepath.Clean(os.Getenv("LEAK_DIRECTORY"))
	cacheDir := filepath.Join(leakDir, "cache")
	if os.Getenv("CACHE_DIRECTORY") != "" {
		cacheDir = filepath.Clean(os.Getenv("CACHE_DIRECTORY"))
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	maxConcurrentSearchesStr := os.Getenv("MAX_CONCURRENT_SEARCHES")
	maxConcurrentSearches := 8
	if maxConcurrentSearchesStr != "" {
		var err error
		maxConcurrentSearches, err = strconv.Atoi(maxConcurrentSearchesStr)
		if err != nil {
			panic("Invalid MAX_CONCURRENT_SEARCHES environment variable: " + err.Error())
		}
	}
	devStr := os.Getenv("DEV")
	dev := false
	if devStr == "TRUE" {
		dev = true
	}

	d, err := leak.OpenDataleaks(leakDir, cacheDir)
	if err != nil {
		panic(err)
	}
	defer d.CloseDataleaks()

	api := api.NewAPI(d, maxConcurrentSearches, dev)
	err = api.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
