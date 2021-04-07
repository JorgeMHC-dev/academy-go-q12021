package pokemon

// CsvPokemon - Defines the properties of a Pokemon to be listed
type CsvPokemon struct {
	ID    int `json:"ID,omitempty"`
	Name  string `json:"Name,omitempty"`
}

// Response - A Response struct to map the Entire Response from API
type Response struct {
    Name    string    `json:"name"`
    Pokemon []Pokemon `json:"pokemon_entries"`
}

// Pokemon - A Pokemon Struct to map every pokemon to from API.
type Pokemon struct {
    EntryNo int            `json:"entry_number"`
    Species PokemonSpecies `json:"pokemon_species"`
}

// PokemonSpecies - A struct to map our Pokemon's Species which includes it's name from API
type PokemonSpecies struct {
    Name string `json:"name"`
}
// PokemonResponse - A struct to get the entire response from API
type PokemonResponse struct {
    Name string `json:"name"`
    ID int `json:"id"`
    Sprites PokemonSprite `json:"sprites"`
}

// PokemonSprite - A struct to get the sprites of the response 
type PokemonSprite struct {
    Other PokemonOther `json:"other"`
}

// PokemonOther - A struct to get the sprites of the response 
type PokemonOther struct {
    Official PokemonFront `json:"official-artwork"`
}

// PokemonFront - A struct to get the sprites of the response 
type PokemonFront struct {
    Front string `json:"front_default"`
}