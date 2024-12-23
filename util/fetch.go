package main

import (
	"fmt"

	"github.com/tschroed/isodisplay"
)

func main() {
	f := isodisplay.NewEmissionsFetcher()
	data, err := f.RawData()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
