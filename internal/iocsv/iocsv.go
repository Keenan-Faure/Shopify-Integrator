package iocsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"objects"
	"os"
	"strings"
	"utils"

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
	file_data2, err := os.Open(file_name + ".csv")
	if err != nil {
		return []objects.CSVProduct{}, err
	}
	defer file_data.Close()
	defer file_data2.Close()
	fileReader := csv.NewReader(file_data)
	records, err := fileReader.ReadAll()
	if err != nil {
		return []objects.CSVProduct{}, err
	}
	products := []objects.CSVProduct{}
	returned_products := []objects.CSVProduct{}
	qty_header_map := make(map[int]string)
	price_header_map := make(map[int]string)
	if err := gocsv.UnmarshalFile(file_data2, &products); err != nil {
		return []objects.CSVProduct{}, err
	}
	for key, value := range records {
		if key == 0 {
			qty_header_map = GetKeysByMatcher(value, "price_")
			price_header_map = GetKeysByMatcher(value, "qty_")
			continue
		}
		break
	}
	for key := range records {
		fmt.Println(key)
		if key == 0 {
			continue
		}
		qty := []objects.CSVQuantity{}
		pricing := []objects.CSVPricing{}
		for qty_key, qty_value := range qty_header_map {
			qty = append(qty, objects.CSVQuantity{
				Name:  qty_value,
				Value: utils.IssetInt(records[key][qty_key]),
			})
		}
		for price_key, price_value := range price_header_map {
			pricing = append(pricing, objects.CSVPricing{
				Name:  price_value,
				Value: utils.IssetString(records[key][price_key]),
			})
		}
		returned_products = append(returned_products, objects.CSVProduct{
			ProductCode:  products[key-1].ProductCode,
			Active:       products[key-1].Active,
			Title:        products[key-1].Title,
			BodyHTML:     products[key-1].BodyHTML,
			Category:     products[key-1].Category,
			Vendor:       products[key-1].Vendor,
			ProductType:  products[key-1].ProductType,
			SKU:          products[key-1].SKU,
			Option1Name:  products[key-1].Option1Name,
			Option1Value: products[key-1].Option1Value,
			Option2Name:  products[key-1].Option2Name,
			Option2Value: products[key-1].Option2Value,
			Option3Name:  products[key-1].Option3Name,
			Option3Value: products[key-1].Option3Value,
			Barcode:      products[key-1].Barcode,
			Warehouses:   qty,
			Pricing:      pricing,
		})
	}
	return returned_products, nil
}

// Removes a file from the server
func RemoveFile(file_name string) error {
	err := os.Remove(file_name + ".csv")
	if err != nil {
		return err
	}
	return nil
}

// returns the keys of all items
// in an array that matches a string
func GetKeysByMatcher(headers []string, match string) map[int]string {
	matcher := make(map[int]string)
	for key, header := range headers {
		if strings.Contains(header, match) {
			matcher[key] = header[len(match):]
		}
	}
	return matcher
}

// TODO does not want to update even though it exists inside the object
