package main

import "objects"

// Applies the database restriction to the product and returns it
func ApplyRestriction(
	dbconfig *DbConfig,
	product objects.Product,
	restriction_type string,
) (objects.Product, error) {
	// check which restriction is being referenced
	// fetch all restrictions from the database for that type
	// update product to contain restrictions
	return objects.Product{}, nil
}
