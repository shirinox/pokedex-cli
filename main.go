package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	api "github.com/shirinox/pokeapi"
)

type CliCommand struct {
	name        string
	description string
	command     func(*api.Config, []string) error
}

func commandHelp(c *api.Config, args []string) error {
	if len(args) != 0 {
		fmt.Println("This command does now allow args.")
	}
	fmt.Printf(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`)
	return nil
}

func commandExit(c *api.Config, args []string) error {
	os.Exit(0)
	return nil
}

func main() {

	mapConfig := api.Config{}

	commands := map[string]CliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			command:     commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			command:     commandExit,
		},
		"map": {
			name:        "map",
			description: "Navigate 20 locations forward",
			command:     api.CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Navigate 20 locations back",
			command:     api.CommandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "See pokemon in a specific zone",
			command:     api.CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			command:     api.CommandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon",
			command:     api.CommandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show all pokemon",
			command:     api.CommandPokedex,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			if _, ok := commands[strings.Fields(scanner.Text())[0]]; ok {
				fullCommandFields := strings.Fields(scanner.Text())
				command := fullCommandFields[0]
				args := fullCommandFields[1:]
				commands[command].command(&mapConfig, args)
				break
			}

			fmt.Printf("Command '%v' does not exit.\n", scanner.Text())
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
}
