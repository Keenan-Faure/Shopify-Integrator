package main

import (
	"context"
	"errors"
	"fmt"
	"integrator/internal/database"
	"objects"
	"strings"
	"utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/*
Converts a Shopify Product (from fetch) into a requestBodyProductStruct
*/
func CompileShopifyToSystemProduct(
	shopifyProduct objects.ShopifySingleProduct,
	shopifyVariant objects.ShopifyProductVariant,
	restrictionMap map[string]string,
) objects.RequestBodyProduct {
	// general product values
	requestBody := objects.RequestBodyProduct{}
	requestBody.Title = ApplyFetchRestriction(restrictionMap, shopifyProduct.Title, "title")
	requestBody.BodyHTML = ApplyFetchRestriction(restrictionMap, shopifyProduct.BodyHTML, "body_html")
	requestBody.ProductType = ApplyFetchRestriction(restrictionMap, shopifyProduct.ProductType, "product_type")
	requestBody.Vendor = ApplyFetchRestriction(restrictionMap, shopifyProduct.Vendor, "vendor")

	// product options
	if DeterFetchRestriction(restrictionMap, "options") {
		if len(shopifyProduct.Options) > 0 {
			if shopifyProduct.Options[0].Name != "Title" {
				productOptions := []objects.ProductOptions{}
				for _, option_value := range shopifyProduct.Options {
					productOptions = append(productOptions, objects.ProductOptions{
						Value:    option_value.Name,
						Position: option_value.Position,
					})
				}
				requestBody.ProductOptions = productOptions
			}
		}
	}

	// product variant
	variant := objects.RequestBodyVariant{}
	variant.Sku = shopifyVariant.Sku
	variant.Option1 = ApplyFetchRestriction(restrictionMap, IgnoreDefaultOption(shopifyVariant.Option1), "options")
	variant.Option2 = ApplyFetchRestriction(restrictionMap, IgnoreDefaultOption(shopifyVariant.Option2), "options")
	variant.Option3 = ApplyFetchRestriction(restrictionMap, IgnoreDefaultOption(shopifyVariant.Option3), "options")
	variant.Barcode = ApplyFetchRestriction(restrictionMap, shopifyVariant.Barcode, "barcode")

	// product variant pricing
	// -- Ignored

	// product variant qty
	// -- Ignored

	requestBody.Variants = append(requestBody.Variants, variant)

	return requestBody
}

// Compile Queue Filter Search into a single object (variable)
func CompileRemoveQueueFilter(
	dbconfig *DbConfig,
	mock bool,
	queue_type,
	status,
	instruction string) (string, error) {
	if queue_type == "" && status == "" && instruction == "" {
		return "success", nil
	}
	baseQuery := `DELETE FROM queue_items WHERE `
	queryWhere := ""
	if queue_type != "" {
		queryWhere += "queue_type = '" + queue_type + "' AND "
	}
	if status != "" {
		queryWhere += "status = '" + status + "' AND "
	}
	if instruction != "" {
		queryWhere += "instruction = '" + instruction + "'"
	}
	if queryWhere == "" {
		return "success", nil
	}
	baseQuery += queryWhere
	baseQuery = RemoveQueryKeywords(baseQuery)
	useLocalhost, _ := dbconfig.GetFlagValue(HOST_RUNTIME_FLAG_NAME)
	customConnection, _ := InitCustomConnection(InitConnectionString(useLocalhost, mock))
	_, err := customConnection.Exec(context.Background(), baseQuery)
	if err != nil {
		return "error", err
	}
	return "success", nil
}

// Convert Database.User to objects.ResponseRegister
func ConvertDatabaseToRegister(user database.User) objects.ResponseRegister {
	return objects.ResponseRegister{
		Name:   user.Name,
		Email:  user.Email,
		ApiKey: user.ApiKey,
	}
}

// Convert database.warehouse into warehouses object
func ConvertDatabaseToWarehouse(warehouses []database.GetWarehousesRow) []objects.Warehouse {
	warehouses_object := []objects.Warehouse{}
	for _, warehouse := range warehouses {
		warehouses_object = append(warehouses_object, objects.Warehouse{
			ID:        uuid.New(),
			Name:      warehouse.Name,
			UpdatedAt: warehouse.UpdatedAt,
		})
	}
	return warehouses_object
}

// Compile Queue Filter Search into a single object (variable)
func CompileQueueFilterSearch(
	dbconfig *DbConfig,
	mock bool,
	page int,
	queue_type,
	status,
	instruction string) ([]objects.ResponseQueueItemFilter, error) {
	if status == "" && instruction == "" && queue_type == "" {
		return []objects.ResponseQueueItemFilter{}, nil
	}
	if page == 0 {
		page = 1
	}
	baseQuery := `SELECT id, queue_type, status, instruction, object, updated_at FROM queue_items WHERE `
	queryWhere := ""
	if page == 0 {
		page = 1
	}
	if queue_type != "" {
		queryWhere += "queue_type = '" + queue_type + "' AND "
	}
	if status != "" {
		queryWhere += "status = '" + status + "' AND "
	}
	if instruction != "" {
		queryWhere += "instruction = '" + instruction + "'"
	}
	if queryWhere == "" {
		return []objects.ResponseQueueItemFilter{}, nil
	}
	baseQuery = RemoveQueryKeywords(baseQuery)
	useLocalhost, _ := dbconfig.GetFlagValue(HOST_RUNTIME_FLAG_NAME)
	customConnection, _ := InitCustomConnection(InitConnectionString(useLocalhost, mock))
	rows, _ := customConnection.Query(context.Background(), baseQuery)
	queueItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[objects.ResponseQueueItemFilter])
	if err != nil {
		return []objects.ResponseQueueItemFilter{}, err
	}
	return queueItems, nil
}

