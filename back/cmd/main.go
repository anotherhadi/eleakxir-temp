package main

import (
	"os"
	"path/filepath"

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

	api := api.NewAPI(d, dev)
	err = api.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
