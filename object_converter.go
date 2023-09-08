package main

import (
	"integrator/internal/database"
	"net/http"
	"objects"
	"utils"
)

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
	customer_id []byte,
	r *http.Request) (objects.Customer, error) {
	customer, err := dbconfig.DB.GetCustomerByID(r.Context(), customer_id)
	if err != nil {
		return objects.Customer{}, err
	}
	customer_address, err := dbconfig.DB.GetAddressByCustomer(r.Context(), customer_id)
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
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Address:   CustomerAddress,
		UpdatedAt: customer.UpdatedAt.String(),
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
			UpdatedAt:     value.UpdatedAt.String(),
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
			UpdatedAt:     value.UpdatedAt.String(),
		})
	}
	return response
}

// Compiles the order data
func CompileOrderData(
	dbconfig *DbConfig,
	order_id []byte,
	r *http.Request) (objects.Order, error) {
	order, err := dbconfig.DB.GetOrderByID(r.Context(), order_id)
	if err != nil {
		return objects.Order{}, err
	}
	order_customer, err := dbconfig.DB.GetCustomerByID(r.Context(), order.CustomerID)
	if err != nil {
		return objects.Order{}, err
	}
	order_customer_address, err := dbconfig.DB.GetAddressByCustomer(r.Context(), order.CustomerID)
	if err != nil {
		return objects.Order{}, err
	}
	order_line_items, err := dbconfig.DB.GetOrderLinesByOrder(r.Context(), order_id)
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
	order_shipping_lines, err := dbconfig.DB.GetShippingLinesByOrder(r.Context(), order_id)
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
			FirstName:  value.FirstName,
			LastName:   value.LastName,
			Address1:   value.Address1.String,
			Address2:   value.Address2.String,
			Suburb:     value.Suburb.String,
			Province:   value.Province.String,
			PostalCode: value.PostalCode.String,
			Company:    value.Company.String,
		})
	}
	OrderCustomer := objects.OrderCustomer{
		FirstName: order_customer.FirstName,
		LastName:  order_customer.LastName,
		UpdatedAt: order_customer.UpdatedAt.String(),
		Address:   OrderCustomerAddress,
	}
	Order := objects.Order{
		Notes:             order.Notes.String,
		WebCode:           order.WebCode.String,
		TaxTotal:          order.TaxTotal.String,
		OrderTotal:        order.OrderTotal.String,
		ShippingTotal:     order.ShippingTotal.String,
		DiscountTotal:     order.DiscountTotal.String,
		UpdatedAt:         order.UpdatedAt.String(),
		CreatedAt:         order.CreatedAt.String(),
		OrderCustomer:     OrderCustomer,
		LineItems:         LineItems,
		ShippingLineItems: ShippingLineItems,
	}
	return Order, nil
}

// Compiles the filter results into one object
func CompileFilterSearch(
	dbconfig *DbConfig,
	r *http.Request,
	page int,
	product_type,
	category,
	vendor string) ([]objects.SearchProduct, error) {
	response := []objects.SearchProduct{}
	if product_type != "" {
		prod_type, err := dbconfig.DB.GetProductsByType(r.Context(), database.GetProductsByTypeParams{
			ProductType: utils.ConvertStringToSQL(utils.ConvertStringToLike(product_type)),
			Limit:       10,
			Offset:      int32((page - 1) * 10),
		})
		if err != nil {
			return response, err
		}
		for _, value := range prod_type {
			response = append(response, objects.SearchProduct{
				ID:          string(value.ID),
				Title:       value.Title.String,
				Category:    value.Category.String,
				ProductType: value.ProductType.String,
				Vendor:      value.Vendor.String,
			})
		}
	}
	if category != "" {
		prod_category, err := dbconfig.DB.GetProductsByCategory(r.Context(), database.GetProductsByCategoryParams{
			Category: utils.ConvertStringToSQL(utils.ConvertStringToLike(category)),
			Limit:    10,
			Offset:   int32((page - 1) * 10),
		})
		if err != nil {
			return response, err
		}
		for _, value := range prod_category {
			response = append(response, objects.SearchProduct{
				ID:          string(value.ID),
				Title:       value.Title.String,
				Category:    value.Category.String,
				ProductType: value.ProductType.String,
				Vendor:      value.Vendor.String,
			})
		}
	}
	if vendor != "" {
		prod_vendor, err := dbconfig.DB.GetProductsByVendor(r.Context(), database.GetProductsByVendorParams{
			Vendor: utils.ConvertStringToSQL(utils.ConvertStringToLike(vendor)),
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			return response, err
		}
		for _, value := range prod_vendor {
			response = append(response, objects.SearchProduct{
				ID:          string(value.ID),
				Title:       value.Title.String,
				Category:    value.Category.String,
				ProductType: value.ProductType.String,
				Vendor:      value.Vendor.String,
			})
		}
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
			ID:          string(value.ID),
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
		})
	}
	for _, value := range title {
		response = append(response, objects.SearchProduct{
			ID:          string(value.ID),
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
	product_id []byte,
	r *http.Request) (objects.Product, error) {
	product, err := dbconfig.DB.GetProductByID(r.Context(), product_id)
	if err != nil {
		return objects.Product{}, err
	}
	product_options, err := dbconfig.DB.GetProductOptions(r.Context(), product_id)
	if err != nil {
		return objects.Product{}, err
	}
	options := []objects.ProductOptions{}
	for _, value := range product_options {
		options = append(options, objects.ProductOptions{
			Value: value.Value,
		})
	}
	variants, err := dbconfig.DB.GetProductVariants(r.Context(), product_id)
	variant_data, err := CompileVariantData(dbconfig, variants, r)
	if err != nil {
		return objects.Product{}, err
	}
	product_data := objects.Product{
		Active:         product.Active,
		Title:          product.Title.String,
		BodyHTML:       product.BodyHtml.String,
		Category:       product.Category.String,
		Vendor:         product.Vendor.String,
		ProductType:    product.ProductType.String,
		Variants:       variant_data,
		ProductOptions: options,
		UpdatedAt:      product.UpdatedAt.String(),
	}
	return product_data, nil
}

// Compiles all variant data for a product into a single variable
func CompileVariantData(
	dbconfig *DbConfig,
	variants []database.GetProductVariantsRow,
	r *http.Request) ([]objects.ProductVariant, error) {
	variantsArray := []objects.ProductVariant{}
	for _, value := range variants {
		qty, err := dbconfig.DB.GetVariantQty(r.Context(), value.ID)
		if err != nil {
			return variantsArray, err
		}
		variant_qty := []objects.VariantQty{}
		for _, sub_value_qty := range qty {
			variant_qty = append(variant_qty, objects.VariantQty{
				Name:  sub_value_qty.Name,
				Value: int(sub_value_qty.Value.Int32),
			})
		}
		pricing, err := dbconfig.DB.GetVariantPricing(r.Context(), value.ID)
		if err != nil {
			return variantsArray, err
		}
		variant_pricing := []objects.VariantPrice{}
		for _, sub_value_price := range pricing {
			variant_pricing = append(variant_pricing, objects.VariantPrice{
				Name:  sub_value_price.Name,
				Value: sub_value_price.Value.String,
			})
		}
		variantsArray = append(variantsArray, objects.ProductVariant{
			Sku:             value.Sku,
			Option1:         value.Option1.String,
			Option2:         value.Option2.String,
			Option3:         value.Option3.String,
			Barcode:         value.Barcode.String,
			VariantPricing:  variant_pricing,
			VariantQuantity: variant_qty,
			UpdatedAt:       value.UpdatedAt.String(),
		})
	}
	return variantsArray, nil
}
