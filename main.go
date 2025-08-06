package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/quduss/gator/internal/config"
	"github.com/quduss/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: gator login <username>")
	}
	username := cmd.args[0]
	// Check if user exists in database
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user '%s' not found", username)
		}
		return fmt.Errorf("database error: %w", err)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: gator register <username>")
	}
	username := cmd.args[0]
	// Check if user already exists
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user '%s' already exists", username)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("database error: %w", err)
	}
	// Create new user
	now := time.Now()
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	// Set the current user in the config
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User '%s' created successfully!\n", username)
	fmt.Printf("User data: ID=%s, Name=%s, CreatedAt=%s\n",
		user.ID, user.Name, user.CreatedAt.Format(time.RFC3339))

	return nil
}

// handlerReset handles the reset command
func handlerReset(s *state, cmd command) error {
	// Delete all users from the database
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset database: %w", err)
	}

	fmt.Println("Database has been reset successfully!")
	return nil
}

// handlerUsers handles the users command
func handlerUsers(s *state, cmd command) error {
	// Get all users from the database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	// Get current user from config
	currentUser := s.cfg.CurrentUserName

	// Print all users
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

// handlerAgg handles the agg command
func handlerAgg(s *state, cmd command) error {
	// Fetch the RSS feed from wagslane.dev
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	// Print the entire feed structure
	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Link: %s\n", feed.Channel.Link)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	fmt.Printf("\nItems:\n")

	for i, item := range feed.Channel.Item {
		fmt.Printf("\nItem %d:\n", i+1)
		fmt.Printf("  Title: %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Description: %s\n", item.Description)
		fmt.Printf("  PubDate: %s\n", item.PubDate)
	}

	return nil
}

// handlerAddFeed handles the addfeed command
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: gator addfeed <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	// Get the current user from the config
	currentUserName := s.cfg.CurrentUserName
	if currentUserName == "" {
		return fmt.Errorf("no user is currently logged in")
	}

	// Get the current user from the database
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get current user: %w", err)
	}

	// Create the feed
	now := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	// Print the feed details
	fmt.Printf("Feed created successfully!\n")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("User ID: %s\n", feed.UserID)
	fmt.Printf("Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))

	return nil
}

// handlerFeeds handles the feeds command
func handlerFeeds(s *state, cmd command) error {
	// Get all feeds from the database
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	// Print all feeds
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feed(s):\n\n", len(feeds))

	for i, feed := range feeds {
		fmt.Printf("Feed %d:\n", i+1)
		fmt.Printf("  Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  User: %s\n", feed.UserName)
		fmt.Println()
	}

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}
	// Open database connection
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		os.Exit(1)
	}
	cmdName := args[1]
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:]
	}
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = cmds.run(programState, cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
