package getinfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	pokemon "github.com/jorgemhc-dev/academy-go-q12021/packages/entity"
)

func TestGetPokemon(t *testing.T) {
	fmt.Println("TestGetPokemon")

	req,err := http.NewRequest("GET", "/pokemon/150",nil)

	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPokemon)

	vars := map[string]string{
        "ID": "150",
    }

	req = mux.SetURLVars(req, vars)

	handler.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler wrong status code: expected: %v got: %v", http.StatusOK, status)
	}

	expected := pokemon.CsvPokemon{
					ID:150,
					Name:"mewtwo",
				}

	responseData, _ := ioutil.ReadAll(rr.Body)
	
	var responseObject pokemon.CsvPokemon
	
	json.Unmarshal(responseData, &responseObject)

	var poke pokemon.CsvPokemon
	for range responseData {
		poke = pokemon.CsvPokemon{
			ID: responseObject.ID,
			Name : responseObject.Name,
		}	
	}

	if !cmp.Equal(poke,expected) {
		t.Errorf("Returned unexpected body: expected: %v got: %v",expected,rr.Body.String())
	}

	fmt.Printf("Pokemon obtained correctly: expected %v got: %v",expected,poke)
}