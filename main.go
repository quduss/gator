package main

import (
	"fmt"
	"log"

	"gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: gator login <username>")
	}
	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User has been set to: %s\n", username)
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
