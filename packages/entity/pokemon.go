package pokemon

// CsvPokemon - Defines the properties of a Pokemon to be listed
type CsvPokemon struct {
	ID    int `json:"ID,omitempty"`
	Name  string `json:"Name,omitempty"`
}

// A Response struct to map the Entire Response from API
type Response struct {
    Name    string    `json:"name"`
    Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon Struct to map every pokemon to from API.
type Pokemon struct {
    EntryNo int            `json:"entry_number"`
    Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map our Pokemon's Species which includes it's name from API
type PokemonSpecies struct {
    Name string `json:"name"`
}
