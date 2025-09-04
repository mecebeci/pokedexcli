package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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


func commandHelp(c *config) error {
	fmt.Println("Available commands:")
	for _, cmd := range table{
		fmt.Printf("	%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}


func commandExit(c *config) error{
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(c *config) error{
	if c.next == ""{
		fmt.Println("No more locations")
		return nil
	}

	res, err := http.Get(c.next)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil{
		return err
	}

	for _, area := range data.Results{
		fmt.Println(area.Name)
	}

	c.next = data.Next
	c.previous = data.Previous

	return nil
}

func commandMapBack(c *config) error {
	if c.previous == ""{
		fmt.Println("No more locations")
		return nil
	}

	res, err := http.Get(c.previous)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	for _, area := range data.Results{
		fmt.Println(area.Name)
	}

	c.next = data.Next
	c.previous = data.Previous
	return nil
}