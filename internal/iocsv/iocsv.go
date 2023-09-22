package iocsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"objects"
	"os"
	"strings"
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
	fileReader := csv.NewReader(file_data)
	records, err := fileReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range records {
		fmt.Println(value)
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

// identifies the columns
// that holds the pricing & qty
func GetPricingHeaderKeys(headers []string) []int {
	pricing_header_keys := []int{}
	for key, header := range headers {
		if strings.Contains(header, "price_") {
			pricing_header_keys = append(pricing_header_keys, key)
		}
	}
	return pricing_header_keys
}
