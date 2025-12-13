package main

import (
	"log"

	"github.com/Parth-11/Codescribe/cmd"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (IGNORED).")
	}

	cmd.Execute()
}
