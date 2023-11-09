package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"objects"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

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
		log.Println(err)
		return errors.New("could not decode feed_id: " + queue_item_product.SystemProductID)
	}
	_, err = uuid.Parse(queue_item_product.SystemVariantID)
	if err != nil {
		log.Println(err)
		return errors.New("could not decode feed_id: " + queue_item_product.SystemVariantID)
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
	setting_keys map[string]string) error {
	for _, map_value := range request_settings_map {
		found := false
		if map_value.Key == "" {
			return errors.New("settings key cannot be blank")
		}
		for key := range setting_keys {
			if map_value.Key == strings.ToLower(key) {
				found = true
			}
		}
		if !found {
			return errors.New("setting " + map_value.Key + " not allowed")
		}
	}
	return nil
}

// Validate: ShopifySetting
func SettingValidation(
	shopify_settings_map objects.RequestSettings,
	setting_keys map[string]string) error {
	found := false
	if shopify_settings_map.Key == "" {
		return errors.New("settings key cannot be blank")
	}
	for key := range setting_keys {
		if shopify_settings_map.Key == strings.ToLower(key) {
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
	if location_map.LocationID == "" || location_map.WarehouseName == "" {
		return errors.New("data validation error")
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
func ProductValidationDatabase(csv_product objects.CSVProduct, dbconfig *DbConfig, r *http.Request) error {
	err := ProductSKUValidation(csv_product.SKU, dbconfig, r)
	if err != nil {
		return err
	}
	option_names := CreateOptionNamesMap(csv_product)
	for _, option_name := range option_names {
		err = ProductOptionNameValidation(csv_product.ProductCode, option_name, dbconfig, r)
		if err != nil {
			return err
		}
	}
	option_values := CreateOptionValuesMap(csv_product)
	for _, option_value := range option_values {
		err = ProductOptionNameValidation(csv_product.ProductCode, option_value, dbconfig, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validation: Product | SKU
func ProductSKUValidation(sku string, dbconfig *DbConfig, r *http.Request) error {
	db_sku, err := dbconfig.DB.GetVariantBySKU(r.Context(), sku)
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

// Validation: Product | Option Values
func ProductOptionValueValidation(
	product_code,
	option_value,
	option_name string,
	dbconfig *DbConfig,
	r *http.Request) error {
	if option_value == "" || option_name == "" {
		return nil
	}
	option_names_db, err := dbconfig.DB.GetProductOptionsByCode(r.Context(), product_code)
	if err != nil {
		return err
	}
	option_names := []objects.ProductOptions{}
	for _, option_name := range option_names_db {
		option_names = append(option_names, objects.ProductOptions{
			Value:    option_name.Name,
			Position: int(option_name.Position),
		})
	}
	variants_db, err := dbconfig.DB.GetVariantOptionsByProductCode(r.Context(), product_code)
	if err != nil {
		return err
	}
	variants := []objects.ProductVariant{}
	for _, variant := range variants_db {
		variants = append(variants, objects.ProductVariant{
			Option1: variant.Option1.String,
			Option2: variant.Option2.String,
			Option3: variant.Option3.String,
		})
	}
	mapp := CreateOptionMap(option_names, variants)
	for _, value := range mapp[option_name] {
		if value == option_value {
			return errors.New("duplicate option values not allowed")
		}
	}
	return nil
}

// Validation: Product | Option Names
func ProductOptionNameValidation(
	product_code,
	option_name string,
	dbconfig *DbConfig,
	r *http.Request) error {
	if option_name == "" {
		return nil
	}
	option_names, err := dbconfig.DB.GetProductOptionsByCode(r.Context(), product_code)
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
			log.Println(errors.New("option name already exists, skipping"))
		}
	}
	return nil
}

// ValidateToken: Data validtion
func ValidateTokenValidation(token_request objects.RequestBodyUser) error {
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
func DecodeValidateTokenRequestBody(r *http.Request) (objects.RequestBodyUser, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyUser{}
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
func PreRegisterValidation(preorder objects.RequestBodyPreRegister) error {
	if preorder.Name == "" || len(preorder.Name) == 0 || preorder.Email == "" || len(preorder.Email) == 0 {
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

// Order: data validation
func OrderValidation(order objects.RequestBodyOrder) error {
	if order.Name == "" {
		return errors.New("data validation error")
	}
	return nil
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
func UserValidation(user objects.RequestBodyUser) error {
	if user.Name == "" {
		return errors.New("empty name not allowed")
	}
	return nil
}

// Product: data validation
func ProductValidation(product objects.RequestBodyProduct) error {
	if product.Title == "" {
		return errors.New("empty title not allowed")
	}
	if len(product.Variants) == 0 {
		return errors.New("product must have a SKU")
	}
	if product.Variants[0].Sku == "" {
		return errors.New("empty SKU codes not allowed")
	}
	if product.Variants[0].VariantPricing[0].Name == "" {
		return errors.New("empty price tier name not allowed")
	}
	if product.Variants[0].VariantQuantity[0].Name == "" {
		return errors.New("empty warehouse name not allowed")
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
					return errors.New("duplicate options not allowed: " + value.Value)
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
	dbconfig *DbConfig,
	r *http.Request) error {
	sku_array := []string{}
	for _, value := range product.Variants {
		if slices.Contains(sku_array, value.Sku) {
			return errors.New("duplicate SKUs not allowed: " + value.Sku)
		}
		sku_array = append(sku_array, value.Sku)
	}
	for _, value := range sku_array {
		db_sku, err := dbconfig.DB.GetVariantBySKU(r.Context(), value)
		if err.Error() == "record not found" {
			return nil
		}
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		if err != nil {
			return err
		}
		if db_sku.Sku == value {
			return errors.New("SKU with code " + value + " already exists")
		}
	}
	return nil
}

// Product: Duplicate Option value validation (variations)
func DuplicateOptionValues(product objects.RequestBodyProduct) error {
	if len(product.ProductOptions) == 1 {
		option_values := []string{}
		for _, value := range product.Variants {
			if slices.Contains(option_values, value.Option1) {
				return errors.New("duplicate option value")
			}
			option_values = append(option_values, value.Option1)
		}
	} else if len(product.ProductOptions) == 2 {
		option_1_values := []string{}
		option_2_values := []string{}
		for _, value := range product.Variants {
			option_1_values = append(option_1_values, value.Option1)
			option_2_values = append(option_2_values, value.Option2)
		}
		counter := 0
		for key := range option_1_values {
			for sub_key := range option_2_values {
				if option_2_values[key] == option_2_values[sub_key] && option_1_values[key] == option_1_values[sub_key] {
					counter += 1
				}
				if counter > 1 {
					return errors.New("duplicate option values not allowed")
				}
			}
		}
	} else if len(product.ProductOptions) != 3 {
		return errors.New("too many option values")
	}
	option_1_values := []string{}
	option_2_values := []string{}
	option_3_values := []string{}
	for _, value := range product.Variants {
		option_1_values = append(option_1_values, value.Option1)
		option_2_values = append(option_2_values, value.Option2)
		option_3_values = append(option_3_values, value.Option3)
	}
	counter := 0
	for key := range option_1_values {
		for sub_key := range option_2_values {
			for primal_key := range option_3_values {
				if (option_3_values[key] == option_3_values[primal_key] &&
					option_2_values[key] == option_2_values[sub_key]) &&
					option_1_values[key] == option_1_values[sub_key] {
					counter += 1
				}
				if counter > 1 {
					return errors.New("duplicate option values not allowed")
				}
			}
		}
	}
	return nil
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
func DecodeUserRequestBody(r *http.Request) (objects.RequestBodyUser, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyUser{}
	err := decoder.Decode(&params)
	if err != nil {
		if err.Error() == "" {
			return params, errors.New("invalid request body")
		}
		return params, err
	}
	return params, nil
}
