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

	"github.com/fatih/structs"
	"github.com/gocarina/gocsv"
)

func CSVProductHeaders(product objects.Product) {
	headers := []string{}
	product_fields := structs.Fields(&objects.ExportProduct{})
	for _, value := range product_fields {
		headers = append(headers, value.Tag("json"))
	}
	variant_fields := structs.Fields(&objects.ExportVariant{})
	for _, value := range variant_fields {
		headers = append(headers, value.Tag("json"))
	}
	headers = append(headers, generateProductOptions()...)
	if len(product.Variants) > 0 {
		headers = append(headers, getVariantPricingCSV(product.Variants[0], true)...)
		headers = append(headers, getVariantQtyCSV(product.Variants[0], true)...)
	}
	fmt.Println(headers)
}

func CSVProductValuesByVariant(product objects.Product, variant objects.ProductVariant) {
	headers := []string{}
	product_fields := structs.Values(product)
	for _, value := range product_fields {
		if reflect.TypeOf(value).String() == "uuid.UUID" {
			headers = append(headers, fmt.Sprintf("%v", value))
			continue
		}
		if structs.IsStruct(value) || reflect.TypeOf(value).Kind() == reflect.Slice {
			continue
		}
		if reflect.TypeOf(value).String() == "time.Time" {
			headers = append(headers, fmt.Sprintf("%v", value))
			continue
		}
		headers = append(headers, fmt.Sprintf("%v", value))
	}
	headers = append(headers, CSVVariantOptions(product, variant)...)
	headers = append(headers, getVariantPricingCSV(variant, false)...)
	headers = append(headers, getVariantQtyCSV(variant, false)...)
	// add variant qty/pricing
	// done
	fmt.Println(headers)
}

func CSVVariantOptions(product objects.Product, variant objects.ProductVariant) []string {
	header := []string{}
	for key, value := range product.ProductOptions {
		if key == 0 {
			header = append(header, value.Value)
			header = append(header, variant.Option1)
		} else if key == 1 {
			header = append(header, value.Value)
			header = append(header, variant.Option2)
		} else if key == 2 {
			header = append(header, value.Value)
			header = append(header, variant.Option3)
		}
	}
	return header
}

// Create function to extract the product_options per variant option
// option1_name, option1_value etc...
func generateProductOptions() []string {
	return []string{"option1_name", "option1_value", "option2_name",
		"option2_value", "option3_name", "option3_value"}
}

// Returns the name of each warehouse
func getVariantPricingCSV(variant objects.ProductVariant, key bool) []string {
	qty_headers := []string{}
	if len(variant.VariantQuantity) > 0 {
		for _, qty := range variant.VariantQuantity {
			if key {
				qty_headers = append(qty_headers, "qty_"+qty.Name)
			} else {
				qty_headers = append(qty_headers, fmt.Sprintf("%v", qty.Value))
			}
		}
	}
	return qty_headers
}

// Returns the name of each price tier
func getVariantQtyCSV(variant objects.ProductVariant, key bool) []string {
	pricing_headers := []string{}
	if len(variant.VariantPricing) > 0 {
		for _, pricing := range variant.VariantPricing {
			if key {
				pricing_headers = append(pricing_headers, "price_"+pricing.Name)
			} else {
				pricing_headers = append(pricing_headers, pricing.Value)
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
