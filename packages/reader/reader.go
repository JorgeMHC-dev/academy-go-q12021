package reader

import (
	"encoding/csv"
	"os"
)

//ReadData - reads all the data of a csv file
func ReadData(fileName string) ([][]string, error) {
	f,err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	if _,err := r.Read(); err != nil {
		return [][]string{},err
	}

	records,err := r.ReadAll()

	if err != nil {
		return [][]string{},err
	}

	return records,nil
}
