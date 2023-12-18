package iocsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"objects"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"utils"

	"github.com/fatih/structs"
	"github.com/gocarina/gocsv"
)

const csv_remove_time = 5 * time.Minute // 5 minutes
const import_directory = "import"

// Handles the import and upload of the file onto the server
func UploadFile(r *http.Request, relative_directory string) (string, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}
	// FormFile returns the first file for the given key `_import`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("_import")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return "", err
	}
	defer file.Close()

	// Displays properties of file that is uploaded
	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)

	// Only accept text/csv file types
	if handler.Header.Get("Content-Type") != "text/csv" {
		return "", errors.New("only CSV extensions are supported")
	}

	// Make new directory for all imports
	err = os.Mkdir(import_directory, os.FileMode(int(0777)))
	if err != nil {
		if err.Error()[len(err.Error())-11:] != "file exists" {
			return "", err
		}
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := os.CreateTemp(import_directory, "upload-*.csv")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return "", err
	}
	// return that we have successfully uploaded our file!
	return (import_directory + tempFile.Name()), nil
}

func CSVProductHeaders(product objects.Product) []string {
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
	return headers
}

func CSVProductValuesByVariant(product objects.Product, variant objects.ProductVariant, pricing_max, qty_max int) []string {
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
	headers = append(headers, CSVProductVariant(variant)...)
	headers = append(headers, CSVVariantOptions(product, variant)...)
	headers = append(headers, getVariantPricingCSV(variant, pricing_max, false)...)
	headers = append(headers, getVariantQtyCSV(variant, qty_max, false)...)
	headers = append(headers, GetProductImagesCSV(product.ProductImages, 0, false)...)
	return headers
}

// generates the variants
func CSVProductVariant(variant objects.ProductVariant) []string {
	headers := []string{}
	variant_fields := structs.Fields(variant)
	for _, value := range variant_fields {
		if value.Tag("json") == "sku" || value.Tag("json") == "barcode" {
			headers = append(headers, value.Value().(string))
		}
	}
	return headers
}

func CSVVariantOptions(product objects.Product, variant objects.ProductVariant) []string {
	option_values := []string{variant.Option1, variant.Option2, variant.Option3}
	header := []string{}
	counter := 3
	counter = counter - len(product.ProductOptions)
	for option_key := range product.ProductOptions {
		header = append(header, product.ProductOptions[option_key].Value)
		header = append(header, option_values[option_key])
	}
	for {
		if counter < 1 {
			return header
		}
		header = append(header, "")
		header = append(header, "")
		counter = counter - 1
	}
}

// Create function to extract the product_options per variant option
// option1_name, option1_value etc...
func generateProductOptions() []string {
	return []string{"option1_name", "option1_value", "option2_name",
		"option2_value", "option3_name", "option3_value"}
}

// Returns the name/qty of each warehouse depending on the key
func getVariantQtyCSV(variant objects.ProductVariant, qty_max int, key bool) []string {
	qty_headers := []string{}
	for _, qty := range variant.VariantQuantity {
		if qty.Value == 0 {
			qty_headers = append(qty_headers, fmt.Sprintf("%v", 0))
		} else {
			qty_headers = append(qty_headers, fmt.Sprintf("%v", qty.Value))
		}
	}
	qty_max_sub := qty_max - len(variant.VariantQuantity)
	for {
		if qty_max_sub < 1 {
			return qty_headers
		}
		qty_headers = append(qty_headers, "0")
		qty_max_sub = qty_max_sub - 1
	}
}

// Returns the name/value of each price tier depending on the key
func getVariantPricingCSV(variant objects.ProductVariant, pricing_max int, key bool) []string {
	pricing_headers := []string{}
	for _, pricing := range variant.VariantPricing {
		if pricing.Value == "" {
			pricing_headers = append(pricing_headers, "0.00")
		} else {
			pricing_headers = append(pricing_headers, pricing.Value)
		}
	}
	pricing_max_sub := pricing_max - len(variant.VariantPricing)
	for {
		if pricing_max_sub < 1 {
			return pricing_headers
		}
		pricing_headers = append(pricing_headers, "0.00")
		pricing_max_sub = pricing_max_sub - 1
	}
}

// Returns the images of each product
func GetProductImagesCSV(images []objects.ProductImages, max int, key bool) []string {
	image_headers := []string{}
	if key {
		count := 1
		for {
			if count <= max {
				image_headers = append(image_headers, "image_"+fmt.Sprint(count))
				count += 1
				continue
			}
			return image_headers
		}
	} else {
		for _, image := range images {
			image_headers = append(image_headers, fmt.Sprintf("%v", image.Src))
		}
	}
	return image_headers
}

// Writes data to a file
func WriteFile(data [][]string, file_name string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if file_name != "" {
		f, err := os.Create(filepath.Clean(path+"/"+file_name) + ".csv")
		if err != nil {
			return "", err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		err = w.WriteAll(data)

		if err != nil {
			return "", err
		}
		return "", nil
	}
	path = path + "/app/export/"
	err = os.MkdirAll(path, os.FileMode(int(0777)))
	if err != nil {
		if err.Error()[len(err.Error())-11:] != "file exists" {
			return "", err
		}
	}
	csv_name := "product_export-" + time.Now().UTC().String() + ".csv"
	f, err := os.Create(filepath.Clean(path + csv_name))
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.WriteAll(data)
	if err != nil {
		return "", err
	}
	return "http://localhost:8080/app/export/" + csv_name, nil
}

// Reads a csv file contents
func ReadFile(file_name string) ([]objects.CSVProduct, error) {
	if file_name == "" {
		return []objects.CSVProduct{}, errors.New("invalid file")
	}
	file_data, err := os.Open(filepath.Clean(file_name) + ".csv")
	if err != nil {
		return []objects.CSVProduct{}, err
	}
	file_data2, err := os.Open(filepath.Clean(file_name) + ".csv")
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
			Image1:       products[key-1].Image1,
			Image2:       products[key-1].Image2,
			Image3:       products[key-1].Image3,
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

// Returns the keys of all items
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

// loop function that uses Goroutine to run
// a function each interval
func LoopRemoveCSV() {
	ticker := time.NewTicker(csv_remove_time)
	for ; ; <-ticker.C {
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		// removes all exported files
		path_export := path + "/app/export/"
		err = os.MkdirAll(path_export, os.FileMode(int(0777)))
		if err != nil {
			if err.Error()[len(err.Error())-11:] != "file exists" {
				log.Println(err)
			}
		}
		files, err := filepath.Glob(path_export + "product_export*")
		if err != nil {
			log.Println(err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				log.Println(err)
			}
		}

		// removes all imported files
		path_import := path + "/import/"
		err = os.MkdirAll(path_import, os.FileMode(int(0777)))
		if err != nil {
			if err.Error()[len(err.Error())-11:] != "file exists" {
				log.Println(err)
			}
		}
		files, err = filepath.Glob(path_import + "upload-*")
		if err != nil {
			log.Println(err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				log.Println(err)
			}
		}
	}
}
