package main

import (
	"context"
	"errors"
	"integrator/internal/database"
	"net/http"
	"objects"
	"strings"
	"time"
	"utils"

	"github.com/google/uuid"
)

// Removes all qty for the warehouse
func RemoveGlobalWarehouse(dbconfig *DbConfig, ctx context.Context, warehouse_name string) error {
	_, err := dbconfig.DB.GetVariants(ctx)
	if err != nil {
		return err
	}
	err = dbconfig.DB.RemoveQtyByWarehouseName(ctx, warehouse_name)
	if err != nil {
		return err
	}
	return nil
}

// Inserts new warehouse for all current variations
func InsertGlobalWarehouse(dbconfig *DbConfig, ctx context.Context, warehouse_name string, reindex bool) error {
	variants := []uuid.UUID{}
	// if it should be reindex, then only retrieve the variant ids that doesn't
	// exist in the variant_qty
	if reindex {
		variants_ids, err := dbconfig.DB.GetUnindexedVariants(ctx)
		if err != nil {
			return err
		}
		variants = append(variants, variants_ids...)
	} else {
		variants_ids, err := dbconfig.DB.GetVariants(ctx)
		if err != nil {
			return err
		}
		variants = append(variants, variants_ids...)
	}
	for _, variant := range variants {
		// update ever variant to contain the new warehouse with a default value of 0
		_, err := dbconfig.DB.CreateVariantQty(ctx, database.CreateVariantQtyParams{
			ID:        uuid.New(),
			VariantID: variant,
			Name:      warehouse_name,
			Value:     utils.ConvertIntToSQL(0),
			Isdefault: false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Parses the data and fills in the missing hourly values
// with a 0 value if it does not exist.
func ParseFetchStats(data []database.GetFetchStatsRow) objects.FetchAmountResponse {
	// get the last record (24 hrs back) of time using the first record
	// which should be the latest

	hours := []string{}
	amount := []int64{}
	for _, fsr := range data {
		splited := strings.Split(fsr.Hour, " ")
		if len(splited) > 1 {
			hours = append(hours, splited[1])
		} else {
			hours = append(hours, "00")
		}
		amount = append(amount, fsr.Amount)
	}
	return objects.FetchAmountResponse{
		Amounts: amount,
		Hours:   hours,
	}

	// TODO should I return the missing values as well?
}

// Parses the data and fills in the missing daily values
// with a 0 value if it does not exist.
func ParseOrderStatsNotPaid(data []database.FetchOrderStatsNotPaidRow) objects.OrderAmountResponse {
	// TODO should I return the missing values
	// if it has 2023-12-05 07, but skips 09 should I make it
	days := []string{}
	count := []int64{}
	for _, pos := range data {
		days = append(days, pos.Day)
		count = append(count, pos.Count)
	}
	return objects.OrderAmountResponse{
		Count: count,
		Days:  days,
	}
}

// Parses the data and fills in the missing daily values
// with a 0 value if it does not exist.
func ParseOrderStatsPaid(data []database.FetchOrderStatsPaidRow) objects.OrderAmountResponse {
	// TODO should I return the missing values
	// if it has 2023-12-05 07, but skips 09 should I make it
	days := []string{}
	count := []int64{}
	for _, pos := range data {
		days = append(days, pos.Day)
		count = append(count, pos.Count)
	}
	return objects.OrderAmountResponse{
		Count: count,
		Days:  days,
	}
}

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

// checks if the credentials in the request body refer to a user
func (dbconfig *DbConfig) CheckUserCredentials(
	request_body objects.RequestBodyLogin,
	r *http.Request,
) (database.GetUserCredentialsRow, bool, error) {
	db_user, err := dbconfig.DB.GetUserCredentials(r.Context(), database.GetUserCredentialsParams{
		Name:     request_body.Username,
		Password: request_body.Password,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return database.GetUserCredentialsRow{}, false, err
		}
		return database.GetUserCredentialsRow{}, false, nil
	}
	return db_user, true, nil
}

// checks if a token already exists in the database
func (dbconfig *DbConfig) CheckTokenExists(email string, r *http.Request) (uuid.UUID, bool, error) {
	token, err := dbconfig.DB.GetToken(r.Context(), email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return uuid.UUID{}, false, err
		}
	}
	if token.Email == email {
		return token.Token, true, nil
	}
	return uuid.UUID{}, false, nil
}

// checks if a token already exists in the database
func (dbconfig *DbConfig) CheckUserEmailType(email, user_type string) (bool, error) {
	db_username, err := dbconfig.DB.GetUserByEmailType(context.Background(), database.GetUserByEmailTypeParams{
		Email:    email,
		UserType: user_type,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
	}
	if db_username == email {
		return true, nil
	}
	return false, nil
}

// Creates an address
func CreateDefaultAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		ID:         uuid.New(),
		CustomerID: customer_id,
		Type:       utils.ConvertStringToSQL("default"),
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
		Type:       utils.ConvertStringToSQL("shipping"),
		FirstName:  order_body.ShippingAddress.FirstName,
		LastName:   order_body.ShippingAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.ShippingAddress.LastName),
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
		Type:       utils.ConvertStringToSQL("billing"),
		FirstName:  order_body.BillingAddress.FirstName,
		LastName:   order_body.BillingAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.BillingAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.BillingAddress.LastName),
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

// Creates a map with images in
func CreateImageMap(csv_product objects.CSVProduct) []string {
	images := []string{}
	images = append(images, csv_product.Image1)
	images = append(images, csv_product.Image2)
	images = append(images, csv_product.Image3)
	return images
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
func CheckExistsPriceTier(dbconfig *DbConfig, ctx context.Context, sku, price_tier string, split bool) (bool, error) {
	price_tiers, err := dbconfig.DB.GetVariantPricingBySKU(ctx, sku)
	if err != nil {
		return false, err
	}
	if split {
		price_tier_split := strings.Split(price_tier, "_")
		for _, value := range price_tiers {
			if value.Name == price_tier_split[0] {
				return true, nil
			}
		}
	} else {
		for _, value := range price_tiers {
			if value.Name == price_tier {
				return true, nil
			}
		}
	}
	return false, nil
}

// Checks if a image already exists on a product
func CheckExistsProductImage(dbconfig *DbConfig, ctx context.Context, product_id uuid.UUID, image_url string, position int) (bool, error) {
	images, err := dbconfig.DB.GetProductImageByProductID(ctx, product_id)
	if err != nil {
		return false, err
	}
	for _, image := range images {
		if image.Position == int32(position) {
			if image.ImageUrl == image_url {
				return true, nil
			}
		}
	}
	return false, nil
}

// Checks if a warehouse already exists
// in the database for a certain SKU
func CheckExistsWarehouse(dbconfig *DbConfig, ctx context.Context, sku, warehouse string) (bool, error) {
	// checks if the SKU has the respective warehouse associated to it
	warehouses, err := dbconfig.DB.GetVariantQtyBySKU(ctx, database.GetVariantQtyBySKUParams{
		Sku:  sku,
		Name: warehouse,
	})
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

func IOGetMax(dbconfig *DbConfig, ctx context.Context, column_type string) (int, error) {
	max := 0
	if column_type == "image" {
		max_db, err := dbconfig.DB.GetMaxImagePosition(ctx)
		if err != nil {
			return 0, err
		}
		max = int(max_db)
	} else if column_type == "price" {
		max_db, err := dbconfig.DB.GetCountOfUniquePrices(ctx)
		if err != nil {
			return 0, err
		}
		max = int(max_db)
	} else if column_type == "qty" {
		max_db, err := dbconfig.DB.GetCountOfUniqueWarehouses(ctx)
		if err != nil {
			return 0, err
		}
		max = int(max_db)
	} else {
		return 0, errors.New("invalid column type to retrieve maximum of")
	}
	return max, nil
}

func AddPricingHeaders(dbconfig *DbConfig, ctx context.Context) ([]string, error) {
	price_tiers := []string{}
	price_tiers_db, err := dbconfig.DB.GetUniquePriceTiers(ctx)
	if err != nil {
		return price_tiers, err
	}
	for _, price := range price_tiers_db {
		price_tiers = append(price_tiers, "price_"+price)
	}
	return price_tiers, nil
}

func AddQtyHeaders(dbconfig *DbConfig, ctx context.Context) ([]string, error) {
	warehouses := []string{}
	warehouses_db, err := dbconfig.DB.GetUniqueWarehouses(ctx)
	if err != nil {
		return warehouses, err
	}
	for _, warehouse := range warehouses_db {
		warehouses = append(warehouses, "qty_"+warehouse)
	}
	return warehouses, nil
}
