package pokeapi


import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/shirinox/pokecache"
)

var cache = pokecache.NewCache(5 * time.Minute)

type Config struct {
	next     string
	previous string
}

func (c *Config) Next() (string, error) {
	if c.next == "" {
		return "", errors.New("No next link.")
	}
	return c.next, nil
}

func (c *Config) Previous() (string, error) {
	if c.previous == "" {
		return "", errors.New("No previous link.")
	}
	return c.previous, nil
}

type MapListResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
}

func checkCache(conf *Config, url string) error {
	if data, ok := cache.Get(url); ok {
		mapList := MapListResponse{}
		err := json.Unmarshal(data, &mapList)

		if err != nil {
			return err
		}

		conf.next = mapList.Next
		conf.previous = mapList.Previous

		for _, location := range mapList.Results {
			fmt.Println(location.Name)
		}
		return nil
	}
	return errors.New("Data does not exist in cache")
}

type PokemonWithData struct {
	Name       string `json:"name"`
	Experience int    `json:"base_experience"`
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`
	Stats      []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Pokemon struct {
	Name string `json:"name"`
}

type PokemonArea struct {
	Name     string `json:"name"`
	Pokemons []struct {
		Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func fetch[T any](url string) (T, error) {
	var result T
	res, err := http.Get(url)
	if err != nil {
		return result, errors.New("There was an error while fetching from the API")
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}

	cache.Add(url, body)

	return result, nil

}

func CommandMap(conf *Config, args []string) error {
	url, err := conf.Next()
	if err != nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	err = checkCache(conf, url)
	if err != nil {
		fmt.Println(err)
	}

	res, err := fetch[MapListResponse](url)

	conf.next = res.Next
	conf.previous = res.Previous

	if err != nil {
		return err
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandMapBack(conf *Config, args []string) error {
	url, err := conf.Previous()

	err = checkCache(conf, url)
	res, err := fetch[MapListResponse](url)

	conf.next = res.Next
	conf.previous = res.Previous

	if err != nil {
		return err
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func CommandExplore(conf *Config, args []string) error {
	if len(args) != 1 {
		return errors.New("check args")
	}
	area := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area)

	pokemonArea, err := fetch[PokemonArea](url)
	if err != nil {
		return err
	}
	fmt.Printf("Explore %s\n", area)
	fmt.Println("Found Pokémon:")
	for _, p := range pokemonArea.Pokemons {
		fmt.Println("-", p.Pokemon.Name)
	}
	return nil
}

var Pokedex map[string]PokemonWithData

func CommandCatch(conf *Config, args []string) error {
	if len(args) != 1 {
		return errors.New("check args")
	}
	pokemon := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	p, err := fetch[PokemonWithData](url)

	if err != nil {
		fmt.Printf("There was an error with the API!\n")
		return err
	}

	if rand.Intn(p.Experience)/2 < p.Experience/4 {
		fmt.Printf("%v was caught!\n", p.Name)
		if len(Pokedex) == 0 {
			Pokedex = map[string]PokemonWithData{}
		}
		Pokedex[p.Name] = p
		return nil
	}

	fmt.Printf("%v escaped!\n", p.Name)
	return nil
}

func CommandInspect(conf *Config, args []string) error {
	if len(args) != 1 {
		return errors.New("check args")
	}
	pokemon := args[0]
	p, ok := Pokedex[pokemon]
	if !ok {
		fmt.Print("you have not caught that pokemon\n")
		return errors.New("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Printf("Stats:\n")
	for _, item := range p.Stats {
		fmt.Printf("   -%s: %d\n", item.Stat.Name, item.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, item := range p.Types {
		fmt.Printf("   - %s\n", item.Type.Name)
	}
	return nil

}

func CommandPokedex(conf *Config, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range Pokedex {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}
