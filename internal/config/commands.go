package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/Gator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login requires username argument")
	}

	username := cmd.Args[0]

	// Check if the user exists in the database
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user '%s' does not exist", username)
	}

	// If we get here, the user exists, so proceed with login
	s.ConfigPtr.CurrentUserName = username

	err = s.ConfigPtr.Save()
	if err != nil {
		return err
	}

	fmt.Println("user has been set")
	return nil
}

type Commands struct {
	CmdHandlers map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CmdHandlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.CmdHandlers[cmd.Name]
	if !exists {
		return fmt.Errorf("handler not found")
	}
	return handler(s, cmd)
}

func HandlerRegister(s *State, cmd Command) error {
	// Ensure name was provided
	if len(cmd.Args) == 0 {
		fmt.Println("Please provide a username")
		os.Exit(1)
	}

	name := cmd.Args[0]

	// Try to get the user by name first to check if they exist
	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		// User already exists
		fmt.Printf("User with name %s already exists\n", name)
		os.Exit(1)
	}

	// Create new user
	user, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		},
	)
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		os.Exit(1)
	}

	// Set current user in config
	s.ConfigPtr.CurrentUserName = name
	if err := s.ConfigPtr.Save(); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		os.Exit(1)
	}

	// Print success message and user data
	fmt.Printf("User created: %s\n", name)
	fmt.Printf("User data: %+v\n", user)

	return nil
}
func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete the users", err)
	}
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.ConfigPtr.CurrentUserName {
			fmt.Println("* ", user.Name, "(current)")
			continue
		}
		fmt.Println("* ", user.Name)
	}
	return nil
}
