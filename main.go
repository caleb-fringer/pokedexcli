package main

import (
	"os"
)

func main() {
	ok := doREPL()
	if !ok {
		os.Exit(1)
	}
}
