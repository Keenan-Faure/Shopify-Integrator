package main

import (
	"context"
	"errors"
	"objects"

	"github.com/google/uuid"
)

func (dbconfig *DbConfig) AddOrder(order_body objects.RequestBodyOrder) error {
	err := OrderValidation(order_body)
	if err != nil {
		return err
	}
	_, err = AddOrder(dbconfig, order_body)
	if err != nil {
		return err
	}
	return nil
}

func (dbconfig *DbConfig) UpdateOrder(order_body objects.RequestBodyOrder) error {
	err := OrderValidation(order_body)
	if err != nil {
		return err
	}
	orderID, err := dbconfig.DB.GetOrderIDByWebCode(context.Background(), order_body.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("order '" + order_body.Name + "' do not exist")
		}
		return err
	}
	if orderID != uuid.Nil {
		err = UpdateOrder(dbconfig, orderID, order_body)
		if err != nil {
			return err
		}
	}
	return nil
}
