package config

import (
	"context"
	"fmt"
	"log"
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
		return fmt.Errorf("could not delete the users, %s", err)
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

func Handleragg(s *State, cmd Command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Print(*feed)
	return nil
}

func HandleAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not enough arguments")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	cUser := s.ConfigPtr.CurrentUserName
	if cUser == "" {
		log.Fatalf("No current username is configured")
	}

	cId, err := s.Db.GetIdOfUser(context.Background(), cUser)
	if err != nil {
		return err
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    cId,
	})
	if err != nil {
		return err
	}

	// After successfully creating the feed, create a feed follow record
	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    cId,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed added successfully: %s\n", feedFollow.FeedName)
	return nil
}
func HandleFeed(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeed(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		userName, err := s.Db.GetUserOfIf(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(userName)
	}
	return nil
}

func HandleFollow(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Error: Missing URL argument. Usage: follow <url>")
		return fmt.Errorf("missing URL argument")
	}

	url := cmd.Args[0]
	curUser := s.ConfigPtr.CurrentUserName
	curUserId, err := s.Db.GetIdOfUser(context.Background(), curUser)
	if err != nil {
		fmt.Println("Error: No user is logged in. Please log in first.")
		return err
	}

	feedId, err := s.Db.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Println("Error: Feed not found for the given URL:", url)
		return err
	}

	newFeed, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    curUserId,
		FeedID:    feedId,
	})
	if err != nil {
		fmt.Println("Error creating feed follow:", err)
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("%s is now following '%s'\n", newFeed.UserName, newFeed.FeedName)
	return nil
}

func HandleFollowing(s *State, cmd Command) error {
	userName := s.ConfigPtr.CurrentUserName

	userId, err := s.Db.GetIdOfUser(context.Background(), userName)
	if err != nil {
		return err
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), userId)
	if err != nil {
		return err
	}

	fmt.Println("Feeds you are following:")
	for _, follow := range feedFollows {
		fmt.Println("- " + follow.FeedName)
	}

	return nil

}
