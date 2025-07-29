package main

import (
	"fmt"
	"log"

	"gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	err = cfg.SetUser("qudus")
	if err != nil {
		log.Fatalf("Error setting user: %v", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config after update: %v", err)
	}
	fmt.Printf("Config contents:\n")
	fmt.Printf("Database URL: %s\n", cfg.DbURL)
	fmt.Printf("Current User: %s\n", cfg.CurrentUserName)
}
