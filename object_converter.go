package main

import (
	"integrator/internal/database"
	"net/http"
	"objects"
)

// Comples the search results into one object
func CompileSearchResult(
	sku []database.GetProductsSearchSKURow,
	title []database.GetProductsSearchTitleRow) []objects.SearchProduct {
	response := []objects.SearchProduct{}
	for _, value := range sku {
		search_product := objects.SearchProduct{
			ID:          string(value.ID),
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
		}
		response = append(response, search_product)
	}
	for _, value := range title {
		search_product := objects.SearchProduct{
			ID:          string(value.ID),
			Title:       value.Title.String,
			Category:    value.Category.String,
			ProductType: value.ProductType.String,
			Vendor:      value.Vendor.String,
		}
		response = append(response, search_product)
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
		options_sub := objects.ProductOptions{
			Name:  value.Name,
			Value: value.Value,
		}
		options = append(options, options_sub)
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
			sub_variant_qty := objects.VariantQty{
				Name:  sub_value_qty.Name,
				Value: int(sub_value_qty.Value.Int32),
			}
			variant_qty = append(variant_qty, sub_variant_qty)
		}
		pricing, err := dbconfig.DB.GetVariantPricing(r.Context(), value.ID)
		if err != nil {
			return variantsArray, err
		}
		variant_pricing := []objects.VariantPrice{}
		for _, sub_value_price := range pricing {
			sub_variant_price := objects.VariantPrice{
				Name:  sub_value_price.Name,
				Value: sub_value_price.Value.String,
			}
			variant_pricing = append(variant_pricing, sub_variant_price)
		}
		variantData := objects.ProductVariant{
			Sku:             value.Sku,
			Option1:         value.Option1.String,
			Option2:         value.Option2.String,
			Option3:         value.Option3.String,
			Barcode:         value.Barcode.String,
			VariantPricing:  variant_pricing,
			VariantQuantity: variant_qty,
			UpdatedAt:       value.UpdatedAt.String(),
		}
		variantsArray = append(variantsArray, variantData)
	}
	return variantsArray, nil
}
