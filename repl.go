package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

var exitFunc = os.Exit

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	NextURL *string
	PrevURL *string
	Cache   *pokecache.Cache
	Client  *pokeapi.Client
}


func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	exitFunc(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	var url string
	if cfg.NextURL != nil {
		url = *cfg.NextURL
	}

	if cachedData, ok := cfg.Cache.Get(url); ok {
		var locList pokeapi.LocationAreaList
		if err := json.Unmarshal(cachedData, &locList); err != nil {
			return err
		}
		for _, item := range locList.Results {
			fmt.Println(item.Name)
		}
		cfg.NextURL = locList.Next
		cfg.PrevURL = locList.Previous
		return nil
	}

	locList, err := cfg.Client.ListLocationAreas(url)
	if err != nil {
		return err
	}

	if cachedBytes, err := json.Marshal(locList); err == nil {
		if url == "" {
			url = "https://pokeapi.co/api/v2/location-area"
		}
		cfg.Cache.Add(url, cachedBytes)
	}

	for _, item := range locList.Results {
		fmt.Println(item.Name)
	}
	cfg.NextURL = locList.Next
	cfg.PrevURL = locList.Previous

	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.PrevURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	url := *cfg.PrevURL
	if url == "https://pokeapi.co/api/v2/location-area?offset=0&limit=20" {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	if cachedData, ok := cfg.Cache.Get(url); ok {
		var locList pokeapi.LocationAreaList
		if err := json.Unmarshal(cachedData, &locList); err != nil {
			return err
		}
		for _, item := range locList.Results {
			fmt.Println(item.Name)
		}
		cfg.NextURL = locList.Next
		cfg.PrevURL = locList.Previous
		return nil
	}

	locList, err := cfg.Client.ListLocationAreas(url)
	if err != nil {
		return err
	}

	if cachedBytes, err := json.Marshal(locList); err == nil {
		cfg.Cache.Add(url, cachedBytes)
	}

	for _, item := range locList.Results {
		fmt.Println(item.Name)
	}
	cfg.NextURL = locList.Next
	cfg.PrevURL = locList.Previous

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Explore command missing location area input. eg. explore <location_area>")
		return nil
	}
	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Println("Found Pokemon:")
	url := "https://pokeapi.co/api/v2/location-area/" + areaName
	
	
	if cachedData, ok := cfg.Cache.Get(url); ok {
		var pokemonList pokeapi.LocationArea
		if err := json.Unmarshal(cachedData, &pokemonList); err != nil {
			return err
		}
		for _, item := range pokemonList.PokemonEncounters {
			fmt.Println(item.Pokemon.Name)
		}
		return nil
	}

	pokemonList, err := cfg.Client.GetLocationArea(areaName)
	if err != nil {
		return err
	}

	if cachedBytes, err := json.Marshal(pokemonList); err == nil {
		if url == "" {
			url = "https://pokeapi.co/api/v2/location-area/" + areaName
		}
		cfg.Cache.Add(url, cachedBytes)
	}

	for _, item := range pokemonList.PokemonEncounters {
		fmt.Println(item.Pokemon.Name)
	}
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
		"explore": {
			name: "explore",
			description: "Lists all Pokemon that can be found at the target location area.",
			callback: commandExplore,
		},
	}
}

func cleanInput(text string) []string {
	formatted := strings.Fields(strings.TrimSpace(strings.ToLower(text)))
	return formatted
}

func startRepl() {
	cfg := &config{
		Cache:  pokecache.NewCache(30 * time.Second),
		Client: pokeapi.NewClient(),
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
		args := words[1:]
		value, exists := registry[command]
		if exists {
			err := value.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command. Type 'help' for command info.")
			continue
		}
	}
}
