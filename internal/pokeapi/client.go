package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    "https://pokeapi.co/api/v2",
	}
}

type LocationAreaList struct {
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []LocationAreaItem `json:"results"`
}

type LocationAreaItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
    Name              string `json:"name"`
    PokemonEncounters []struct {
        Pokemon struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"pokemon"`
    } `json:"pokemon_encounters"`
}

func (c *Client) ListLocationAreas(url string) (*LocationAreaList, error) {
	if url == "" {
		url = c.baseURL + "/location-area"
	}

	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var locList LocationAreaList
	err = json.Unmarshal(body, &locList)
	if err != nil {
		return nil, err
	}

	return &locList, nil
}

func (c *Client) GetLocationArea(name string) (*LocationArea, error) {
    url := c.baseURL + "/location-area/" + name

    res, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    var loc LocationArea
    if err := json.Unmarshal(body, &loc); err != nil {
        return nil, err
    }

    return &loc, nil
}
