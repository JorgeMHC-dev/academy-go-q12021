package external

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	pokemon "github.com/jorgemhc-dev/academy-go-q12021/packages/entity"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/reader"
)

func TestUpdateCsvPokedex(t *testing.T) {
	fmt.Println("TessUpdateCsvPokedex")

	mockResponse := pokemon.Response{
		Name: "Kanto",
		Pokemon: []pokemon.Pokemon{
			pokemon.Pokemon{
				EntryNo: 1,
				Species: pokemon.PokemonSpecies{
					Name: "Charmander",
				},
			},
		},
	}

	jsn,erro := json.Marshal(mockResponse)

	if erro != nil {
		t.Errorf("Wrong object")
	}

	httpmock.Activate()
    httpmock.RegisterResponder("GET", "http://pokeapi.co/api/v2/pokedex/1000",
    httpmock.NewStringResponder(200, string(jsn)))

	err := UpdateCsvPokedex("1000")

	httpmock.DeactivateAndReset()

	records,err := reader.ReadData("./pokemons.csv")

	var response pokemon.Response

	for _,record := range records {
		id,_ := strconv.Atoi(record[0])

		response = pokemon.Response{
			Name: "Kanto",
			Pokemon: []pokemon.Pokemon{
				pokemon.Pokemon{
					EntryNo: id,
					Species: pokemon.PokemonSpecies{
						Name: record[1],
					},
				},
			},
		}
	}

	if err != nil {
		t.Error(err)
	} else if !cmp.Equal(response,mockResponse) {
		t.Errorf("Wrong save on CSV response: %v expected: %v", response,mockResponse)
	}

	fmt.Printf("Record saved correctly on CSV: %v", records)
}