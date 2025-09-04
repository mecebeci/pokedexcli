package main

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
}

var table = map[string]cliCommand{}

func main() {
	cfg := &config{
		next: "https://pokeapi.co/api/v2/location-area/",
	}

	table["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(c *config) error {
			return commandHelp(c)
		},
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

	StartREPL(cfg)
}
