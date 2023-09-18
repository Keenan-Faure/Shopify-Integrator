package iocsv

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Reads a csv file contents
func ReadFile(file_name string) error {
	file_data, err := os.Open(file_name + ".csv")
	if err != nil {
		return err
	}
	defer file_data.Close()
	fileReader := csv.NewReader(file_data)
	records, err := fileReader.ReadAll()
	if err != nil {
		return err
	}
	// parse through data
	fmt.Println(records[0])
	return nil
}
