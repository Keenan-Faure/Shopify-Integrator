package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"objects"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

// Validate: PushRestriction
func RestrictionValidation(
	request_settings_map []objects.RestrictionRequest) error {
	setting_keys := []string{
		"title",
		"body_html",
		"category",
		"vendor",
		"product_type",
		"barcode",
		"options",
		"pricing",
		"warehousing",
	}
	for _, map_value := range request_settings_map {
		found := false
		if map_value.Field == "" {
			return errors.New("settings key cannot be blank")
		}
		for _, value := range setting_keys {
			if map_value.Field == strings.ToLower(value) {
				found = true
			}
		}
		if !found {
			return errors.New("restriction '" + map_value.Field + "' not allowed")
		}
	}
	return nil
}

// Decode: PushRestriction
func DecodeRestriction(dbconfig *DbConfig, r *http.Request) ([]objects.RestrictionRequest, error) {
	decoder := json.NewDecoder(r.Body)
	params := []objects.RestrictionRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return []objects.RestrictionRequest{}, err
	}
	return params, nil
}

// Validate: InsertGlobalWarehouse
func GlobalWarehouseValidation(warehouse objects.RequestGlobalWarehouse) error {
	if warehouse.Name == "" {
		return errors.New("invalid warehouse name")
	}
	return nil
}

// Decode: InsertGlobalWarehouse
func DecodeGlobalWarehouse(dbconfig *DbConfig, r *http.Request) (objects.RequestGlobalWarehouse, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestGlobalWarehouse{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return objects.RequestGlobalWarehouse{}, err
	}
	return params, nil
}

