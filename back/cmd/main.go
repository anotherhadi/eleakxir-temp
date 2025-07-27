package main

import (
	"github.com/anotherhadi/eleakxir/api"
	"github.com/anotherhadi/eleakxir/leak"
)

func main() {
	d, err := leak.OpenDataleaks("/mnt/external/leaks/", "/mnt/external/leaks/cache")
	if err != nil {
		panic(err)
	}
	defer d.CloseDataleaks()

	// fmt.Println(d.TotalDataleaks)
	// fmt.Println(d.TotalRows)
	// fmt.Println(d.TotalSize)

	api := api.NewAPI(d)
	err = api.Run(":8080")
	if err != nil {
		panic(err)
	}
}
