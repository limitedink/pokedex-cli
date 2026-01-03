# Pokedex CLI

**High-performance Pokedex REPL built in Go, featuring a custom thread-safe caching system with asynchronous TTL eviction to optimize API latency.**

## Motivation

This project was engineered to explore advanced Go concepts in a practical environment. The primary goal was to build a stateful CLI application that handles real-world constraints—network latency, data persistence, and concurrency.

Key architectural goals included:
*   **Separation of Concerns:** Distinct modules for API interaction (`pokeapi`), data caching (`pokecache`), and the REPL interface.
*   **Concurrency:** Implementing a thread-safe cache using `sync.Mutex` and background goroutines for automatic data reaping.
*   **Performance:** Minimizing expensive network calls to the PokeAPI by serving repeated requests from memory.

## Technical Implementation

### Custom Caching System
To ensure the application remains responsive, I implemented a custom in-memory caching layer (`internal/pokecache`).
*   **Thread Safety:** Utilizes `sync.Mutex` to ensure safe concurrent access to the map, preventing race conditions.
*   **Automatic Eviction:** A background goroutine (ticker) runs periodically to "reap" or remove entries that have exceeded their time-to-live (TTL), preventing memory bloat during long sessions.

### REPL Architecture
The Read-Eval-Print Loop is designed to be extensible. It uses a registry pattern to map string commands to function callbacks, making it trivial to add new features without modifying the core input loop. Input is sanitized and tokenized to handle variable arguments robustly.

## Quick Start

### Prerequisites
*   **Go 1.20+**

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
*   **`map` / `mapb`**: Navigate through location areas in paginated batches (cached for speed).
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

## Future Enhancements & Ideas

The project is designed to be extensible. Here are some potential features for future development or forks:

*   **Persistence:** Save the Pokedex state to disk (JSON/SQLite) to persist progress between sessions.
*   **Battle System:** Simulate turn-based battles between your caught Pokemon and wild encounters.
*   **RPG Elements:** Implement an XP system to level up Pokemon and an evolution mechanic based on time or battles.
*   **Inventory System:** Add support for different Pokeball types (Great Ball, Ultra Ball) with varying catch rates.
*   **Improved UX:** Add command history (up-arrow support) and a more immersive navigation system (e.g., "move left/right").
*   **Testing:** Expand unit test coverage and refactor for greater testability.

## Contributing

Contributions are welcome! If you'd like to improve the project, please follow these steps:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/YourFeature`).
3.  Make your changes and commit them (`git commit -m 'Add some feature'`).
4.  Push to the branch (`git push origin feature/YourFeature`).
5.  Open a Pull Request.
