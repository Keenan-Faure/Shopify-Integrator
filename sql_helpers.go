package main

import (
	"errors"
	"net/http"
)

// checks if a username already exists inside database
func (dbconfig *DbConfig) CheckUserExist(name string, r *http.Request) (bool, error) {
	username, err := dbconfig.DB.GetUserByName(r.Context(), name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
	}
	if username.Name == name {
		return true, errors.New("name already exists")
	}
	return false, nil
}