// Convert objects.Product into objects.ShopifyProduct
func ConvertProductToShopify(product objects.Product) objects.ShopifyProduct {
	return objects.ShopifyProduct{
		ShopifyProd: objects.ShopifyProd{
			Title:    product.Title,
			BodyHTML: product.BodyHTML,
			Vendor:   product.Vendor,
			Type:     product.ProductType,
			Status:   "active",
			Variants: ConvertVariantToShopifyProdVariant(product),
			Options:  CompileShopifyOptions(product),
		},
	}
}

// Convert objects.Product.Variant into objects.ShopifyProdVariant
func ConvertVariantToShopifyProdVariant(product objects.Product) []objects.ShopifyProdVariant {
	variants := []objects.ShopifyProdVariant{}
	for _, value := range product.Variants {
		variants = append(variants, objects.ShopifyProdVariant{
			Sku:                 value.Sku,
			Price:               "0",
			CompareAtPrice:      "0",
			Option1:             value.Option1,
			Option2:             value.Option2,
			Option3:             value.Option3,
			Barcode:             value.Barcode,
			InventoryManagement: "shopify",
		})
	}
	return variants
}

// Converts a objects.ShopifyProductResponse to a objects.ShopifyPID struct
func ConvertToShopifyIDs(product objects.ShopifyProductResponse) objects.ShopifyIDs {
	ids := objects.ShopifyIDs{}
	ids.ProductID = fmt.Sprint(product.Product.ID)
	variants := []objects.ShopifyVIDs{}
	for _, value := range product.Product.Variants {
		variants = append(variants, objects.ShopifyVIDs{
			VariantID: fmt.Sprint(value.ID),
		})
	}
	ids.Variants = variants
	return ids
}

// Convert objects.Variant into objects.ShopifyVariant
func ConvertVariantToShopify(variant objects.ProductVariant) objects.ShopifyVariant {
	return objects.ShopifyVariant{
		ShopifyVar: objects.ShopifyVar{
			Sku:                 variant.Sku,
			Price:               "0",
			CompareAtPrice:      "0",
			Option1:             variant.Option1,
			Option2:             variant.Option2,
			Option3:             variant.Option3,
			Barcode:             variant.Barcode,
			InventoryManagement: "shopify", // TODO create a product field for this?
		},
	}
}

// Convert objects.Variant into objects.ShopifyVariant
func ConvertVariantToShopifyVariant(variant objects.ProductVariant) objects.ShopifyProdVariant {
	return objects.ShopifyProdVariant{
		ID:                   0,
		ProductID:            0,
		Title:                "",
		Price:                "0",
		Sku:                  variant.Sku,
		Position:             0,
		InventoryPolicy:      "",
		CompareAtPrice:       "0",
		InventoryManagement:  "",
		Option1:              variant.Option1,
		Option2:              variant.Option2,
		Option3:              variant.Option3,
		Barcode:              variant.Barcode,
		Grams:                0,
		InventoryItemID:      0,
		InventoryQuantity:    0,
		OldInventoryQuantity: 0,
	}
}

