package main

import (
	"errors"
	"integrator/internal/database"
	"net/http"
	"objects"
	"strings"
	"time"
	"utils"

	"github.com/google/uuid"
)

// Checks if the VID exists internally.
// Returns an empty string if it doesn't
// and the VID if it does
func (dbconfig *DbConfig) ExistsVID(sku string, r *http.Request) (string, error) {
	pid, err := dbconfig.DB.GetVIDBySKU(r.Context(), sku)
	if err != nil {
		return "", err
	}
	if len(pid.ShopifyVariantID) > 0 && pid.ShopifyVariantID != "" {
		return pid.ShopifyVariantID, nil
	}
	return "", nil
}

// Checks if the PID exists internally
// Returns an empty string if it doesn't
// and the PID if it does
func (dbconfig *DbConfig) ExistsPID(product_code string, r *http.Request) (string, error) {
	pid, err := dbconfig.DB.GetPIDByProductCode(r.Context(), product_code)
	if err != nil {
		return "", err
	}
	if len(pid.ShopifyProductID) > 0 && pid.ShopifyProductID != "" {
		return pid.ShopifyProductID, nil
	}
	return "", nil
}

// checks if a username already exists inside database
func (dbconfig *DbConfig) CheckUserExist(name string, r *http.Request) (bool, error) {
	username, err := dbconfig.DB.GetUserByName(r.Context(), name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
		return false, err
	}
	if username == name {
		return true, errors.New("name already exists")
	}
	return false, nil
}

// checks if a token already exists in the database
func (dbconfig *DbConfig) CheckTokenExists(request_body objects.RequestBodyPreRegister, r *http.Request) (bool, error) {
	token, err := dbconfig.DB.GetToken(r.Context(), database.GetTokenParams{
		Name:  request_body.Name,
		Email: request_body.Email,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
	}
	if token.Email == request_body.Email && token.Name == request_body.Email {
		return true, nil
	}
	return false, nil
}

// Creates an address
func CreateDefaultAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		ID:         uuid.New(),
		CustomerID: customer_id,
		Name:       utils.ConvertStringToSQL("default"),
		FirstName:  order_body.Customer.DefaultAddress.FirstName,
		LastName:   order_body.Customer.DefaultAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.FirstName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.City),
		Province:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Province),
		PostalCode: utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Zip),
		Company:    utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Creates an address
func CreateShippingAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		ID:         uuid.New(),
		CustomerID: customer_id,
		Name:       utils.ConvertStringToSQL("shipping"),
		FirstName:  order_body.ShippingAddress.FirstName,
		LastName:   order_body.ShippingAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(order_body.ShippingAddress.City),
		Province:   utils.ConvertStringToSQL(order_body.ShippingAddress.Province),
		PostalCode: utils.ConvertStringToSQL(order_body.ShippingAddress.Zip),
		Company:    utils.ConvertStringToSQL(order_body.ShippingAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Creates an address
func CreateBillingAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		ID:         uuid.New(),
		CustomerID: customer_id,
		Name:       utils.ConvertStringToSQL("billing"),
		FirstName:  order_body.BillingAddress.FirstName,
		LastName:   order_body.BillingAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.BillingAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.BillingAddress.FirstName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(order_body.BillingAddress.City),
		Province:   utils.ConvertStringToSQL(order_body.BillingAddress.Province),
		PostalCode: utils.ConvertStringToSQL(order_body.BillingAddress.Zip),
		Company:    utils.ConvertStringToSQL(order_body.BillingAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Creates a map of product options vs their names
// map[OptionName][OptionValue]
func CreateOptionMap(
	option_names []objects.ProductOptions,
	variants []objects.ProductVariant) map[string][]string {
	mapp := make(map[string][]string)
	for _, option_name := range option_names {
		for _, variant := range variants {
			if option_name.Position == 1 {
				mapp[option_name.Value] = append(mapp[option_name.Value], variant.Option1)
			} else if option_name.Position == 2 {
				mapp[option_name.Value] = append(mapp[option_name.Value], variant.Option2)
			} else if option_name.Position == 3 {
				mapp[option_name.Value] = append(mapp[option_name.Value], variant.Option3)
			}
			// TODO what happens here?
		}
	}
	return mapp
}

// Create Option Name array
func CreateOptionNamesMap(csv_product objects.CSVProduct) []string {
	mapp := []string{}
	mapp = append(mapp, csv_product.Option1Name)
	mapp = append(mapp, csv_product.Option2Name)
	mapp = append(mapp, csv_product.Option3Name)
	return mapp
}

// Create option Value array
func CreateOptionValuesMap(csv_product objects.CSVProduct) []string {
	mapp := []string{}
	mapp = append(mapp, csv_product.Option1Value)
	mapp = append(mapp, csv_product.Option2Value)
	mapp = append(mapp, csv_product.Option3Value)
	return mapp
}

// Convert Product (POST) into CSVProduct
func ConvertProductToCSV(products objects.RequestBodyProduct) []objects.CSVProduct {
	csv_products := []objects.CSVProduct{}
	for _, variant := range products.Variants {
		csv_product := objects.CSVProduct{
			ProductCode:  products.ProductCode,
			Active:       "1",
			Title:        products.Title,
			BodyHTML:     products.BodyHTML,
			Category:     products.Category,
			Vendor:       products.Vendor,
			ProductType:  products.ProductType,
			SKU:          variant.Sku,
			Option1Value: variant.Option1,
			Option2Value: variant.Option2,
			Option3Value: variant.Option3,
			Barcode:      variant.Barcode,
		}
		if len(products.ProductOptions) == 1 {
			csv_product.Option1Name = utils.IssetString(products.ProductOptions[0].Value)
			if len(products.ProductOptions) == 2 {
				csv_product.Option2Name = utils.IssetString(products.ProductOptions[1].Value)
				if len(products.ProductOptions) == 3 {
					csv_product.Option3Name = utils.IssetString(products.ProductOptions[2].Value)
				}
			}
		}
		csv_products = append(csv_products, csv_product)
	}
	return csv_products
}

// Checks if a price tier already exists
// in the database for a certain SKU
func CheckExistsPriceTier(dbconfig *DbConfig, r *http.Request, sku, price_tier string) (bool, error) {
	price_tiers, err := dbconfig.DB.GetVariantPricingBySKU(r.Context(), sku)
	if err != nil {
		return false, err
	}
	price_tier_split := strings.Split(price_tier, "_")
	for _, value := range price_tiers {
		if value.Name == price_tier_split[0] {
			return true, nil
		}
	}
	return false, nil
}

// Checks if a warehouse already exists
// in the database for a certain SKU
func CheckExistsWarehouse(dbconfig *DbConfig, r *http.Request, sku, warehouse string) (bool, error) {
	warehouses, err := dbconfig.DB.GetVariantQtyBySKU(r.Context(), sku)
	if err != nil {
		return false, err
	}
	warehouse_split := strings.Split(warehouse, "_")
	for _, value := range warehouses {
		if value.Name == warehouse_split[0] {
			return true, nil
		}
	}
	return false, nil
}
