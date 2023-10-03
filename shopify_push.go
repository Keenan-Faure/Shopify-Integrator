package main

// Pushes a product to Shopify
func (dbconfig *DbConfig) PushProduct() {

	// Check if product ids exist internally
	// If yes, then update (UpdateProductShopify)
	// done

	// If product does not exist
	// Create the product on website (save IDs)
	// Add variants to website as well (save IDs)
	// Add Collection to website
	// done
}

// Pushes a variant to Shopify
func PushVariant() {
	// Check if variant ids exist internally
	// If yes, then update (UpdateVariantShopify)
	// done

	// check if the product exists on the website (getProductBySKU)
	// If it does then retrieve the IDs (and save them)
	// Then update variant using those IDs

	// If no, then create new variant on website under the respective product
	// retrieve the IDs to use in future updates and save id's
}

// Pushes all products in database to Shopify
func Syncronize() {
	// Retrieve all products from database in batches and process them
	// by (loop) products -> (loop) variants

	// TODO errors are logged?
}
