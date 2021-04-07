package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	pokemon "github.com/jorgemhc-dev/academy-go-q12021/packages/entity"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/writer"
)

func UpdateCsvPokedex(pokedex interface{}) error {
	Dex := "http://pokeapi.co/api/v2/pokedex/"
	switch v := pokedex.(type) {
	case string:
		Dex = Dex+v
	case int:
		Dex = Dex+strconv.Itoa(v)
	}
	response, err := http.Get(Dex)
    if err != nil {
        return err
    } else if response.StatusCode == http.StatusNotFound{
		return fmt.Errorf("Pokedex not found")
	} else if response.StatusCode != http.StatusOK{
		return fmt.Errorf("Error in api")
	}

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return err
    }
	var responseObject pokemon.Response
	
	json.Unmarshal(responseData, &responseObject)

	newCsvData := make([][]string,0)

	newCsvData = append(newCsvData,[]string{"ID","Name"})
	for i := range responseObject.Pokemon {
		newCsvData = append(newCsvData,[]string{strconv.Itoa(responseObject.Pokemon[i].EntryNo), responseObject.Pokemon[i].Species.Name})
	}
	fail := writer.WriteData(newCsvData,"./pokemons.csv")

	if fail != nil {
		return fail
	}

    return nil
}

func FetchInfo(ID interface{}) (pokemon.PokemonResponse,error) {

	Poke := "http://pokeapi.co/api/v2/pokemon/"
	switch v := ID.(type) {
	case string:
		Poke = Poke+v
	case int:
		Poke = Poke+strconv.Itoa(v)
	}
	response, err := http.Get(Poke)
    if err != nil {
        return pokemon.PokemonResponse{},err
    } else if response.StatusCode == http.StatusNotFound{
		return pokemon.PokemonResponse{},fmt.Errorf("Pokedex not found")
	} else if response.StatusCode != http.StatusOK{
		return pokemon.PokemonResponse{},fmt.Errorf("Error in api")
	}

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return pokemon.PokemonResponse{},err
    }

	var responseObject pokemon.PokemonResponse

	json.Unmarshal(responseData, &responseObject)

	return responseObject,nil
}