// Compiles the ShopifyOptions array
func CompileShopifyOptions(product objects.Product) []objects.ShopifyOptions {
	shopify_options := []objects.ShopifyOptions{}
	options_map := CreateOptionMap(product.ProductOptions, product.Variants)
	for key, value := range options_map {
		shopify_options = append(shopify_options, objects.ShopifyOptions{
			Name:   key,
			Values: value,
		})
	}
	return shopify_options
}

// Compiles the customer data
func CompileCustomerData(
	dbconfig *DbConfig,
	customer_id uuid.UUID,
	ignore_address bool) (objects.Customer, error) {
	customer, err := dbconfig.DB.GetCustomerByID(context.Background(), customer_id)
	if err != nil {
		if err.Error() == "sql: no rows found in result set" {
			return objects.Customer{}, errors.New("customer with ID '" + customer_id.String() + "' not found")
		}
		return objects.Customer{}, err
	}
	if ignore_address {
		return objects.Customer{
			ID:        customer_id,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     customer.Email.String,
			Phone:     customer.Phone.String,
			Address:   []objects.CustomerAddress{},
			UpdatedAt: customer.UpdatedAt,
		}, nil
	}
	customer_address, err := dbconfig.DB.GetAddressByCustomer(context.Background(), customer_id)
	if err != nil {
		return objects.Customer{}, err
	}
	CustomerAddress := []objects.CustomerAddress{}
	for _, value := range customer_address {
		CustomerAddress = append(CustomerAddress, objects.CustomerAddress{
			Type:         value.Type,
			FirstName:    value.FirstName,
			LastName:     value.LastName,
			Address1:     value.Address1.String,
			Address2:     value.Address2.String,
			City:         value.City.String,
			Province:     value.Province.String,
			ProvinceCode: value.ProvinceCode.String,
			Company:      value.Company.String,
		})
	}
	return objects.Customer{
		ID:        customer_id,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email.String,
		Phone:     customer.Phone.String,
		Address:   CustomerAddress,
		UpdatedAt: customer.UpdatedAt,
	}, nil
}

// Compiles the order search data
func CompileOrderSearchResult(
	customer_fl []database.GetOrdersSearchByCustomerRow,
	webcode []database.GetOrdersSearchWebCodeRow) []objects.SearchOrder {
	response := []objects.SearchOrder{}
	for _, value := range customer_fl {
		response = append(response, objects.SearchOrder{
			Notes:         value.Notes.String,
			WebCode:       value.WebCode,
			TaxTotal:      value.TaxTotal.String,
			OrderTotal:    value.OrderTotal.String,
			ShippingTotal: value.ShippingTotal.String,
			DiscountTotal: value.DiscountTotal.String,
			UpdatedAt:     value.UpdatedAt,
		})
	}
	for _, value := range webcode {
		response = append(response, objects.SearchOrder{
			Notes:         value.Notes.String,
			WebCode:       value.WebCode,
			TaxTotal:      value.TaxTotal.String,
			OrderTotal:    value.OrderTotal.String,
			ShippingTotal: value.ShippingTotal.String,
			DiscountTotal: value.DiscountTotal.String,
			UpdatedAt:     value.UpdatedAt,
		})
	}
	return response
}

