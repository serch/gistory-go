package main

import (
	"fmt"
	"os"

	"github.com/serch/gistory-go/gistory"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gistory-go <gemName>")
		os.Exit(1)
	}
	gemName := os.Args[1]
	gistory.Run(gemName, ".")
}
