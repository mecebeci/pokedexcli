package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type LocationAreaResponse struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []struct{
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	Pokemon []struct{
		Pokemon struct {
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	ID int32 `json:"id"`
	Name string `json:"name"`
	BaseExperience int32 `json:"base_experience"`
	Height int32 `json:"height"`
	Weight int32 `json:"weight"`
	Stats []struct{
		BaseStat int32 `json:"base_stat"`
		Stat struct{
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct{
		Slot int `json:"slot"`
		Type struct{
			Name string `json:"name"`
			URL string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func commandHelp(c *config, args []string) error {
	fmt.Println("Available commands:")
	for _, cmd := range table{
		fmt.Printf("	%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}


func commandExit(c *config, args []string) error{
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(c *config, args []string) error {
	return fetchAndPrintLocations(c, c.next)
}

func commandMapBack(c *config, args []string) error {
	return fetchAndPrintLocations(c, c.previous)
}


func commandExplore(c *config, args []string) error {
	if len(args) < 1{
		fmt.Println("You must provide a location to explore. Example: explore kanto-route-1")
		return nil
	}

	location := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", location)

	var body []byte
	if data, ok := c.cache.Get(url); ok {
		body = data
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		c.cache.Add(url, body)
	}
	var loc Location
	if err := json.Unmarshal(body, &loc); err != nil{
		return err
	}

    fmt.Printf("Exploring %s...\n", loc.Name)
    fmt.Println("Found Pokémon:")
    for _, p := range loc.Pokemon {
        fmt.Printf("- %s\n", p.Pokemon.Name)
    }
    return nil
}

func commandCatch(c *config, args []string) error{
	if len(args) < 1 {
		fmt.Println("Please provide a Pokemon name")
		return nil
	}

	name := args[0]

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var p Pokemon
	err = json.Unmarshal(body, &p)
	if err != nil {
		return err
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	chance := r.Intn(int(p.BaseExperience) + 50)
	fmt.Printf("Throwing a Pokeball at %s...\n", p.Name)
	if chance < 50 {
		fmt.Printf("%s was caught!\n", p.Name)
		c.pokedex[p.Name] = p
	} else {
		fmt.Printf("%s escaped!\n", p.Name)
	}
	return nil
}

func commandInspect(c *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please specify a Pokemon name")
		return nil
	}

	name := args[0]
	p, ok := c.pokedex[name]
	if !ok {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *config, args []string) error {
	if len(c.pokedex) == 0 {
		return fmt.Errorf("you have not caught any pokemon")
	}
	for _, p := range c.pokedex {
		fmt.Printf("- %s\n", p.Name)
	}
	return nil
}

func fetchAndPrintLocations(c *config, url string) error {
	if url == "" {
		fmt.Println("No more locations")
		return nil
	}

	var body []byte
	var err error

	if data, ok := c.cache.Get(url); ok {
		fmt.Println("Using cached data")
		body = data
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		c.cache.Add(url, body)
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	c.next = data.Next
	c.previous = data.Previous

	return nil
}