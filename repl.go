package main
import(
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokédex... Goodbye!")
	os.Exit(0)
}

func commandHelp() error{
	fmt.Println("Welcome to the Pokédex!")
	fmt.Println("Usage:")
	for _, cmd := range(registry){
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
	
}

var registry = map[string]cliCommand{
    "exit": {
        name:"exit",
        description:"Exit the Pokedex",
        callback:commandExit,
    },
		"help": {
				name:"help",
				description:"Prints a help message describing how to use the REPL",
				callback:commandHelp,
	}
}


func cleanInput (text string) []string{
	formatted := strings.Fields(strings.TrimSpace(strings.ToLower(text)))
	return formatted	
}

func startRepl () {
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
			err := value.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command. Type 'help' for command info.")
			continue
		}
	}
}


