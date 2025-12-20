package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"pokedexcli/internal/pokecache"
)

var exitFunc = os.Exit

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	NextURL *string
	PrevURL *string
	Cache *pokecache.Cache
}

type locationAreaList struct {
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []locationAreaItem `json:"results"`
}

type locationAreaItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	exitFunc(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.NextURL != nil {
        url = *cfg.NextURL
  }
	fmt.Println("commandMap URL:", url)
	var bodyBytes []byte

	if data, ok := cfg.Cache.Get(url); ok {
	bodyBytes = data
  } else {
			fmt.Println("making HTTP request")
			res, err := http.Get(url)
			if err != nil {
					return err
			}
			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			if err != nil {
					return err
			}
			bodyBytes = b
			
			cfg.Cache.Add(url, bodyBytes)
  }	
	var locList locationAreaList
	err := json.Unmarshal(bodyBytes, &locList)
	if err != nil {
		return err
	}

	for _, item := range locList.Results {
		fmt.Println(item.Name)
	}
	cfg.NextURL = locList.Next
	cfg.PrevURL = locList.Previous

	return nil
}

func commandMapb(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.PrevURL == nil {
        fmt.Println("you're on the first page")
        return nil
    }
  url = *cfg.PrevURL
	if url == "https://pokeapi.co/api/v2/location-area?offset=0&limit=20" {
        url = "https://pokeapi.co/api/v2/location-area"
  }

	var bodyBytes []byte

	if data, ok := cfg.Cache.Get(url); ok {
	bodyBytes = data
  } else {
			fmt.Println("making HTTP req")
			res, err := http.Get(url)
			if err != nil {
					return err
			}
			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			if err != nil {
					return err
			}
			bodyBytes = b

			cfg.Cache.Add(url, bodyBytes)
  }	
	
	var locList locationAreaList
	err := json.Unmarshal(bodyBytes, &locList)
	if err != nil {
		return err
	}

	for _, item := range locList.Results {
		fmt.Println(item.Name)
	}
  cfg.NextURL = locList.Next
	cfg.PrevURL = locList.Previous

	return nil
}

var registry = map[string]cliCommand{}

func init() {
	registry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex.",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Prints a help message describing how to use the REPL.",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Like the map command but displays the previous 20 locations/previous page of locations.",
			callback:    commandMapb,
		},
	}
}

func cleanInput(text string) []string {
	formatted := strings.Fields(strings.TrimSpace(strings.ToLower(text)))
	return formatted
}

func startRepl() {
	cfg := &config{
		Cache: pokecache.NewCache(30 * time.Second),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokédex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		command := words[0]
		value, exists := registry[command]
		if exists {
			err := value.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command. Type 'help' for command info.")
			continue
		}
	}
}
