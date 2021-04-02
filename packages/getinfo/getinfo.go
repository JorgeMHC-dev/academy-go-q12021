package getinfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	pokemon "github.com/jorgemhc-dev/academy-go-q12021/packages/entity"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/reader"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/writer"
)

func getPokemons() ([]pokemon.CsvPokemon, error) {
	records,err := reader.ReadData("./pokemons.csv")

	var result []pokemon.CsvPokemon

	if err != nil {
		return nil,err
	}

	for _, record := range records {

		id,_ :=  strconv.Atoi(record[0])

		pokemon := pokemon.CsvPokemon{
			ID : id,
			Name: record[1],
		}

		result = append(result,pokemon)
	}

	return result,nil
}

func getPokemon(ID int) (pokemon.CsvPokemon,error) {
	records,err := reader.ReadData("./pokemons.csv")

	var result pokemon.CsvPokemon

	if err != nil {
		return pokemon.CsvPokemon{},err
	}

	for _, record := range records {

		id,_ :=  strconv.Atoi(record[0])
		if id == ID {
			result = pokemon.CsvPokemon{
				ID : id,
				Name: record[1],
			}
		}
	}
	if result == (pokemon.CsvPokemon{}) {
		err := fmt.Errorf("Pokemon not found")
		return result,err
	}
	return result,nil

}

func getPokemonsByArg(T string, Items int, items_per_worker int, job int, workers int) ([]pokemon.CsvPokemon,error) {
	records,err := reader.ReadData("./pokemons.csv")

	var result []pokemon.CsvPokemon
	var res pokemon.CsvPokemon
	var index int

	if err != nil {
		return nil,err
	}
	switch T {
	case "even":
			for _,record := range records {
				if index >= items_per_worker{
					break
				}
	
				id,_ :=  strconv.Atoi(record[0])
				if id%2 == 0 {
					index,res = processInfoJobs(job,workers,Items,items_per_worker,index,id,record[1])
					result = append(result,res);
					result = delete_empty(result)
				}
			}
		
	case "odd":
			for _, record := range records {
				if index >= items_per_worker{
					break
				}
				id,_ :=  strconv.Atoi(record[0])
				if id%2 != 0 {
					index,res = processInfoJobs(job,workers,Items,items_per_worker,index,id,record[1])
					result = append(result,res);
					result = delete_empty(result)
				}
			}
	default:
		return nil,fmt.Errorf("Not supported type")
	}

	return result,nil
}

func delete_empty (s []pokemon.CsvPokemon) []pokemon.CsvPokemon {
	var r []pokemon.CsvPokemon
	for _, pkmn := range s {
        if pkmn != (pokemon.CsvPokemon{}) {
			r = append(r, pkmn)
        }
	}
	return r
}

func processInfoJobs(job int, workers int, Items int, items_per_worker int, index int, id int, record string) (int, pokemon.CsvPokemon) {
	var result pokemon.CsvPokemon

	if job+1 == workers{ 
		nwIndex:= Items%items_per_worker
		if index == nwIndex && nwIndex>0{
			return items_per_worker,result
		}	
		nextId := (job+job) * items_per_worker
		if id > nextId{
			result = pokemon.CsvPokemon{
				ID : id,
				Name: record,
			}
			index++
			return index,result
		}
	} else if job > 0 {
		nextId := (job+job) * items_per_worker
		if id > nextId{
			result = pokemon.CsvPokemon{
				ID : id,
				Name: record,
			}
			index++
			return index,result
		}
	} else {
		result = pokemon.CsvPokemon{
			ID : id,
			Name: record,
		}
		index++
		return index,result
	}
	return index,result
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Pokemon Not found")
		return
	}

	json.NewEncoder(w).Encode(pokemon)
}

func updateCsvPokedex(pokedex interface{}) error {
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

func UpdateCsv(w http.ResponseWriter, r *http.Request){
	vars:= mux.Vars(r)

	err := updateCsvPokedex(vars["POKEDEX"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(http.StatusOK);
}

func splitFloat(n string) int {
	s := strings.Split(n,".")
	x := s[0]
	y := s[1]

	n1,_:= strconv.Atoi(x)
	n2,_:= strconv.Atoi(y)

	if n2>0 {
		n2 = 1
	}

	return n1+n2
}

func floorCeilFloat(n string) int {
	s := strings.Split(n,".")
	y := s[1]

	n2,_:= strconv.Atoi(y)

	var nw float64

	if n2<=5 && n2>0{
		nw,_ = strconv.ParseFloat(n,64)

		nw = math.Floor(nw)
	} else if n2 > 5 {
		nw,_ = strconv.ParseFloat(n,64)

		nw = math.Ceil(nw)
	} else {
		nw,_ = strconv.ParseFloat(n,64)
	}

	return int(nw)
}

func ObtaingPokemonConcurrent(w http.ResponseWriter, r *http.Request) {
	vars:= mux.Vars(r)

	items,_ := strconv.Atoi(vars["ITEMS"])
	items_per_worker,_  := strconv.Atoi(vars["ITEMS-PER-WORKER"])
	typeS := vars["TP"]
	var workers int
	if items <= items_per_worker {
		workers = 1
	} else {
		var res float64
		res = float64(items)/float64(items_per_worker)
		
		workers = splitFloat(fmt.Sprintf("%f",res))
	}

	jobs := make(chan int, workers)
	results := make(chan []pokemon.CsvPokemon, items)
	err := make(chan error)

	go Worker(jobs,results,err,typeS,items, items_per_worker, workers)

	for i := 0; i < workers; i++{
		jobs <- i
	}
	close(jobs)

	var pokemons [][]pokemon.CsvPokemon

	w.Header().Set("Content-Type", "application/json")

	for i:= 0; i<workers ;i++ {
		select {
		case results:= <-results:
			if results != nil{
				pokemons := append(pokemons,results)
				json.NewEncoder(w).Encode(pokemons)
			}	
		case erro := <-err:
			json.NewEncoder(w).Encode(erro)
			continue
		}
	}
}

func Worker(jobs <- chan int, result chan<- []pokemon.CsvPokemon,erro chan<- error, typeS string, items int, items_per_worker int, workers int){
	for i:= range jobs {
		response,err := getPokemonsByArg(typeS,items, items_per_worker,i, workers)
		fmt.Println(err!=nil)
		if err != nil {
			erro <- err 
			close(result)
			break
		}
		result <- response
	}	
}