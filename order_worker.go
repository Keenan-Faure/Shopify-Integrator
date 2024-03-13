package main

import (
	"context"
	"log"
	"objects"
)

func (dbconfig *DbConfig) AddOrder(order_body objects.RequestBodyOrder) error {
	err := OrderValidation(order_body)
	if err != nil {
		log.Println(err)
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
	exists, err := CheckExistsOrder(dbconfig, context.Background(), order_body.Name)
	if err != nil {
		return err
	}
	orderID, err := dbconfig.DB.GetOrderIDByWebCode(context.Background(), order_body.Name)
	if err != nil {
		return err
	}
	if exists {
		err = UpdateOrder(dbconfig, orderID, order_body)
		if err != nil {
			return err
		}
	}
	return nil
}
