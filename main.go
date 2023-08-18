package main

import (
	"log"
	"os"

	_ "go-in-action/matchers"
	"go-in-action/search"
)

func init() {
	// All init functions in any code file that are part of the program will get called before the main function.
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("president")
}