// Compiles the order data
func CompileOrderData(
	dbconfig *DbConfig,
	order_id uuid.UUID,
	ignore_ship_cust bool) (objects.Order, error) {
	order, err := dbconfig.DB.GetOrderByID(context.Background(), order_id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return objects.Order{}, errors.New("order with ID '" + order_id.String() + "' not found")
		}
		return objects.Order{}, err
	}
	customer_id, err := dbconfig.DB.GetCustomerByOrderID(context.Background(), order_id)
	if err != nil {
		return objects.Order{}, err
	}
	order_customer, err := dbconfig.DB.GetCustomerByID(context.Background(), customer_id)
	if err != nil {
		return objects.Order{}, err
	}
	OrderCustomer := objects.OrderCustomer{
		FirstName: order_customer.FirstName,
		LastName:  order_customer.LastName,
		Address:   []objects.CustomerAddress{},
		UpdatedAt: order_customer.UpdatedAt,
	}
	if ignore_ship_cust {
		Order := objects.Order{
			ID:                order.ID,
			Notes:             order.Notes.String,
			Status:            order.Status,
			WebCode:           order.WebCode,
			TaxTotal:          order.TaxTotal.String,
			OrderTotal:        order.OrderTotal.String,
			ShippingTotal:     order.ShippingTotal.String,
			DiscountTotal:     order.DiscountTotal.String,
			UpdatedAt:         order.UpdatedAt,
			CreatedAt:         order.CreatedAt,
			OrderCustomer:     OrderCustomer,
			LineItems:         []objects.OrderLines{},
			ShippingLineItems: []objects.OrderLines{},
		}
		return Order, nil
	}
	order_customer_address, err := dbconfig.DB.GetAddressByCustomer(context.Background(), customer_id)
	if err != nil {
		return objects.Order{}, err
	}
	order_line_items, err := dbconfig.DB.GetOrderLinesByOrder(context.Background(), order_id)
	if err != nil {
		return objects.Order{}, err
	}
	LineItems := []objects.OrderLines{}
	for _, value := range order_line_items {
		LineItems = append(LineItems, objects.OrderLines{
			SKU:      value.Sku,
			Price:    value.Price.String,
			Qty:      int(value.Qty.Int32),
			TaxRate:  value.TaxRate.String,
			TaxTotal: value.TaxTotal.String,
		})
	}
	order_shipping_lines, err := dbconfig.DB.GetShippingLinesByOrder(context.Background(), order_id)
	if err != nil {
		return objects.Order{}, err
	}
	ShippingLineItems := []objects.OrderLines{}
	for _, value := range order_shipping_lines {
		ShippingLineItems = append(ShippingLineItems, objects.OrderLines{
			SKU:      value.Sku,
			Price:    value.Price.String,
			Qty:      int(value.Qty.Int32),
			TaxRate:  value.TaxRate.String,
			TaxTotal: value.TaxTotal.String,
		})
	}
	OrderCustomerAddress := []objects.CustomerAddress{}
	for _, value := range order_customer_address {
		OrderCustomerAddress = append(OrderCustomerAddress, objects.CustomerAddress{
			Type:         value.Type,
			FirstName:    value.FirstName,
			LastName:     value.LastName,
			Address1:     value.Address1.String,
			Address2:     value.Address2.String,
			City:         value.City.String,
			Province:     value.Province.String,
			ProvinceCode: value.ProvinceCode.String,
			Company:      value.Company.String,
		})
	}
	OrderCustomer.Address = OrderCustomerAddress
	Order := objects.Order{
		ID:                order_id,
		Notes:             order.Notes.String,
		Status:            order.Status,
		WebCode:           order.WebCode,
		TaxTotal:          order.TaxTotal.String,
		OrderTotal:        order.OrderTotal.String,
		ShippingTotal:     order.ShippingTotal.String,
		DiscountTotal:     order.DiscountTotal.String,
		UpdatedAt:         order.UpdatedAt,
		CreatedAt:         order.CreatedAt,
		OrderCustomer:     OrderCustomer,
		LineItems:         LineItems,
		ShippingLineItems: ShippingLineItems,
	}
	return Order, nil
}

// Compiles the filter results into one object
func CompileFilterSearch(dbconfig *DbConfig, mock bool, page int, product_type, category, vendor string) ([]objects.SearchProduct, error) {
	baseQuery := `SELECT id, active, product_code, title, category, vendor, product_type, updated_at FROM products WHERE `
	queryWhere := ""
	if page == 0 {
		page = 1
	}
	if product_type != "" {
		queryWhere += "product_type = '" + product_type + "' AND "
	}
	if category != "" {
		queryWhere += "category = '" + category + "' AND "
	}
	if vendor != "" {
		queryWhere += "vendor = '" + vendor + "'"
	}
	if queryWhere == "" {
		return []objects.SearchProduct{}, nil
	}
	baseQuery = RemoveQueryKeywords(baseQuery)
	useLocalhost, _ := dbconfig.GetFlagValue(HOST_RUNTIME_FLAG_NAME)
	customConnection, _ := InitCustomConnection(InitConnectionString(useLocalhost, mock))
	rows, _ := customConnection.Query(context.Background(), baseQuery)
	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[objects.SearchProduct])
	for _, product := range products {
		images, err := CompileProductImages(product.ID, dbconfig)
		if err != nil {
			return []objects.SearchProduct{}, err
		}
		product.Images = images
	}
	if err != nil {
		return []objects.SearchProduct{}, err
	}
	return products, nil
}

