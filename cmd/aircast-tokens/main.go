package main

import (
	"os"

	tokens "github.com/aircast-one/aircast-tokens"
)

func main() {
	if _, err := os.Stdout.Write(tokens.CSS); err != nil {
		os.Exit(1)
	}
}