// Decode: QueueItem - Order
func DecodeQueueItemOrder(rawJSON json.RawMessage) (objects.RequestBodyOrder, error) {
	var params objects.RequestBodyOrder
	err := json.Unmarshal(rawJSON, &params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// Decode: QueueItem - Product
func DecodeQueueItemProduct(rawJSON json.RawMessage) (objects.RequestQueueItemProducts, error) {
	var params objects.RequestQueueItemProducts
	err := json.Unmarshal(rawJSON, &params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// Validate: QueueItem - Product
// Validate the UUID's of the object
func QueueItemProductValidation(queue_item_product objects.RequestQueueItemProducts) error {
	_, err := uuid.Parse(queue_item_product.SystemProductID)
	if err != nil {
		return errors.New("could not decode system product id '" + queue_item_product.SystemProductID + "'")
	}
	_, err = uuid.Parse(queue_item_product.SystemVariantID)
	if err != nil {
		return errors.New("could not decode system variant id '" + queue_item_product.SystemVariantID + "'")
	}
	return nil
}

// Validate: QueueItem
func QueueItemValidation(
	request objects.RequestQueueItem) error {
	requests_allowed := []string{
		"product",
		"product_variant",
		"order",
	}
	instructions_allowed := []string{
		"add_product",
		"add_variant",
		"add_order",
		"update_product",
		"update_variant",
		"update_order",
		"zsync_channel",
	}
	if request.Type == "" {
		return errors.New("invalid request type")
	}
	for _, requests := range requests_allowed {
		if requests == request.Type {
			break
		}
	}
	for _, instructions := range instructions_allowed {
		if instructions == request.Instruction {
			break
		}
	}
	return nil
}

// Decode: QueueItem
func DecodeQueueItem(r *http.Request) (objects.RequestQueueItem, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestQueueItem{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return objects.RequestQueueItem{}, err
	}
	return params, nil
}

// Validate: ShopifySettings
func SettingsValidation(
	request_settings_map []objects.RequestSettings,
	setting_keys []string) error {
	for _, map_value := range request_settings_map {
		found := false
		if map_value.Key == "" {
			return errors.New("settings key cannot be blank")
		}
		for _, value := range setting_keys {
			if map_value.Key == strings.ToLower(value) {
				found = true
			}
		}
		if !found {
			return errors.New("setting '" + map_value.Key + "' not allowed")
		}
	}
	return nil
}

// Validate: ShopifySetting
func SettingValidation(
	shopify_settings_map objects.RequestSettings,
	setting_keys []string) error {
	found := false
	if shopify_settings_map.Key == "" {
		return errors.New("settings key cannot be blank")
	}
	for _, value := range setting_keys {
		if shopify_settings_map.Key == strings.ToLower(value) {
			found = true
		}
	}
	if !found {
		return errors.New("setting " + shopify_settings_map.Key + " not allowed")
	}
	return nil
}

// Decode: ShopifySettings
func DecodeSettings(r *http.Request) ([]objects.RequestSettings, error) {
	decoder := json.NewDecoder(r.Body)
	params := []objects.RequestSettings{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return []objects.RequestSettings{}, err
	}
	return params, nil
}

// Decode: ShopifySetting
func DecodeSetting(r *http.Request) (objects.RequestSettings, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestSettings{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return objects.RequestSettings{}, err
	}
	return params, nil
}

// Validation: Inventory Map
func InventoryMapValidation(location_map objects.RequestWarehouseLocation) error {
	if location_map.LocationID == "" || len(location_map.LocationID) == 0 {
		return errors.New("empty location id not allowed")
	}
	if location_map.WarehouseName == "" || len(location_map.WarehouseName) == 0 {
		return errors.New("empty warehouse name not allowed")
	}
	return nil
}

// Decode: Inventory Map
func DecodeInventoryMap(r *http.Request) (objects.RequestWarehouseLocation, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestWarehouseLocation{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// Validation: Product Import
func ProductValidationDatabase(csv_product objects.CSVProduct, dbconfig *DbConfig) error {
	err := ProductSKUValidation(csv_product.SKU, dbconfig)
	if err != nil {
		return err
	}
	option_names := CreateOptionNamesMap(csv_product)
	for _, option_name := range option_names {
		err = ProductOptionNameValidation(csv_product.ProductCode, option_name, dbconfig)
		if err != nil {
			return err
		}
	}
	option_values := CreateOptionValuesMap(csv_product)
	for _, option_value := range option_values {
		err = ProductOptionNameValidation(csv_product.ProductCode, option_value, dbconfig)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validation: Product | SKU
func ProductSKUValidation(sku string, dbconfig *DbConfig) error {
	if sku == "" || len(sku) == 0 {
		return errors.New("invalid sku not allowed")
	}
	db_sku, err := dbconfig.DB.GetVariantBySKU(context.Background(), sku)
	if err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return err
	}
	if db_sku.Sku == sku {
		return errors.New("SKU with code " + sku + " already exists")
	}
	return nil
}

// Validation: Product | Option Names
func ProductOptionNameValidation(
	product_code,
	option_name string,
	dbconfig *DbConfig,
) error {
	if product_code == "" || len(product_code) == 0 {
		return errors.New("invalid product code not allowed")
	}
	if option_name == "" || len(option_name) == 0 {
		return nil
	}
	option_names, err := dbconfig.DB.GetProductOptionsByCode(context.Background(), product_code)
	if err != nil {
		return err
	}
	if len(option_names) > 3 {
		return errors.New("cannot exceed 3 option names")
	}
	for _, value := range option_names {
		if value.Position > 3 || value.Position < 0 {
			return errors.New("invalid option name position for product " + product_code)
		}
		if value.Name == option_name {
			return errors.New("option name already exists")
		}
	}
	return nil
}

// ValidateToken: Data validtion
func ValidateTokenValidation(token_request objects.RequestBodyRegister) error {
	if token_request.Name == "" || len(token_request.Name) == 0 {
		return errors.New("data validation error")
	} else if token_request.Email == "" || len(token_request.Email) == 0 {
		return errors.New("data validation error")
	}
	_, err := uuid.Parse(token_request.Token)
	if err != nil {
		return err
	}
	return nil
}

// ValidateToken: decode the request body
func DecodeValidateTokenRequestBody(r *http.Request) (objects.RequestBodyRegister, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyRegister{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// PreRegister: Data validation
func PreRegisterValidation(prereg objects.RequestBodyPreRegister) error {
	if prereg.Name == "" || len(prereg.Name) == 0 || prereg.Email == "" || len(prereg.Email) == 0 {
		return errors.New("data validation error")
	}
	return nil
}

// PreRegister: decode the request body
func DecodePreRegisterRequestBody(r *http.Request) (objects.RequestBodyPreRegister, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyPreRegister{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// Customer: Data validation
func CustomerValidation(order objects.RequestBodyCustomer) error {
	if order.FirstName == "" {
		return errors.New("data validation error")
	}
	return nil
}

// Customer: decode the request body
func DecodeCustomerRequestBody(r *http.Request) (objects.RequestBodyCustomer, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyCustomer{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// Order: data validation
func OrderValidation(order objects.RequestBodyOrder) error {
	if order.Name == "" {
		return errors.New("data validation error")
	}
	return nil
}

// Order: decodes the request body
func DecodeOrderRequestBody(r *http.Request) (objects.RequestBodyOrder, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyOrder{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// User: data validation
func TokenValidation(key string) error {
	if key == "" || len(key) <= 0 || len(key) > 64 {
		return errors.New("invalid product id")
	}
	return nil
}

// Product: data validation
func IDValidation(id string) error {
	if id == "" || len(id) <= 0 || len(id) > 36 {
		return errors.New("invalid product id")
	}
	return nil
}

// User: data validation
func UserValidation(username, password string) error {
	if username == "" || len(username) == 0 {
		return errors.New("empty username not allowed")
	}
	if password == "" || len(password) == 0 {
		return errors.New("empty password not allowed")
	}
	return nil
}

// Product: data validation
func ProductValidation(dbconfig *DbConfig, product objects.RequestBodyProduct) error {
	if product.Title == "" {
		return errors.New("empty title not allowed")
	}
	if len(product.Variants) == 0 {
		return errors.New("product must have a SKU")
	}
	if product.Variants[0].Sku == "" {
		return errors.New("empty SKU codes not allowed")
	}
	if len(product.Variants[0].VariantPricing) > 0 {
		if product.Variants[0].VariantPricing[0].Name == "" {
			return errors.New("empty price tier name not allowed")
		}
	} else {
		return errors.New("product must have a price")
	}
	for _, variant := range product.Variants {
		for key_qty := range variant.VariantQuantity {
			// check if the warehouse exists, then we update the quantity
			warehouse_name := variant.VariantQuantity[key_qty].Name
			_, err := dbconfig.DB.GetWarehouseByName(context.Background(), warehouse_name)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					return errors.New("warehouse " + warehouse_name + " not found")
				}
				return err
			}
		}
	}
	return nil
}

// Product: Duplicate Option validation
func ValidateDuplicateOption(product objects.RequestBodyProduct) error {
	options_names := []string{}
	if len(product.ProductOptions) > 1 {
		for _, value := range product.ProductOptions {
			if value.Value != "" && len(value.Value) > 0 {
				if slices.Contains(options_names, value.Value) {
					return errors.New("duplicate product option names not allowed: " + value.Value)
				}
				options_names = append(options_names, value.Value)
			}
		}
	}
	return nil
}

// Product: Duplicate SKU validation
func ValidateDuplicateSKU(
	product objects.RequestBodyProduct,
	dbconfig *DbConfig) error {
	sku_array := []string{}
	for _, value := range product.Variants {
		if slices.Contains(sku_array, value.Sku) {
			return errors.New("duplicate SKUs not allowed: " + value.Sku)
		}
		sku_array = append(sku_array, value.Sku)
	}
	for _, value := range sku_array {
		db_sku, err := dbconfig.DB.GetVariantBySKU(context.Background(), value)
		if err != nil {
			if err.Error() == "record not found" {
				return nil
			}
			if err.Error() == "sql: no rows in result set" {
				return nil
			}
			return err
		}
		if db_sku.Sku == value {
			return errors.New("SKU with code " + value + " already exists")
		}
	}
	return nil
}

// Product: Duplicate Option value validation (variations)
func DuplicateOptionValues(dbconfig *DbConfig, variantData objects.RequestBodyVariant, productID uuid.UUID) error {
	products, err := CompileProduct(dbconfig, productID, false)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}
	for _, variant := range products.Variants {
		duplicatedOptions := 0
		requestVariantOptions := CreateProductOptionSlice(variantData.Option1, variantData.Option2, variantData.Option3)
		variantOptions := CreateProductOptionSlice(variant.Option1, variant.Option2, variant.Option3)
		requestOptionsLen := fmt.Sprint(len(requestVariantOptions))
		variantOptionsLen := fmt.Sprint(len(variantOptions))

		for key := range requestVariantOptions {
			if requestOptionsLen == variantOptionsLen {
				if requestVariantOptions[key] == variantOptions[key] {
					duplicatedOptions++
				}
			} else {
				return errors.New("invalid variant option amount, expected '" + variantOptionsLen + "' but found '" + requestOptionsLen + "'")
			}
		}
		if duplicatedOptions >= len(variantOptions) {
			return errors.New("duplicate option values not allowed")
		}
	}
	return nil
}

// Product: Creates a slice containing valid strings of product option values
func CreateProductOptionSlice(option1, option2, option3 string) []string {
	options := []string{}
	if len(option1) > 0 && option1 != "" {
		options = append(options, option1)
		if len(option2) > 0 && option2 != "" {
			options = append(options, option2)
			if len(option3) > 0 && option3 != "" {
				options = append(options, option3)
			}
		}
	}
	return options
}

// Product: decodes the request body
func DecodeProductRequestBody(r *http.Request) (objects.RequestBodyProduct, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyProduct{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// User: decodes the request body
func DecodeUserRequestBody(r *http.Request) (objects.RequestBodyRegister, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyRegister{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}

// Login: decodes the request body
func DecodeLoginRequestBody(r *http.Request) (objects.RequestBodyLogin, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyLogin{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}
