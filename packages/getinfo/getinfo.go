package getinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	pokemon "github.com/jorgemhc-dev/academy-go-q12021/packages/entity"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/external"
	"github.com/jorgemhc-dev/academy-go-q12021/packages/reader"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./templates/*.html",
	}

	rnd = renderer.New(opts)
}

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

func getPokemon(ID string) (pokemon.CsvPokemon,error) {
	records,err := reader.ReadData("./pokemons.csv")

	var result pokemon.CsvPokemon

	if err != nil {
		return pokemon.CsvPokemon{},err
	}
	for _, record := range records {

		id :=  record[0]
		name := record[1]
		if id == ID || name == ID {
			idI,_ := strconv.Atoi(id)
			result = pokemon.CsvPokemon{
				ID : idI,
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

//FetchAll - Obtains all the pokemons in the list
func FetchAll(w http.ResponseWriter, r *http.Request) {
	pokemons,err := getPokemons()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pokemons)
	}
}

//FetchPokemon - Obtains one pokemon based on the ID or Name
func FetchPokemon(w http.ResponseWriter, r *http.Request) {
	vars:= mux.Vars(r)

	id,_ := vars["ID"]
	pokemon,err := getPokemon(id)
	
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Pokemon Not found")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pokemon)
}

//FetchPokemonInfo - obtains the pokemon info from the api
func FetchPokemonInfo(w http.ResponseWriter, r *http.Request){
	vars:= mux.Vars(r)

	poke,err := external.FetchInfo(vars["ID"])

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		rnd.HTML(w, http.StatusNotFound, "notFound", nil)
		return
	}
	rnd.HTML(w, http.StatusOK, "body", poke)

}

//UpdateCsv - updates the csv with a pokedex obtained from the pokeapi endpoint
func UpdateCsv(w http.ResponseWriter, r *http.Request){
	vars:= mux.Vars(r)

	err := external.UpdateCsvPokedex(vars["POKEDEX"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
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

// ObtaingPokemonConcurrent - Get the vars of the URL and runs our go routine and obtains all the results sending them in the response header
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
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(pokemons)
			}	
		case erro := <-err:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(erro)
			continue
		}
	}
}


// Worker - Initialize all the workers running on our go routine
func Worker(jobs <- chan int, result chan<- []pokemon.CsvPokemon,erro chan<- error, typeS string, items int, items_per_worker int, workers int){
	for i:= range jobs {
		response,err := getPokemonsByArg(typeS,items, items_per_worker,i, workers)
		if err != nil {
			erro <- err 
			close(result)
			break
		}
		result <- response
	}	
}