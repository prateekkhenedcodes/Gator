package main

import (
	"fmt"
	"github.com/prateekkhenedcodes/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	err = cfg.SetUser("prateek")
	if err != nil {
		fmt.Println("Error setting the user:", err)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config after setting user:", err)
	}

	fmt.Println(cfg.CurrentUserName)
	fmt.Println(cfg.DBUrl)
}