func CompileProductImages(
	product_id uuid.UUID,
	dbconfig *DbConfig) ([]objects.ProductImages, error) {
	response := []objects.ProductImages{}
	images, err := dbconfig.DB.GetProductImageByProductID(context.Background(), product_id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return []objects.ProductImages{}, errors.New("product with ID '" + product_id.String() + "' not found")
		}
		return response, err
	}
	for _, image := range images {
		response = append(response, objects.ProductImages{
			Src:       image.ImageUrl,
			Position:  int(image.Position),
			UpdatedAt: image.UpdatedAt,
		})
	}
	return response, nil
}

// Comples the search results into one object
func CompileSearchResult(
	dbconfig *DbConfig,
	search []database.GetProductsSearchRow) ([]objects.SearchProduct, error) {
	response := []objects.SearchProduct{}
	for _, value := range search {
		images, err := CompileProductImages(value.ID, dbconfig)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return []objects.SearchProduct{}, errors.New("product with ID '" + value.ID.String() + "' not found")

			}
			return response, err
		}
		response = append(response, objects.SearchProduct{
			ID:          value.ID,
			Active:      value.Active,
			Images:      images,
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
			UpdatedAt:   value.UpdatedAt,
		})
	}
	return response, nil
}

// Convert Product (POST) into CSVProduct
func ConvertProductToCSVProduct(products objects.RequestBodyProduct) []objects.CSVProduct {
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

// Compiles the product data
func CompileProduct(
	dbconfig *DbConfig,
	product_id uuid.UUID,
	ignore_variant bool) (objects.Product, error) {
	product, err := dbconfig.DB.GetProductByID(context.Background(), product_id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return objects.Product{}, errors.New("product with ID '" + product_id.String() + "' not found")
		}
		return objects.Product{}, err
	}
	product_options, err := dbconfig.DB.GetProductOptions(context.Background(), product_id)
	if err != nil {
		return objects.Product{}, err
	}
	options := []objects.ProductOptions{}
	for _, value := range product_options {
		options = append(options, objects.ProductOptions{
			Value:    value.Name,
			Position: int(value.Position),
		})
	}
	images, err := CompileProductImages(product_id, dbconfig)
	if err != nil {
		return objects.Product{}, err
	}
	if ignore_variant {
		product_data := objects.Product{
			ID:             product_id,
			ProductCode:    strings.ReplaceAll(product.ProductCode, "\"", "'"),
			Active:         product.Active,
			Title:          strings.ReplaceAll(product.Title.String, "\"", "'"),
			BodyHTML:       strings.ReplaceAll(product.BodyHtml.String, "\"", "'"),
			Category:       strings.ReplaceAll(product.Category.String, "\"", "'"),
			Vendor:         strings.ReplaceAll(product.Vendor.String, "\"", "'"),
			ProductType:    strings.ReplaceAll(product.ProductType.String, "\"", "'"),
			Variants:       []objects.ProductVariant{},
			ProductOptions: options,
			ProductImages:  images,
			UpdatedAt:      product.UpdatedAt,
		}
		return product_data, nil
	}
	variant_data, err := CompileVariants(dbconfig, product_id)
	if err != nil {
		return objects.Product{}, err
	}
	product_data := objects.Product{
		ID:             product_id,
		ProductCode:    strings.ReplaceAll(product.ProductCode, "\"", "'"),
		Active:         product.Active,
		Title:          strings.ReplaceAll(product.Title.String, "\"", "'"),
		BodyHTML:       strings.ReplaceAll(product.BodyHtml.String, "\"", "'"),
		Category:       strings.ReplaceAll(product.Category.String, "\"", "'"),
		Vendor:         strings.ReplaceAll(product.Vendor.String, "\"", "'"),
		ProductType:    strings.ReplaceAll(product.ProductType.String, "\"", "'"),
		Variants:       variant_data,
		ProductOptions: options,
		ProductImages:  images,
		UpdatedAt:      product.UpdatedAt,
	}
	return product_data, nil
}

