package main

import (
	"context"
	"errors"
	"integrator/internal/database"
	"net/http"
	"objects"
	"strings"
	"time"
	"utils"

	"github.com/google/uuid"
)

/*
This file contains utility functions that attempts to shorten the amount of lines
And keep the code used in the application
Functions are mostly used to interact with the database.
*/

/* Adds an order to the application */
func AddOrder(dbconfig *DbConfig, orderBody objects.RequestBodyOrder) error {
	exists, err := CheckExistsOrder(dbconfig, context.Background(), orderBody.Name)
	if err != nil {
		return err
	}
	if !exists {
		if err := OrderValidation(orderBody); err != nil {
			return err
		}
		return nil
	}
	return nil
}

/* Adds an order's line items to the database under the specific orderID */
func AddOrderLine(orderID uuid.UUID) {

}

/* Adds a customer to the application */
func AddCustomer(dbconfig *DbConfig, customer objects.RequestBodyCustomer, WebCustomerCode string) (uuid.UUID, error) {
	exists, err := CheckExistsCustomer(dbconfig, context.Background(), WebCustomerCode)
	if err != nil {
		return uuid.UUID{}, err
	}
	if !exists {
		dbCustomer, err := dbconfig.DB.CreateCustomer(context.Background(), database.CreateCustomerParams{
			ID:              uuid.New(),
			WebCustomerCode: WebCustomerCode,
			FirstName:       customer.FirstName,
			LastName:        customer.LastName,
			Email:           utils.ConvertStringToSQL(customer.Email),
			Phone:           utils.ConvertStringToSQL(customer.Phone),
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		})
		if err != nil {
			return dbCustomer.ID, err
		}
	}
	return uuid.UUID{}, nil
}

/* Adds a customer address */
func AddCustomerAddress(dbconfig *DbConfig, orderData objects.RequestBodyOrder, customerID uuid.UUID) {
	// Add default, shipping, billing address
}

/* Adds a product to the application */
func AddProduct() {

}

/* Adds a product variant to the application. The productID needs to point to a valid product*/
func AddVariant(productID uuid.UUID) {

}

/*
Returns an array of strings containing the unique price tiers appended with the keyword price_
*/
func AddPricingHeaders(dbconfig *DbConfig, ctx context.Context) ([]string, error) {
	price_tiers := []string{}
	price_tiers_db, err := dbconfig.DB.GetUniquePriceTiers(ctx)
	if err != nil {
		return price_tiers, err
	}
	for _, price := range price_tiers_db {
		price_tiers = append(price_tiers, "price_"+price)
	}
	return price_tiers, nil
}

/*
Returns an array of strings containing the unique warehouses appended with the keyword qty_
*/
func AddQtyHeaders(dbconfig *DbConfig, ctx context.Context) ([]string, error) {
	warehouses := []string{}
	warehouses_db, err := dbconfig.DB.GetUniqueWarehouses(ctx)
	if err != nil {
		return warehouses, err
	}
	for _, warehouse := range warehouses_db {
		warehouses = append(warehouses, "qty_"+warehouse)
	}
	return warehouses, nil
}

