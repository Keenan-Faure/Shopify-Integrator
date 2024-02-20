package main

import (
	"integrator/internal/database"
	"objects"
	"strings"
)

/*
This file contains various functions that act as utilities when adding, returning, responding with data.
*/

/* Parses the data and fills in the missing hourly values with a 0 value if it does not exist. */
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

/* Parses the data and fills in the missing daily values with a 0 value if it does not exist. */
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

/* Parses the data and fills in the missing daily values with a 0 value if it does not exist. */
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

/* Creates a map of product options vs their names map[OptionName][OptionValue] */
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

/* Create Option Name array */
func CreateOptionNamesMap(csv_product objects.AppProduct) []string {
	mapp := []string{}
	mapp = append(mapp, csv_product.Option1Name)
	mapp = append(mapp, csv_product.Option2Name)
	mapp = append(mapp, csv_product.Option3Name)
	return mapp
}

/* Create option Value array */
func CreateOptionValuesMap(csv_product objects.AppProduct) []string {
	mapp := []string{}
	mapp = append(mapp, csv_product.Option1Value)
	mapp = append(mapp, csv_product.Option2Value)
	mapp = append(mapp, csv_product.Option3Value)
	return mapp
}

/* Creates an array map with images  */
func CreateImageMap(csv_product objects.AppProduct) []string {
	images := []string{}
	images = append(images, csv_product.Image1)
	images = append(images, csv_product.Image2)
	images = append(images, csv_product.Image3)
	return images
}
