package main

import (
	"context"
	"errors"
	"integrator/internal/database"
	"time"

	"github.com/google/uuid"
)

// Internally saves a flag to the database
// requires a flag name and flag value
// if the value already exists
// then we update - uses upsert
func (dbconfig *DbConfig) AddRuntimeFlags(flagName string, flagValue bool) error {
	err := dbconfig.DB.UpsertRunTimeFlag(context.Background(), database.UpsertRunTimeFlagParams{
		ID:        uuid.New(),
		FlagName:  flagName,
		FlagValue: flagValue,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	return nil
}

// Returns the possible flag value from the database
// in the case that the value does
// not exist it returns false
func (dbconfig *DbConfig) GetFlagValue(flagName string) (bool, error) {
	dbFlagValue, err := dbconfig.DB.GetRuntimeFlag(context.Background(), flagName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, errors.New("runtime flag '" + flagName + "' not found in database")
		}
	}
	return dbFlagValue.FlagValue, nil
}