// Compiles all variant data for a product
func CompileVariants(
	dbconfig *DbConfig,
	product_id uuid.UUID,
) ([]objects.ProductVariant, error) {
	variants, err := dbconfig.DB.GetProductVariants(context.Background(), product_id)
	if err != nil {
		return []objects.ProductVariant{}, err
	}
	variantsArray := []objects.ProductVariant{}
	for _, value := range variants {
		qty, err := dbconfig.DB.GetVariantQty(context.Background(), value.ID)
		if err != nil {
			return variantsArray, err
		}
		variant_qty := []objects.VariantQty{}
		for _, sub_value_qty := range qty {
			variant_qty = append(variant_qty, objects.VariantQty{
				IsDefault: sub_value_qty.Isdefault,
				Name:      sub_value_qty.Name,
				Value:     int(sub_value_qty.Value.Int32),
			})
		}
		pricing, err := dbconfig.DB.GetVariantPricing(context.Background(), value.ID)
		if err != nil {
			return variantsArray, err
		}
		variant_pricing := []objects.VariantPrice{}
		for _, sub_value_price := range pricing {
			variant_pricing = append(variant_pricing, objects.VariantPrice{
				IsDefault: sub_value_price.Isdefault,
				Name:      sub_value_price.Name,
				Value:     sub_value_price.Value.String,
			})
		}
		variantsArray = append(variantsArray, objects.ProductVariant{
			ID:              value.ID,
			Sku:             value.Sku,
			Option1:         value.Option1.String,
			Option2:         value.Option2.String,
			Option3:         value.Option3.String,
			Barcode:         value.Barcode.String,
			VariantPricing:  variant_pricing,
			VariantQuantity: variant_qty,
			UpdatedAt:       value.UpdatedAt,
		})
	}
	return variantsArray, nil
}

// Compiles a variant data of a single variant
func CompileVariantByID(
	dbconfig *DbConfig,
	variant_id uuid.UUID,
) (objects.ProductVariant, error) {
	variant, err := dbconfig.DB.GetVariantByVariantID(context.Background(), variant_id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return objects.ProductVariant{}, errors.New("variant with ID '" + variant_id.String() + "' not found")
		}
		return objects.ProductVariant{}, err
	}
	variant_data := objects.ProductVariant{}
	qty, err := dbconfig.DB.GetVariantQty(context.Background(), variant.ID)
	if err != nil {
		return variant_data, err
	}
	variant_qty := []objects.VariantQty{}
	for _, sub_value_qty := range qty {
		variant_qty = append(variant_qty, objects.VariantQty{
			IsDefault: sub_value_qty.Isdefault,
			Name:      sub_value_qty.Name,
			Value:     int(sub_value_qty.Value.Int32),
		})
	}
	pricing, err := dbconfig.DB.GetVariantPricing(context.Background(), variant.ID)
	if err != nil {
		return variant_data, err
	}
	variant_pricing := []objects.VariantPrice{}
	for _, sub_value_price := range pricing {
		variant_pricing = append(variant_pricing, objects.VariantPrice{
			IsDefault: sub_value_price.Isdefault,
			Name:      sub_value_price.Name,
			Value:     sub_value_price.Value.String,
		})
	}
	variant_data = objects.ProductVariant{
		ID:              variant.ID,
		Sku:             variant.Sku,
		Option1:         variant.Option1.String,
		Option2:         variant.Option2.String,
		Option3:         variant.Option3.String,
		Barcode:         variant.Barcode.String,
		VariantPricing:  variant_pricing,
		VariantQuantity: variant_qty,
		UpdatedAt:       variant.UpdatedAt,
	}
	return variant_data, nil
}

// Removes WHERE and AND keywords from baseQuery
func RemoveQueryKeywords(baseQuery string) string {
	if baseQuery == "" {
		return ""
	}
	if baseQuery[len(baseQuery)-5:] == " AND " {
		baseQuery = baseQuery[:len(baseQuery)-5]
	}
	if baseQuery[len(baseQuery)-7:] == " WHERE " {
		baseQuery = baseQuery[:len(baseQuery)-7]
	}
	return baseQuery
}
