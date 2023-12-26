package main

import (
	"integrator/internal/database"
)

// Converts the []database.FetchRestriction slice into a map
func FetchRestrictionsToMap(fetch_map []database.FetchRestriction) map[string]string {
	mapp := make(map[string]string)
	for _, value := range fetch_map {
		mapp[value.Field] = value.Flag
	}
	return mapp
}

// Converts the []database.PushRestriction slice into a map
func PushRestrictionsToMap(push_map []database.PushRestriction) map[string]string {
	mapp := make(map[string]string)
	for _, value := range push_map {
		mapp[value.Field] = value.Flag
	}
	return mapp
}

// Applies the database restriction to the product and returns it
func ApplyFetchRestriction(
	restrictions map[string]string,
	value,
	restriction_type string,
) string {
	// check which restriction is being referenced
	// fetch all restrictions from the database for that type
	// update product to contain restrictions
	// only if the restriction is set to 'app' do we
	// not update the 'app' value with the shopify value

	if map_value, ok := restrictions[restriction_type]; ok {
		if map_value == "app" {
			// do not update
			return ""
		} else {
			return value
		}
	} else {
		// if the key does not exist
		// we should update it
		return value
	}
}

// Applies the database restriction to the product and returns it
func ApplyPushRestriction(
	restrictions map[string]string,
	value,
	restriction_type string,
) string {
	// check which restriction is being referenced
	// fetch all restrictions from the database for that type
	// update product to contain restrictions
	// only if the restriction is set to 'shopify' do we
	// not update the 'shopify' value with the app value

	if map_value, ok := restrictions[restriction_type]; ok {
		if map_value == "shopify" {
			// do not update
			return ""
		} else {
			return value
		}
	} else {
		// if the key does not exist
		// we should update it
		return value
	}
}

// Determines whether a field should be updated or not
// by default if a key does not exist in the database
// it updates the app value with the shopify value during a fetch
func DeterFetchRestriction(
	restriction_map map[string]string,
	restriction_type string,
) bool {
	if map_value, ok := restriction_map[restriction_type]; ok {
		if map_value == "app" {
			// do not update
			return false
		} else {
			return true
		}
	} else {
		// if the key does not exist
		// we should update it
		return true
	}
}

// Determines whether a field should be updated or not
// by default if a key does not exist in the database
// it updates the shopify value with the app value during a push
func DeterPushRestriction(
	restriction_map map[string]string,
	restriction_type string,
) bool {
	if map_value, ok := restriction_map[restriction_type]; ok {
		if map_value == "shopify" {
			// do not update
			return false
		} else {
			return true
		}
	} else {
		// if the key does not exist
		// we should update it
		return true
	}
}
