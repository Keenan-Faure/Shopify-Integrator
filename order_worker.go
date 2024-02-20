package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"integrator/internal/database"
	"log"
	"objects"
	"time"
	"utils"

	"github.com/google/uuid"
)

func (dbconfig *DbConfig) AddOrder(order_body objects.RequestBodyOrder) error {
	err := OrderValidation(order_body)
	if err != nil {
		log.Println(err)
		return err
	}
	customer, err := dbconfig.DB.CreateCustomer(context.Background(), database.CreateCustomerParams{
		ID:        uuid.New(),
		FirstName: order_body.Customer.FirstName,
		LastName:  order_body.Customer.LastName,
		Email:     utils.ConvertStringToSQL(order_body.Customer.Email),
		Phone:     utils.ConvertStringToSQL(order_body.Customer.Phone),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), AddAddress(order_body, customer.ID, "default"))
	if err != nil {
		return err
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), AddAddress(order_body, customer.ID, "shipping"))
	if err != nil {
		return err
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), AddAddress(order_body, customer.ID, "billing"))
	if err != nil {
		return err
	}
	order, err := dbconfig.DB.CreateOrder(context.Background(), database.CreateOrderParams{
		ID:            uuid.New(),
		Status:        order_body.FinancialStatus,
		Notes:         utils.ConvertStringToSQL(""),
		WebCode:       utils.ConvertStringToSQL(order_body.Name),
		TaxTotal:      utils.ConvertStringToSQL(order_body.TotalTax),
		OrderTotal:    utils.ConvertStringToSQL(order_body.TotalPrice),
		ShippingTotal: utils.ConvertStringToSQL(order_body.TotalShippingPriceSet.ShopMoney.Amount),
		DiscountTotal: utils.ConvertStringToSQL(order_body.TotalDiscounts),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	for _, value := range order_body.LineItems {
		if len(value.TaxLines) > 0 {
			_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("product"),
				Sku:       value.Sku,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(value.Quantity),
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				return err
			}
		} else {
			_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("product"),
				Sku:       value.Sku,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(value.Quantity),
				TaxRate:   sql.NullString{},
				TaxTotal:  sql.NullString{},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				return err
			}
		}
	}
	for _, value := range order_body.ShippingLines {
		if len(value.TaxLines) > 0 {
			_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("shipping"),
				Sku:       value.Code,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(1),
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				return err
			}
		} else {
			_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("shipping"),
				Sku:       value.Code,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(1),
				TaxRate:   sql.NullString{},
				TaxTotal:  sql.NullString{},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				return err
			}
		}
	}
	err = dbconfig.DB.CreateCustomerOrder(context.Background(), database.CreateCustomerOrderParams{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		OrderID:    order.ID,
		UpdatedAt:  time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (dbconfig *DbConfig) UpdateOrder(order_body objects.RequestBodyOrder) error {
	err := OrderValidation(order_body)
	if err != nil {
		return nil
	}
	db_order, err := dbconfig.DB.GetOrderByWebCode(context.Background(), utils.ConvertStringToSQL(fmt.Sprint(order_body.Name)))
	if err != nil {
		return nil
	}
	if db_order.WebCode.String == fmt.Sprint(order_body.Name) {
		// delete order
		err := dbconfig.DB.RemoveOrder(context.Background(), db_order.ID)
		if err != nil {
			return err
		}
		customer, err := dbconfig.DB.CreateCustomer(context.Background(), database.CreateCustomerParams{
			ID:        uuid.New(),
			FirstName: order_body.Customer.FirstName,
			LastName:  order_body.Customer.LastName,
			Email:     utils.ConvertStringToSQL(order_body.Customer.Email),
			Phone:     utils.ConvertStringToSQL(order_body.Customer.Phone),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		_, err = dbconfig.DB.CreateAddress(context.Background(), AddAddress(order_body, customer.ID, "default"))
		if err != nil {
			return err
		}
		_, err = dbconfig.DB.CreateAddress(context.Background(), AddAddress(order_body, customer.ID, "shipping"))
		if err != nil {
			return err
		}
		_, err = dbconfig.DB.CreateAddress(context.Background(), AddAddress(order_body, customer.ID, "billing"))
		if err != nil {
			return err
		}
		order, err := dbconfig.DB.CreateOrder(context.Background(), database.CreateOrderParams{
			ID:            db_order.ID,
			Status:        order_body.FinancialStatus,
			Notes:         utils.ConvertStringToSQL(""),
			WebCode:       utils.ConvertStringToSQL(order_body.Name),
			TaxTotal:      utils.ConvertStringToSQL(order_body.TotalTax),
			OrderTotal:    utils.ConvertStringToSQL(order_body.TotalPrice),
			ShippingTotal: utils.ConvertStringToSQL(order_body.TotalShippingPriceSet.ShopMoney.Amount),
			DiscountTotal: utils.ConvertStringToSQL(order_body.TotalDiscounts),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		for _, value := range order_body.LineItems {
			if len(value.TaxLines) > 0 {
				_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
					ID:        uuid.New(),
					OrderID:   order.ID,
					LineType:  utils.ConvertStringToSQL("product"),
					Sku:       value.Sku,
					Price:     utils.ConvertStringToSQL(value.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(value.Quantity),
					TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(value.TaxLines[0].Rate)),
					TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				})
				if err != nil {
					return err
				}
			} else {
				_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
					ID:        uuid.New(),
					OrderID:   order.ID,
					LineType:  utils.ConvertStringToSQL("product"),
					Sku:       value.Sku,
					Price:     utils.ConvertStringToSQL(value.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(value.Quantity),
					TaxRate:   sql.NullString{},
					TaxTotal:  sql.NullString{},
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				})
				if err != nil {
					return err
				}
			}
		}
		for _, value := range order_body.ShippingLines {
			if len(value.TaxLines) > 0 {
				_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
					ID:        uuid.New(),
					OrderID:   order.ID,
					LineType:  utils.ConvertStringToSQL("shipping"),
					Sku:       value.Code,
					Price:     utils.ConvertStringToSQL(value.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(1),
					TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(value.TaxLines[0].Rate)),
					TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				})
				if err != nil {
					return err
				}
			} else {
				_, err := dbconfig.DB.CreateOrderLine(context.Background(), database.CreateOrderLineParams{
					ID:        uuid.New(),
					OrderID:   order.ID,
					LineType:  utils.ConvertStringToSQL("shipping"),
					Sku:       value.Code,
					Price:     utils.ConvertStringToSQL(value.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(1),
					TaxRate:   sql.NullString{},
					TaxTotal:  sql.NullString{},
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				})
				if err != nil {
					return err
				}
			}
		}
		err = dbconfig.DB.CreateCustomerOrder(context.Background(), database.CreateCustomerOrderParams{
			ID:         uuid.New(),
			CustomerID: customer.ID,
			OrderID:    order.ID,
			UpdatedAt:  time.Now().UTC(),
			CreatedAt:  time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		return nil
	} else {
		// the Order web code is not found
		return errors.New("could not find order to update with code " + db_order.WebCode.String)
	}
}
