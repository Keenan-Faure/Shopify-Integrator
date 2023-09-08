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
func CreateAddressUtil(address objects.OrderAddress, customer_id uuid.UUID) database.CreateAddressParams {
	return database.CreateAddressParams{
		CustomerID: customer_id,
		FirstName:  address.Customer.DefaultAddress.FirstName,
		LastName:   address.Customer.DefaultAddress.FirstName,
		Address1:   utils.ConvertStringToSQL(address.Customer.DefaultAddress.Address1),
		Address2:   utils.ConvertStringToSQL(address.Customer.DefaultAddress.Address2),
		Suburb:     utils.ConvertStringToSQL(""),
		City:       utils.ConvertStringToSQL(address.Customer.DefaultAddress.City),
		Province:   utils.ConvertStringToSQL(address.Customer.DefaultAddress.Province),
		PostalCode: utils.ConvertStringToSQL(address.Customer.DefaultAddress.ProvinceCode),
		Company:    utils.ConvertStringToSQL(address.Customer.DefaultAddress.Company),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
