package getinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	pokemon "github.com/jorgemhc-dev/academy-go-q12021/packages/entity"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/reader"
)

func getPokemons() ([]pokemon.Pokemon, error) {
	records,err := reader.ReadData("./pokemons.csv")

	var result []pokemon.Pokemon

	if err != nil {
		return result,nil
	}

	for _, record := range records {

		id,_ :=  strconv.Atoi(record[0])

		pokemon := pokemon.Pokemon{
			ID : id,
			Name: record[1],
		}

		result = append(result,pokemon)
	}

	return result,nil
}

func getPokemon(ID int) (pokemon.Pokemon,error) {
	records,err := reader.ReadData("./pokemons.csv")

	var result pokemon.Pokemon

	if err != nil {
		return result,nil
	}

	for _, record := range records {

		id,_ :=  strconv.Atoi(record[0])
		if id == ID {
			result = pokemon.Pokemon{
				ID : id,
				Name: record[1],
			}
		}
	}
	if result == (pokemon.Pokemon{}) {
		err := fmt.Errorf("Pokemon not found")
		return result,err
	}
	return result,nil

}

func FetchAll(w http.ResponseWriter, r *http.Request) {
	pokemons,err := getPokemons()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(pokemons)
	}
}
func FetchPokemon(w http.ResponseWriter, r *http.Request) {
	vars:= mux.Vars(r)

	id,_ := strconv.Atoi(vars["ID"])
	pokemon,err := getPokemon(id)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // We use not found for simplicity
		json.NewEncoder(w).Encode("Pokemon Not found")
		return
	}

	json.NewEncoder(w).Encode(pokemon)
}