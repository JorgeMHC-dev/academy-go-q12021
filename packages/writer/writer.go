package writer

import (
	"encoding/csv"
	"os"
)

//WriteData - writes data to the csv file
func WriteData(data [][]string, fileName string) error {
	f,err := os.Create(fileName)

	if err != nil {
		return err
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
 
	for _, csvRow := range data {
		_ = csvwriter.Write(csvRow)
	}

	csvwriter.Flush()

	return nil
}