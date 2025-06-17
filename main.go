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

	"github.com/chandanbsd/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
	config      *config
}

type config struct {
	Next     *string
	Previous *string
}

var c *config = &config{
	Previous: nil,
	Next:     nil,
}

type result struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaPokemon struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

var commands map[string]cliCommand

var cache *pokecache.Cache = pokecache.NewCache(5 * time.Millisecond)

func init() {
	commands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Lists the locations",
			callback:    commandLocationAreaNext,
			config:      c,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the locations back",
			callback:    commandLocationAreaPrevious,
			config:      c,
		},
		"explore": {
			name:        "explore",
			description: "Used to explore the pokemons at the given location",
			callback:    commandExplore,
			config:      c,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      c,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      c,
		},
	}
}

func main() {

	var word string
	scanner := bufio.NewScanner(os.Stdin)

	var cleanedSlice []string
	for true {
		fmt.Printf("Pokedex > ")
		isComplete := scanner.Scan()

		if isComplete {
			word = scanner.Text()
			cleanedSlice = cleanInput(word)

			if len(cleanedSlice) > 0 {
				commandStr := cleanedSlice[0]

				command, ok := commands[commandStr]

				if !ok {
					fmt.Println("Unknown command")
				}
				
				if len(cleanedSlice) == 2 {
					command.callback(cleanedSlice[1])
				} else {
					command.callback("")
				}
			}

		}

	}
}

func commandExit(area_name string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(area_name string) error {

	fmt.Println("Welcome to the Pokedex!\nUsage: ")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandLocationAreaNext(area_name string) error {

	config := commands["map"].config

	basedUrl := "https://pokeapi.co/api/v2/location-area/"

	if config.Next == nil {
		config.Next = &basedUrl
	}

	cacheRes, ok := cache.Get(basedUrl)

	var locationArea locationArea
	var err error

	if ok {
		locationArea, err = printHelper(cacheRes)
		if err != nil {
			return err
		}
	} else {
		bytes, err := fetchHelper(*config.Next)
		if err != nil {
			return err
		}

		locationArea, err = printHelper(bytes)
		if err != nil {
			return err
		}
	}

	config.Previous = config.Next
	config.Next = locationArea.Next

	//decoder := json.NewDecoder(res.Body)

	//err = decoder.Decode(&locationArea)

	return nil
}

func fetchHelper(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func printHelper(val []byte) (locationArea, error) {

	var locationArea locationArea = locationArea{}

	err := json.Unmarshal(val, &locationArea)
	if err != nil {
		return locationArea, err
	}

	for _, res := range locationArea.Results {
		fmt.Println(res.Name)
	}

	return locationArea, nil

}

func commandLocationAreaPrevious(area_name string) error {

	config := commands["map"].config

	basedUrl := "https://pokeapi.co/api/v2/location-area/"
	config.Next = &basedUrl

	if config.Previous == nil {
		fmt.Println("you're )on the first page")
		return nil
	}

	cacheRes, ok := cache.Get(basedUrl)

	var locationArea locationArea
	var err error

	if ok {
		locationArea, err = printHelper(cacheRes)
		if err != nil {
			return err
		}
	} else {
		bytes, err := fetchHelper(*config.Next)
		if err != nil {
			return err
		}

		locationArea, err = printHelper(bytes)
		if err != nil {
			return err
		}
	}

	config.Previous = locationArea.Previous
	config.Next = config.Previous
	return nil
}

func cleanInput(text string) []string {
	res := strings.Fields(text)

	values := []string{}

	for _, s := range res {
		values = append(values, strings.ToLower(s))
	}

	return values
}

func printPokemonHelper(val []byte) {

	var locationArea LocationAreaPokemon = LocationAreaPokemon{}

	err := json.Unmarshal(val, &locationArea)
	if err != nil {
		return
	}

	for _, res := range locationArea.PokemonEncounters {
		fmt.Println(res.Pokemon.Name)
	}
}

func commandExplore(area_name string) error{
	url := "https://pokeapi.co/api/v2/location-area/" + area_name
	
	res, err := fetchHelper(url)
	if err != nil {
		return err
	}

	printPokemonHelper(res) 

	return  nil
}