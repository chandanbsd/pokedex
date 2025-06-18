package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        interface{} `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       string      `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      string      `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale string      `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var bag map[string]Pokemon = map[string]Pokemon{}

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
		"catch": {
			name:        "catch",
			description: "catches a pokemon",
			callback:    commandCatch,
			config:      c,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Pokedex",
			callback:    commandPokedex,
			config:      c,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect the pokemon",
			callback:    commandInspect,
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
		config.Previous = nil
	}

	cacheRes, ok := cache.Get(*config.Next)

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
	config.Next = locationArea.Next

	return nil
}

func fetchHelper(url string) ([]byte, error) {
	client := &http.Client{}

	cacheRes, ok := cache.Get(url)

	if ok {
		return cacheRes, nil
	}


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

	cache.Add(url, bytes)

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

	config := commands["mapb"].config
	
	basedUrl := "https://pokeapi.co/api/v2/location-area/"

	if config.Previous == nil {
		fmt.Println("you're on the first page")
		config.Next = &basedUrl
		return nil
	}

	cacheRes, ok := cache.Get(*config.Previous)

	var locationArea locationArea
	var err error

	if ok {
		locationArea, err = printHelper(cacheRes)
		if err != nil {
			return err
		}
	} else {
		bytes, err := fetchHelper(*config.Previous)
		if err != nil {
			return err
		}

		locationArea, err = printHelper(bytes)
		if err != nil {
			return err
		}
	}

	config.Previous = locationArea.Previous
	config.Next = locationArea.Next

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

var attemptedCatches map[string]int = map[string]int{}

func commandCatch(pokemonName string) error {

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	resBytes, err := fetchHelper(url)
	if err != nil {
		return err
	}

	var pokemon Pokemon

	err = json.Unmarshal(resBytes, &pokemon)

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)

	if attemptedCatches[pokemonName] == 4 {
		bag[pokemonName] = pokemon
	} else if pokemon.BaseExperience < 100 && rand.Float64() < 0.90{
			fmt.Printf("%v was caught!\nYou may now inspect it with the inspect command.\n", pokemonName)
		bag[pokemonName] = pokemon

	} else if pokemon.BaseExperience < 200 && rand.Float64() < 0.75 {
			fmt.Printf("%v was caught!\n", pokemonName)
			bag[pokemonName] = pokemon

	} else if pokemon.BaseExperience < 300 && rand.Float64() < 0.5  {
		fmt.Printf("%v was caught!\n", pokemonName)
			bag[pokemonName] = pokemon
	} else if pokemon.BaseExperience >= 300 && rand.Float64() < 0.25 {
			fmt.Printf("%v was caught!\n", pokemonName)
			bag[pokemonName] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
		attemptedCatches[pokemonName] += 1
	}

	return nil
}

func commandPokedex(ignoreArg string) error {
	fmt.Printf("Your Pokedex:\n")

	for key, _ := range bag {
		fmt.Printf(" - %v\n", key)
	}
	return nil
}

func commandInspect(pokemonName string) error {
	
	pokemon, ok := bag[pokemonName]

	if !ok {
		return nil
	}

	fmt.Printf(`
Name: %v
Height: %v
Weight: %v
`, pokemonName, pokemon.Height, pokemon.Weight);


	fmt.Println("Stats:")
	interestedStats := map[string]int{
		"hp" : 0,
		"attack": 0,
		"defense": 0,
		"special-attack": 0,
		"special-defense": 0,
		"speed": 0,
	}

	for _, stat := range pokemon.Stats {

		_, ok := interestedStats[stat.Stat.Name]

		if ok {
			interestedStats[stat.Stat.Name] = stat.BaseStat
		}
	}

	for key, value := range interestedStats {
		fmt.Printf(" -%s: %v\n", key, value)
	}


	fmt.Println("Types:")
	
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	return nil
}