package main

import (
	"fmt"
	"log"
	"os"

	"github.com/prateekkhenedcodes/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	s := &config.State{
		ConfigPtr: &cfg,
	}
	cmds := config.Commands{
		CmdHandlers: make(map[string]func(*config.State, config.Command) error),
	}
	cmds.Register("login", config.HandlerLogin)

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
