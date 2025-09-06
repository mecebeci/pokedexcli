package main

import (
	"time"

	"github.com/mecebeci/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	next     string
	previous string
	cache *pokecache.Cache
	pokedex map[string]Pokemon
}

var table = map[string]cliCommand{}

func main() {
	cfg := &config{
		next: "https://pokeapi.co/api/v2/location-area/",
		cache: pokecache.NewCache(5 * time.Second),
		pokedex: make(map[string]Pokemon),
	}

	table["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: commandHelp,
	}

	table["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	table["map"] = cliCommand{
		name:        "map",
		description: "Displays the names of 20 location areas in the Pokemon world",
		callback:    commandMap,
	}

    table["mapb"] = cliCommand{
        name: "mapb",
        description: "Displays the previous 20 location areas",
        callback: commandMapBack,
    }

	table["explore"] = cliCommand{
		name: "explore",
		description: "Explore a specific location area",
		callback: commandExplore,
	}

	table["catch"] = cliCommand{
		name: "catch",
		description: "Catch a specific pokemon",
		callback: commandCatch,
	}

	table["inspect"] = cliCommand{
		name: "inspect",
		description: "Inspect a specific pokemon",
		callback: commandInspect,
	}

	table["pokedex"] = cliCommand{
		name: "pokedex",
		description: "Show all catched pokemons",
		callback: commandPokedex,
	}

	StartREPL(cfg)
}
