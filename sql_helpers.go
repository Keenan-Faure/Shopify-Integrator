package main

import (
	"errors"
	"integrator/internal/database"
	"net/http"
	"objects"
	"time"
	"utils"

	"github.com/google/uuid"
)

// checks if a username already exists inside database
func (dbconfig *DbConfig) CheckUserExist(name string, r *http.Request) (bool, error) {
	username, err := dbconfig.DB.GetUserByName(r.Context(), name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
	}
	if username == name {
		return true, errors.New("name already exists")
	}
	return false, nil
}

// Creates an address
func CreateDefaultAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		CustomerID: customer_id,
		Name:       utils.ConvertStringToSQL("default"),
		FirstName:  order_body.Customer.DefaultAddress.FirstName,
		LastName:   order_body.Customer.DefaultAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.FirstName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.City),
		Province:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Province),
		PostalCode: utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Zip),
		Company:    utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Creates an address
func CreateShippingAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		CustomerID: customer_id,
		Name:       utils.ConvertStringToSQL("shipping"),
		FirstName:  order_body.ShippingAddress.FirstName,
		LastName:   order_body.ShippingAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(order_body.ShippingAddress.City),
		Province:   utils.ConvertStringToSQL(order_body.ShippingAddress.Province),
		PostalCode: utils.ConvertStringToSQL(order_body.ShippingAddress.Zip),
		Company:    utils.ConvertStringToSQL(order_body.ShippingAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Creates an address
func CreateBillingAddress(order_body objects.RequestBodyOrder, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		CustomerID: customer_id,
		Name:       utils.ConvertStringToSQL("billing"),
		FirstName:  order_body.BillingAddress.FirstName,
		LastName:   order_body.BillingAddress.LastName,
		Address1:   utils.ConvertStringToSQL(order_body.BillingAddress.FirstName),
		Address2:   utils.ConvertStringToSQL(order_body.BillingAddress.FirstName),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(order_body.BillingAddress.City),
		Province:   utils.ConvertStringToSQL(order_body.BillingAddress.Province),
		PostalCode: utils.ConvertStringToSQL(order_body.BillingAddress.Zip),
		Company:    utils.ConvertStringToSQL(order_body.BillingAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
