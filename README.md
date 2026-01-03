# Pokedex CLI

A command-line interface (CLI) tool that acts as a Pokedex, allowing users to explore the Pokemon world, catch Pokemon, and inspect their details using the [PokeAPI](https://pokeapi.co/).

## Motivation

This application provides a robust and interactive command-line interface for exploring data from the PokeAPI. It features a custom internal caching system to optimize network requests and ensure a responsive user experience. The architecture creates a clear separation of concerns between the API client, caching layer, and the REPL interface, prioritizing maintainability and extensibility. It handles complex data relationships and state management to simulate a fully functional Pokedex environment.

## Quick Start

Ensure you have [Go](https://go.dev/dl/) installed on your machine.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/limitedink/pokedex-cli.git
    cd pokedex-cli
    ```

2.  **Build the project:**

    ```bash
    go build -o pokedexcli
    ```

3.  **Run the executable:**

    ```bash
    ./pokedexcli
    ```

## Usage

Once the REPL is running, you can use the following commands to interact with the Pokedex:

*   **`help`**: Displays a help message describing how to use the REPL.
*   **`map`**: Displays the names of the next 20 location areas in the Pokemon world.
*   **`mapb`**: Displays the previous 20 location areas.
*   **`explore <location_area>`**: Lists all Pokemon that can be found at the specified location area.
*   **`catch <pokemon_name>`**: Attempts to catch a Pokemon. The chance of catching depends on the Pokemon's base experience.
*   **`inspect <pokemon_name>`**: Displays detailed information (height, weight, stats, types) about a Pokemon you have caught.
*   **`pokedex`**: Lists all the Pokemon you have currently caught.
*   **`exit`**: Exits the Pokedex.

### Example Session

```text
Pokédex > map
pastoria-city-area
canalave-city-area
...

Pokédex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
- bibarel
- wingull
...

Pokédex > catch bibarel
Throwing a Pokeball at bibarel...
bibarel was caught!

Pokédex > inspect bibarel
Name:  bibarel
Height:  10
Weight:  315
Stats:
  -hp: 79
  -attack: 85
...

Pokédex > pokedex
Your Pokedex:
 - bibarel
```

## Contributing

Contributions are welcome! If you'd like to improve the project, please follow these steps:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/YourFeature`).
3.  Make your changes and commit them (`git commit -m 'Add some feature'`).
4.  Push to the branch (`git push origin feature/YourFeature`).
5.  Open a Pull Request.
