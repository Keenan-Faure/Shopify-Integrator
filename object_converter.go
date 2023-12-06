package main

import (
	"context"
	"fmt"
	"integrator/internal/database"
	"objects"
	"utils"

	"github.com/google/uuid"
)

// Compile Queue Filter Search into a single object (variable)
func CompileQueueFilterSearch(
	dbconfig *DbConfig,
	ctx context.Context,
	page int,
	queue_type,
	status,
	instruction string) ([]objects.ResponseQueueItemFilter, error) {
	response := []objects.ResponseQueueItemFilter{}
	if queue_type == "" {
		if status == "" {
			queue_items, err := dbconfig.DB.GetQueueItemsByInstruction(
				ctx,
				database.GetQueueItemsByInstructionParams{
					Instruction: instruction,
					Limit:       10,
					Offset:      int32((page - 1) * 10),
				})
			if err != nil {
				return []objects.ResponseQueueItemFilter{}, err
			}
			for _, value := range queue_items {
				response = append(response, objects.ResponseQueueItemFilter{
					ID:          value.ID,
					QueueType:   value.QueueType,
					Status:      value.Status,
					Instruction: value.Instruction,
					Object:      value.Object,
					UpdatedAt:   value.UpdatedAt,
				})
			}
			return response, nil
		} else {
			queue_items, err := dbconfig.DB.GetQueueItemsByInstructionAndStatus(
				ctx,
				database.GetQueueItemsByInstructionAndStatusParams{
					Instruction: instruction,
					Status:      status,
					Limit:       10,
					Offset:      int32((page - 1) * 10),
				})
			if err != nil {
				return []objects.ResponseQueueItemFilter{}, err
			}
			for _, value := range queue_items {
				response = append(response, objects.ResponseQueueItemFilter{
					ID:          value.ID,
					QueueType:   value.QueueType,
					Status:      value.Status,
					Instruction: value.Instruction,
					Object:      value.Object,
					UpdatedAt:   value.UpdatedAt,
				})
			}
			return response, nil
		}
	}
	if status == "" {
		if instruction == "" {
			queue_items, err := dbconfig.DB.GetQueueItemsByType(
				ctx,
				database.GetQueueItemsByTypeParams{
					QueueType: queue_type,
					Limit:     10,
					Offset:    int32((page - 1) * 10),
				})
			if err != nil {
				return []objects.ResponseQueueItemFilter{}, err
			}
			for _, value := range queue_items {
				response = append(response, objects.ResponseQueueItemFilter{
					ID:          value.ID,
					QueueType:   value.QueueType,
					Status:      value.Status,
					Instruction: value.Instruction,
					Object:      value.Object,
					UpdatedAt:   value.UpdatedAt,
				})
			}
			return response, nil
		} else {
			queue_items, err := dbconfig.DB.GetQueueItemsByInstructionAndType(
				ctx,
				database.GetQueueItemsByInstructionAndTypeParams{
					Instruction: instruction,
					QueueType:   queue_type,
					Limit:       10,
					Offset:      int32((page - 1) * 10),
				})
			if err != nil {
				return []objects.ResponseQueueItemFilter{}, err
			}
			for _, value := range queue_items {
				response = append(response, objects.ResponseQueueItemFilter{
					ID:          value.ID,
					QueueType:   value.QueueType,
					Status:      value.Status,
					Instruction: value.Instruction,
					Object:      value.Object,
					UpdatedAt:   value.UpdatedAt,
				})
			}
			return response, nil
		}
	}
	if instruction == "" {
		if queue_type == "" {
			queue_items, err := dbconfig.DB.GetQueueItemsByStatus(
				ctx,
				database.GetQueueItemsByStatusParams{
					Status: status,
					Limit:  10,
					Offset: int32((page - 1) * 10),
				})
			if err != nil {
				return []objects.ResponseQueueItemFilter{}, err
			}
			for _, value := range queue_items {
				response = append(response, objects.ResponseQueueItemFilter{
					ID:          value.ID,
					QueueType:   value.QueueType,
					Status:      value.Status,
					Instruction: value.Instruction,
					Object:      value.Object,
					UpdatedAt:   value.UpdatedAt,
				})
			}
			return response, nil
		} else {
			queue_items, err := dbconfig.DB.GetQueueItemsByStatusAndType(
				ctx,
				database.GetQueueItemsByStatusAndTypeParams{
					Status:    status,
					QueueType: queue_type,
					Limit:     10,
					Offset:    int32((page - 1) * 10),
				})
			if err != nil {
				return []objects.ResponseQueueItemFilter{}, err
			}
			for _, value := range queue_items {
				response = append(response, objects.ResponseQueueItemFilter{
					ID:          value.ID,
					QueueType:   value.QueueType,
					Status:      value.Status,
					Instruction: value.Instruction,
					Object:      value.Object,
					UpdatedAt:   value.UpdatedAt,
				})
			}
			return response, nil
		}
	}
	queue_items, err := dbconfig.DB.GetQueueItemsByFilter(
		ctx,
		database.GetQueueItemsByFilterParams{
			Status:      status,
			QueueType:   queue_type,
			Instruction: instruction,
			Limit:       10,
			Offset:      int32((page - 1) * 10),
		})
	if err != nil {
		return []objects.ResponseQueueItemFilter{}, err
	}
	for _, value := range queue_items {
		response = append(response, objects.ResponseQueueItemFilter{
			ID:          value.ID,
			QueueType:   value.QueueType,
			Status:      value.Status,
			Instruction: value.Instruction,
			Object:      value.Object,
			UpdatedAt:   value.UpdatedAt,
		})
	}
	return response, nil
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

// Convert objects.Product into objects.ShopifyProduct
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

// Compile the customer search results
func CompileCustomerSearchData(
	customers_name []database.GetCustomersByNameRow,
	customer_by_id []database.GetCustomersByNameRow) []objects.SearchCustomer {
	customer := []objects.SearchCustomer{}
	for _, value := range customers_name {
		customer = append(customer, objects.SearchCustomer{
			FirstName: value.FirstName,
			LastName:  value.LastName,
		})
	}
	for _, value := range customer_by_id {
		customer = append(customer, objects.SearchCustomer{
			FirstName: value.FirstName,
			LastName:  value.LastName,
		})
	}
	return customer
}

// Compiles the customer data
func CompileCustomerData(
	dbconfig *DbConfig,
	customer_id uuid.UUID,
	ctx context.Context,
	ignore_address bool) (objects.Customer, error) {
	customer, err := dbconfig.DB.GetCustomerByID(ctx, customer_id)
	if err != nil {
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
	customer_address, err := dbconfig.DB.GetAddressByCustomer(ctx, customer_id)
	if err != nil {
		return objects.Customer{}, err
	}
	CustomerAddress := []objects.CustomerAddress{}
	for _, value := range customer_address {
		CustomerAddress = append(CustomerAddress, objects.CustomerAddress{
			FirstName:  value.FirstName,
			LastName:   value.LastName,
			Address1:   value.Address1.String,
			Address2:   value.Address2.String,
			Suburb:     value.Suburb.String,
			City:       value.City.String,
			Province:   value.Province.String,
			PostalCode: value.PostalCode.String,
			Company:    value.Company.String,
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
			WebCode:       value.WebCode.String,
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
			WebCode:       value.WebCode.String,
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
	ctx context.Context,
	ignore_ship_cust bool) (objects.Order, error) {
	order, err := dbconfig.DB.GetOrderByID(ctx, order_id)
	if err != nil {
		return objects.Order{}, err
	}
	customer_id, err := dbconfig.DB.GetCustomerByOrderID(ctx, order_id)
	if err != nil {
		return objects.Order{}, err
	}
	order_customer, err := dbconfig.DB.GetCustomerByID(ctx, customer_id)
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
			WebCode:           order.WebCode.String,
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
	order_customer_address, err := dbconfig.DB.GetAddressByCustomer(ctx, customer_id)
	if err != nil {
		return objects.Order{}, err
	}
	order_line_items, err := dbconfig.DB.GetOrderLinesByOrder(ctx, order_id)
	if err != nil {
		return objects.Order{}, err
	}
	LineItems := []objects.OrderLines{}
	for _, value := range order_line_items {
		LineItems = append(LineItems, objects.OrderLines{
			SKU:      value.Sku,
			Price:    value.Price.String,
			Barcode:  int(value.Barcode.Int32),
			Qty:      int(value.Qty.Int32),
			TaxRate:  value.TaxRate.String,
			TaxTotal: value.TaxTotal.String,
		})
	}
	order_shipping_lines, err := dbconfig.DB.GetShippingLinesByOrder(ctx, order_id)
	if err != nil {
		return objects.Order{}, err
	}
	ShippingLineItems := []objects.OrderLines{}
	for _, value := range order_shipping_lines {
		ShippingLineItems = append(ShippingLineItems, objects.OrderLines{
			SKU:      value.Sku,
			Price:    value.Price.String,
			Barcode:  int(value.Barcode.Int32),
			Qty:      int(value.Qty.Int32),
			TaxRate:  value.TaxRate.String,
			TaxTotal: value.TaxTotal.String,
		})
	}
	OrderCustomerAddress := []objects.CustomerAddress{}
	for _, value := range order_customer_address {
		OrderCustomerAddress = append(OrderCustomerAddress, objects.CustomerAddress{
			Type:       value.Name.String,
			FirstName:  value.FirstName,
			LastName:   value.LastName,
			Address1:   value.Address1.String,
			Address2:   value.Address2.String,
			Suburb:     value.Suburb.String,
			City:       value.City.String,
			Province:   value.Province.String,
			PostalCode: value.PostalCode.String,
			Company:    value.Company.String,
		})
	}
	OrderCustomer.Address = OrderCustomerAddress
	Order := objects.Order{
		ID:                order_id,
		Notes:             order.Notes.String,
		Status:            order.Status,
		WebCode:           order.WebCode.String,
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
func CompileFilterSearch(
	dbconfig *DbConfig,
	ctx context.Context,
	page int,
	product_type,
	category,
	vendor string) ([]objects.SearchProduct, error) {
	response := []objects.SearchProduct{}
	if product_type != "" {
		if category == "" {
			// vendor
			results, err := dbconfig.DB.GetProductsByVendor(ctx, database.GetProductsByVendorParams{
				Vendor: utils.ConvertStringToSQL(vendor),
				Limit:  10,
				Offset: int32((page - 1) * 10),
			})
			if err != nil {
				return response, err
			}
			for _, value := range results {
				response = append(response, objects.SearchProduct{
					ID:          value.ID,
					Title:       value.Title.String,
					Category:    value.Category.String,
					ProductType: value.ProductType.String,
					Vendor:      value.Vendor.String,
				})
			}
			return response, nil
		} else {
			// category & vendor
			results, err := dbconfig.DB.GetProductsByVendorAndCategory(ctx, database.GetProductsByVendorAndCategoryParams{
				Vendor:   utils.ConvertStringToSQL(vendor),
				Category: utils.ConvertStringToSQL(category),
				Limit:    10,
				Offset:   int32((page - 1) * 10),
			})
			if err != nil {
				return response, err
			}
			for _, value := range results {
				response = append(response, objects.SearchProduct{
					ID:          value.ID,
					Title:       value.Title.String,
					Category:    value.Category.String,
					ProductType: value.ProductType.String,
					Vendor:      value.Vendor.String,
				})
			}
			return response, nil
		}
	}
	if category != "" {
		if product_type != "" {
			results, err := dbconfig.DB.GetProductsByVendor(ctx, database.GetProductsByVendorParams{
				Vendor: utils.ConvertStringToSQL(vendor),
				Limit:  10,
				Offset: int32((page - 1) * 10),
			})
			if err != nil {
				return response, err
			}
			for _, value := range results {
				response = append(response, objects.SearchProduct{
					ID:          value.ID,
					Title:       value.Title.String,
					Category:    value.Category.String,
					ProductType: value.ProductType.String,
					Vendor:      value.Vendor.String,
				})
			}
			return response, nil
		}
		results, err := dbconfig.DB.GetProductsByTypeAndVendor(ctx, database.GetProductsByTypeAndVendorParams{
			ProductType: utils.ConvertStringToSQL(product_type),
			Vendor:      utils.ConvertStringToSQL(vendor),
			Limit:       10,
			Offset:      int32((page - 1) * 10),
		})
		if err != nil {
			return response, err
		}
		for _, value := range results {
			response = append(response, objects.SearchProduct{
				ID:          value.ID,
				Title:       value.Title.String,
				Category:    value.Category.String,
				ProductType: value.ProductType.String,
				Vendor:      value.Vendor.String,
			})
		}
		return response, nil
	}
	if vendor != "" {
		if product_type == "" {
			results, err := dbconfig.DB.GetProductsByCategory(ctx, database.GetProductsByCategoryParams{
				Category: utils.ConvertStringToSQL(category),
				Limit:    10,
				Offset:   int32((page - 1) * 10),
			})
			if err != nil {
				return response, err
			}
			for _, value := range results {
				response = append(response, objects.SearchProduct{
					ID:          value.ID,
					Title:       value.Title.String,
					Category:    value.Category.String,
					ProductType: value.ProductType.String,
					Vendor:      value.Vendor.String,
				})
			}
			return response, nil
		} else {
			results, err := dbconfig.DB.GetProductsByTypeAndVendor(ctx, database.GetProductsByTypeAndVendorParams{
				ProductType: utils.ConvertStringToSQL(product_type),
				Vendor:      utils.ConvertStringToSQL(vendor),
				Limit:       10,
				Offset:      int32((page - 1) * 10),
			})
			if err != nil {
				return response, err
			}
			for _, value := range results {
				response = append(response, objects.SearchProduct{
					ID:          value.ID,
					Title:       value.Title.String,
					Category:    value.Category.String,
					ProductType: value.ProductType.String,
					Vendor:      value.Vendor.String,
				})
			}
			return response, nil
		}
	}
	results, err := dbconfig.DB.GetProductsFilter(ctx, database.GetProductsFilterParams{
		Category:    utils.ConvertStringToSQL(category),
		ProductType: utils.ConvertStringToSQL(product_type),
		Vendor:      utils.ConvertStringToSQL(vendor),
		Limit:       10,
		Offset:      int32((page - 1) * 10),
	})
	if err != nil {
		return response, err
	}
	for _, value := range results {
		response = append(response, objects.SearchProduct{
			ID:          value.ID,
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
		})
	}
	return response, nil
}

// Comples the search results into one object
func CompileSearchResult(
	sku []database.GetProductsSearchSKURow,
	title []database.GetProductsSearchTitleRow) []objects.SearchProduct {
	response := []objects.SearchProduct{}
	for _, value := range sku {
		response = append(response, objects.SearchProduct{
			ID:          value.ID,
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
		})
	}
	for _, value := range title {
		response = append(response, objects.SearchProduct{
			ID:          value.ID,
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
		})
	}
	return response
}

// Compiles the product data
func CompileProductData(
	dbconfig *DbConfig,
	product_id uuid.UUID,
	ctx context.Context,
	ignore_variant bool) (objects.Product, error) {
	product, err := dbconfig.DB.GetProductByID(ctx, product_id)
	if err != nil {
		return objects.Product{}, err
	}
	product_options, err := dbconfig.DB.GetProductOptions(ctx, product_id)
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
	product_images, err := dbconfig.DB.GetProductImageByProductID(ctx, product_id)
	if err != nil {
		return objects.Product{}, err
	}
	images := []objects.ProductImages{}
	for _, product_image := range product_images {
		images = append(images, objects.ProductImages{
			Src:       product_image.ImageUrl,
			Position:  int(product_image.Position),
			UpdatedAt: product_image.UpdatedAt,
		})
	}
	if ignore_variant {
		product_data := objects.Product{
			ID:             product_id,
			ProductCode:    product.ProductCode,
			Active:         product.Active,
			Title:          product.Title.String,
			BodyHTML:       product.BodyHtml.String,
			Category:       product.Category.String,
			Vendor:         product.Vendor.String,
			ProductType:    product.ProductType.String,
			Variants:       []objects.ProductVariant{},
			ProductOptions: options,
			ProductImages:  images,
			UpdatedAt:      product.UpdatedAt,
		}
		return product_data, nil
	}
	variant_data, err := CompileVariantsData(dbconfig, product_id, ctx)
	if err != nil {
		return objects.Product{}, err
	}
	product_data := objects.Product{
		ID:             product_id,
		ProductCode:    product.ProductCode,
		Active:         product.Active,
		Title:          product.Title.String,
		BodyHTML:       product.BodyHtml.String,
		Category:       product.Category.String,
		Vendor:         product.Vendor.String,
		ProductType:    product.ProductType.String,
		Variants:       variant_data,
		ProductOptions: options,
		ProductImages:  images,
		UpdatedAt:      product.UpdatedAt,
	}
	return product_data, nil
}

// Compiles all variant data for a product
func CompileVariantsData(
	dbconfig *DbConfig,
	product_id uuid.UUID,
	ctx context.Context) ([]objects.ProductVariant, error) {
	variants, err := dbconfig.DB.GetProductVariants(ctx, product_id)
	if err != nil {
		return []objects.ProductVariant{}, err
	}
	variantsArray := []objects.ProductVariant{}
	for _, value := range variants {
		qty, err := dbconfig.DB.GetVariantQty(ctx, value.ID)
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
		pricing, err := dbconfig.DB.GetVariantPricing(ctx, value.ID)
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

// Compiles a variant data of a single product
func CompileVariantData(
	dbconfig *DbConfig,
	variant_id uuid.UUID,
	ctx context.Context) (objects.ProductVariant, error) {
	variant, err := dbconfig.DB.GetVariantByVariantID(ctx, variant_id)
	if err != nil {
		return objects.ProductVariant{}, err
	}
	variant_data := objects.ProductVariant{}
	qty, err := dbconfig.DB.GetVariantQty(ctx, variant.ID)
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
	pricing, err := dbconfig.DB.GetVariantPricing(ctx, variant.ID)
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
