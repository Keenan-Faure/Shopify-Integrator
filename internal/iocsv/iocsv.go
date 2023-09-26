package iocsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"objects"
	"os"
	"reflect"
	"strings"
	"utils"

	"github.com/gocarina/gocsv"
)

// Writes data onto a csv file
func WriteCSV(file_name string, product_data []objects.Product) error {
	// need to generate csv headers that are dynamic
	fmt.Println(product_data)
	fmt.Println("----")
	headers := generateHeaders(product_data[0])
	// create method to convert objects.product to []string
	// have to use the key => value somehow

	data := [][]string{
		headers,
		{"productdata"},
	}
	file, err := os.Create(file_name + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}

// use reflect to return the same
// values from the product object

// Generates a []string containing the export headers
func generateHeaders(product objects.Product) []string {
	headers_products := reflect.TypeOf(objects.Product{})
	headers_variants := reflect.TypeOf(objects.ProductVariant{})
	headers_variant_pricing := getVariantPricingCSV(product, true)
	headers_variant_qty := getVariantQtyCSV(product, true)
	product_headers := make([]string, headers_products.NumField())
	for key := range product_headers {
		product_headers[key] = headers_products.Field(key).Name
	}
	variant_headers := make([]string, headers_variants.NumField())
	for key := range variant_headers {
		variant_headers[key] = headers_variants.Field(key).Name
	}
	product_headers = append(product_headers, generateProductOptions()...)
	product_headers = append(product_headers, variant_headers...)
	product_headers = append(product_headers, headers_variant_pricing...)
	product_headers = append(product_headers, headers_variant_qty...)
	return product_headers
}

func generateReference() []string {
	return []string{
		"active",
		"title",
	}
}

// Create function to extract the product_options per variant option
// option1_name, option1_value etc...
func generateProductOptions() []string {
	return []string{"option1_name", "option1_value", "option2_name",
		"option2_value", "option3_name", "option3_value"}
}

// Returns the name of each warehouse
func getVariantPricingCSV(product objects.Product, key bool) []string {
	qty_headers := []string{}
	if len(product.Variants) > 0 {
		for _, variant := range product.Variants {
			if len(variant.VariantQuantity) > 0 {
				for _, qty := range variant.VariantQuantity {
					if key {
						qty_headers = append(qty_headers, "qty_"+qty.Name)
					} else {
						qty_headers = append(qty_headers, string(qty.Value))
					}
				}
			}
		}
	}
	return qty_headers
}

// Returns the name of each price tier
func getVariantQtyCSV(product objects.Product, key bool) []string {
	pricing_headers := []string{}
	if len(product.Variants) > 0 {
		for _, variant := range product.Variants {
			if len(variant.VariantPricing) > 0 {
				for _, pricing := range variant.VariantPricing {
					if key {
						pricing_headers = append(pricing_headers, "price_"+pricing.Name)
					} else {
						pricing_headers = append(pricing_headers, pricing.Value)
					}
				}
			}
		}
	}
	return pricing_headers
}

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
			qty_header_map = GetKeysByMatcher(value, "qty_")
			price_header_map = GetKeysByMatcher(value, "price_")
			continue
		}
		break
	}
	for key := range records {
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
