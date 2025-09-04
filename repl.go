package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartREPL(cfg *config) {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        if scanner.Scan() {
            command := scanner.Text()
            if len(command) > 0 {
                words := cleanInput(command)
                cmd, ok := table[words[0]]
                if ok {
                    args := words[1:]
                    err := cmd.callback(cfg, args)
                    if err != nil {
                        fmt.Println("Error:", err)
                    }
                } else {
                    fmt.Println("Unknown command")
                }
            } else {
                fmt.Println("You cannot enter empty command!")
            }
        }
    }
}


func cleanInput(text string) []string {
    text = strings.ToLower(text)
    return strings.Fields(text)
}
