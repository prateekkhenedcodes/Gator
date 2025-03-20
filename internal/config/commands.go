package config

import "fmt"

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login requires username argument")
	}
	err := s.ConfigPtr.SetUser(cmd.Args[0])
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
