package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
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
	Caught  map[string]pokeapi.Pokemon
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

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Catch command missing pokemon name input. eg. catch <pokemon_name>")
		return nil
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := cfg.Client.GetPokemon(pokemonName)
	if err != nil {
		fmt.Println("Can't find that Pokemon here.")
		return nil
	}
	roll := rand.Float64()
	minExp := 50.0
	maxExp := 300.0
	minProb := 0.2
	maxProb := 0.9

	exp := float64(pokemon.BaseExperience)
	if exp < minExp {
		exp = minExp
	}
	if exp > maxExp {
		exp = maxExp
	}

	t := (exp - minExp) / (maxExp - minExp)

	prob := maxProb*(1.0-t) + minProb*t

	if roll > prob {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	cfg.Caught[pokemon.Name] = *pokemon
	fmt.Printf("%s was caught!\n", pokemonName)

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
			name:        "explore",
			description: "Lists all Pokemon that can be found at the target location area.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon.",
			callback:    commandCatch,
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
		Caught: make(map[string]pokeapi.Pokemon),
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
