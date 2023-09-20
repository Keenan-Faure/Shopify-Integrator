package iocsv

import (
	"errors"
	"objects"
	"os"

	"github.com/gocarina/gocsv"
)

// Reads a csv file contents
func ReadFile(file_name string) ([]objects.CSVProduct, error) {
	if file_name == "" {
		return []objects.CSVProduct{}, errors.New("invalid file")
	}
	file_data, err := os.Open(file_name + ".csv")
	if err != nil {
		return []objects.CSVProduct{}, err
	}
	defer file_data.Close()
	products := []objects.CSVProduct{}
	if err := gocsv.UnmarshalFile(file_data, &products); err != nil {
		return []objects.CSVProduct{}, err
	}
	return products, nil
}

// Removes a file from the server
func RemoveFile(file_name string) error {
	err := os.Remove(file_name + ".csv")
	if err != nil {
		return err
	}
	return nil
}
