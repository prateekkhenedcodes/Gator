package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/prateekkhenedcodes/Gator/internal/config"
	"github.com/prateekkhenedcodes/Gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	// Initialize database BEFORE creating the state
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("Failed to open or connect to database: %v\n", err)
		os.Exit(1)
	}

	// Create queries object
	dbQueries := database.New(db)

	// Now create state with both config and db
	s := &config.State{
		ConfigPtr: &cfg,
		Db:        dbQueries, // Make sure 'db' is the correct field name in your State struct
	}

	cmds := config.Commands{
		CmdHandlers: make(map[string]func(*config.State, config.Command) error),
	}
	cmds.Register("login", config.HandlerLogin)
	cmds.Register("register", config.HandlerRegister)
	cmds.Register("reset", config.HandlerReset)
	cmds.Register("users", config.HandlerGetUsers)
	cmds.Register("agg", config.Handleragg)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments: command name required")
	}
	c := config.Command{
		Name: args[1],
		Args: args[2:],
	}
	err = cmds.Run(s, c)
	if err != nil {
		log.Fatal(err)
	}
}