/* Inserts new warehouse for all current variations */
func AddGlobalWarehouse(dbconfig *DbConfig, ctx context.Context, warehouse_name string, reindex bool) error {
	variants := []uuid.UUID{}
	// if it should be reindex, then only retrieve the variant ids that doesn't
	// exist in the variant_qty
	if reindex {
		variants_ids, err := dbconfig.DB.GetUnindexedVariants(ctx)
		if err != nil {
			return err
		}
		variants = append(variants, variants_ids...)
	} else {
		variants_ids, err := dbconfig.DB.GetVariants(ctx)
		if err != nil {
			return err
		}
		variants = append(variants, variants_ids...)
	}
	for _, variant := range variants {
		// update ever variant to contain the new warehouse with a default value of 0
		_, err := dbconfig.DB.CreateVariantQty(ctx, database.CreateVariantQtyParams{
			ID:        uuid.New(),
			VariantID: variant,
			Name:      warehouse_name,
			Value:     utils.ConvertIntToSQL(0),
			Isdefault: false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

/* Updates or creates the specific price tier for a specific SKU */
func AddPricing(dbconfig *DbConfig, sku string, variant_id uuid.UUID, pricing_name string, price string) error {
	exists, err := CheckExistsPriceTier(
		dbconfig,
		context.Background(),
		sku,
		pricing_name,
		false,
	)
	if err != nil {
		return err
	}
	if exists {
		err = dbconfig.DB.UpdateVariantPricing(context.Background(), database.UpdateVariantPricingParams{
			Name:      pricing_name,
			Value:     utils.ConvertStringToSQL(price),
			Isdefault: false,
			UpdatedAt: time.Now().UTC(),
			Sku:       sku,
			Name_2:    pricing_name,
		})
		if err != nil {
			return err
		}
	} else {
		_, err = dbconfig.DB.CreateVariantPricing(
			context.Background(),
			database.CreateVariantPricingParams{
				ID:        uuid.New(),
				VariantID: variant_id,
				Name:      pricing_name,
				Value:     utils.ConvertStringToSQL(price),
				Isdefault: false,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

/* Updates or creates an image for a certain product */
func AddImagery(dbconfig *DbConfig, product_id uuid.UUID, image_url string, position int) error {
	exists, err := CheckExistsProductImage(
		dbconfig,
		context.Background(),
		product_id,
		image_url,
		position,
	)
	if err != nil {
		return err
	}
	if exists {
		err = dbconfig.DB.UpdateProductImage(context.Background(), database.UpdateProductImageParams{
			ImageUrl:  image_url,
			UpdatedAt: time.Now().UTC(),
			ProductID: product_id,
			Position:  int32(position),
		})
		if err != nil {
			return err
		}
	} else {
		err = dbconfig.DB.CreateProductImage(context.Background(), database.CreateProductImageParams{
			ID:        uuid.New(),
			ProductID: product_id,
			ImageUrl:  image_url,
			Position:  int32(position),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

/* Updates or creates a warehouse for a certain variant */
func AddWarehouse(dbconfig *DbConfig, sku string, variant_id uuid.UUID, warehouse_name string, qty int) error {
	_, err := dbconfig.DB.GetWarehouseByName(context.Background(), warehouse_name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("warehouse " + warehouse_name + " not found")
		}
		return err
	}
	// if warehouse is found, we update the qty, we cannot create a new one
	err = dbconfig.DB.UpdateVariantQty(context.Background(), database.UpdateVariantQtyParams{
		Name:      warehouse_name,
		Value:     utils.ConvertIntToSQL(qty),
		Isdefault: false,
		Sku:       sku,
		Name_2:    warehouse_name,
	})
	if err != nil {
		return err
	}
	return nil
}

/* Updates or creates product options for a certain product */
func AddProductOptions(dbconfig *DbConfig, product_id uuid.UUID, product_code string, option_names []string) error {
	product_options, err := dbconfig.DB.GetProductOptions(context.Background(), product_id)
	if err != nil {
		return err
	}
	// product does not have any options
	if len(product_options) == 0 {
		for key, option_name := range option_names {
			if option_name != "" {
				_, err := dbconfig.DB.CreateProductOption(context.Background(), database.CreateProductOptionParams{
					ID:        uuid.New(),
					ProductID: product_id,
					Name:      option_name,
					Position:  int32(key + 1),
				})
				if err != nil {
					return err
				}
			}
		}
	} else {
		// product has options, we should update
		for key, option_name := range option_names {
			if option_name != "" {
				_, err := dbconfig.DB.UpdateProductOption(context.Background(), database.UpdateProductOptionParams{
					Name:       option_name,
					Position:   int32(key + 1),
					ProductID:  product_id,
					Position_2: int32(key + 1),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

/* Creates an address */
func AddAddress(dbconfig *DbConfig, address objects.CustomerAddress, customer_id uuid.UUID, address_type string) error {
	_, err := dbconfig.DB.CreateAddress(context.Background(), database.CreateAddressParams{
		ID:         uuid.New(),
		CustomerID: customer_id,
		Type:       utils.ConvertStringToSQL(address_type),
		FirstName:  address.FirstName,
		LastName:   address.LastName,
		Address1:   utils.ConvertStringToSQL(address.FirstName),
		Address2:   utils.ConvertStringToSQL(address.LastName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(address.City),
		Province:   utils.ConvertStringToSQL(address.Province),
		PostalCode: utils.ConvertStringToSQL(address.Zip),
		Company:    utils.ConvertStringToSQL(address.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	return nil
}

/*
Removes the inventory warehouse internally.

Note that doing this will remove all quantity from current products in that warehouse
*/
func RemoveGlobalWarehouse(dbconfig *DbConfig, ctx context.Context, warehouse_name string) error {
	err := dbconfig.DB.RemoveQtyByWarehouseName(ctx, warehouse_name)
	if err != nil {
		return err
	}
	return nil
}

/*
Checks if the VID exists internally.

Returns an empty string if it doesn't and the VID if it does
*/
func (dbconfig *DbConfig) CheckVID(sku string, r *http.Request) (string, error) {
	pid, err := dbconfig.DB.GetVIDBySKU(r.Context(), sku)
	if err != nil {
		return "", err
	}
	if len(pid.ShopifyVariantID) > 0 && pid.ShopifyVariantID != "" {
		return pid.ShopifyVariantID, nil
	}
	return "", nil
}

/*
Checks if the PID exists internally

Returns an empty string if it doesn't and the PID if it does
*/
func (dbconfig *DbConfig) CheckPID(product_code string, r *http.Request) (string, error) {
	pid, err := dbconfig.DB.GetPIDByProductCode(r.Context(), product_code)
	if err != nil {
		return "", err
	}
	if len(pid.ShopifyProductID) > 0 && pid.ShopifyProductID != "" {
		return pid.ShopifyProductID, nil
	}
	return "", nil
}

/* Checks if an order already exists inside the database using it's web code */
func CheckExistsOrder(dbconfig *DbConfig, ctx context.Context, order_web_code string) (bool, error) {
	dbOrder, err := dbconfig.DB.GetOrderByWebCode(ctx, order_web_code)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	if dbOrder.WebCode == order_web_code {
		return true, nil
	}
	return false, nil
}

/* Checks if a customer already exists inside the database using it's customer id on the order payload */
func CheckExistsCustomer(dbconfig *DbConfig, ctx context.Context, customer_id string) (bool, error) {
	dbCustomer, err := dbconfig.DB.GetCustomerByWebCode(ctx, customer_id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	if dbCustomer.WebCustomerCode == customer_id {
		return true, nil
	}
	return false, nil
}

/* Checks if a customer already exists inside the database using it's customer id on the order payload */
func CheckExistsCustomerAddress(dbconfig *DbConfig, ctx context.Context, customerID, addressType string) (bool, error) {
	customerUuid, err := uuid.Parse(customerID)
	if err != nil {
		return false, errors.New("could not decode product id: " + customerID)
	}
	_, err = dbconfig.DB.GetAddressByCustomerAndType(ctx, database.GetAddressByCustomerAndTypeParams{
		CustomerID:  customerUuid,
		AddressType: addressType,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/* Checks if a price tier already exists in the database for a certain SKU */
func CheckExistsPriceTier(dbconfig *DbConfig, ctx context.Context, sku, price_tier string, split bool) (bool, error) {
	price_tiers, err := dbconfig.DB.GetVariantPricingBySKU(ctx, sku)
	if err != nil {
		return false, err
	}
	if split {
		price_tier_split := strings.Split(price_tier, "_")
		for _, value := range price_tiers {
			if value.Name == price_tier_split[0] {
				return true, nil
			}
		}
	} else {
		for _, value := range price_tiers {
			if value.Name == price_tier {
				return true, nil
			}
		}
	}
	return false, nil
}

/* Checks if a image already exists on a product */
func CheckExistsProductImage(dbconfig *DbConfig, ctx context.Context, product_id uuid.UUID, image_url string, position int) (bool, error) {
	images, err := dbconfig.DB.GetProductImageByProductID(ctx, product_id)
	if err != nil {
		return false, err
	}
	for _, image := range images {
		if image.Position == int32(position) {
			if image.ImageUrl == image_url {
				return true, nil
			}
		}
	}
	return false, nil
}

/* Checks if a warehouse already exists in the database for a certain SKU */
func CheckExistsWarehouse(dbconfig *DbConfig, ctx context.Context, sku, warehouse string) (bool, error) {
	// checks if the SKU has the respective warehouse associated to it
	warehouses, err := dbconfig.DB.GetVariantQtyBySKU(ctx, database.GetVariantQtyBySKUParams{
		Sku:  sku,
		Name: warehouse,
	})
	if err != nil {
		return false, err
	}
	warehouse_split := strings.Split(warehouse, "_")
	for _, value := range warehouses {
		if value.Name == warehouse_split[0] {
			return true, nil
		}
	}
	return false, nil
}

/* Checks if a username already exists inside database */
func (dbconfig *DbConfig) CheckUExistsUser(name string, r *http.Request) (bool, error) {
	username, err := dbconfig.DB.GetUserByName(r.Context(), name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
		return false, err
	}
	if username == name {
		return true, errors.New("name already exists")
	}
	return false, nil
}

/* Checks if the credentials in the request body refer to a user */
func (dbconfig *DbConfig) CheckUserCredentials(
	request_body objects.RequestBodyLogin,
	r *http.Request,
) (database.GetUserCredentialsRow, bool, error) {
	db_user, err := dbconfig.DB.GetUserCredentials(r.Context(), database.GetUserCredentialsParams{
		Name:     request_body.Username,
		Password: request_body.Password,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return database.GetUserCredentialsRow{}, false, err
		}
		return database.GetUserCredentialsRow{}, false, nil
	}
	return db_user, true, nil
}

/* Checks if a token already exists in the database */
func (dbconfig *DbConfig) CheckExistsToken(email string, r *http.Request) (uuid.UUID, bool, error) {
	token, err := dbconfig.DB.GetToken(r.Context(), email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return uuid.UUID{}, false, err
		}
	}
	if token.Email == email {
		return token.Token, true, nil
	}
	return uuid.UUID{}, false, nil
}

/* Checks if a token already exists in the database */
func (dbconfig *DbConfig) CheckUserEmailType(email, user_type string) (bool, error) {
	db_username, err := dbconfig.DB.GetUserByEmailType(context.Background(), database.GetUserByEmailTypeParams{
		Email:    email,
		UserType: user_type,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
	}
	if db_username == email {
		return true, nil
	}
	return false, nil
}

/*
Returns the maximum count inside the database for the respective column type.

Note images will count the amount of columns where Prices/Warehouse will count the unique price tiers.
*/
func IOGetMax(dbconfig *DbConfig, ctx context.Context, column_type string) (int, error) {
	max := 0
	if column_type == "image" {
		max_db, err := dbconfig.DB.GetMaxImagePosition(ctx)
		if err != nil {
			return 0, err
		}
		max = int(max_db)
	} else if column_type == "price" {
		max_db, err := dbconfig.DB.GetCountOfUniquePrices(ctx)
		if err != nil {
			return 0, err
		}
		max = int(max_db)
	} else if column_type == "qty" {
		max_db, err := dbconfig.DB.GetCountOfUniqueWarehouses(ctx)
		if err != nil {
			return 0, err
		}
		max = int(max_db)
	} else {
		return 0, errors.New("invalid column type to retrieve maximum of")
	}
	return max, nil
}
