package main

import (
	"context"
	"database/sql"
	"fmt"
	"integrator/internal/database"
	"log"
	"objects"
	"time"
	"utils"

	"github.com/google/uuid"
)

func (dbconfig *DbConfig) AddOrder(order_body objects.RequestBodyOrder) {
	err := OrderValidation(order_body)
	if err != nil {
		log.Println(err)
		return
	}
	customer, err := dbconfig.DB.CreateCustomer(context.Background(), database.CreateCustomerParams{
		ID:        uuid.New(),
		FirstName: order_body.Customer.FirstName,
		LastName:  order_body.Customer.FirstName,
		Email:     utils.ConvertStringToSQL(order_body.Customer.Email),
		Phone:     utils.ConvertStringToSQL(order_body.Customer.Phone),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), CreateDefaultAddress(order_body, customer.ID))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), CreateShippingAddress(order_body, customer.ID))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), CreateBillingAddress(order_body, customer.ID))
	if err != nil {
		log.Println(err)
		return
	}
	order, err := dbconfig.DB.CreateOrder(context.Background(), database.CreateOrderParams{
		ID:            uuid.New(),
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
		log.Println(err)
		return
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
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprintf("%v", value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
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
				log.Println(err)
				return
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
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprintf("%v", value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
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
				log.Println(err)
				return
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
		log.Println(err)
		return
	}
	// TODO it should return something
	// In the parent thread it should update the queue item to completed if there is no error
}

func (dbconfig *DbConfig) UpdateOrder(order_body objects.RequestBodyOrder) {
	err := OrderValidation(order_body)
	if err != nil {
		log.Println(err)
		return
	}
	db_order, err := dbconfig.DB.GetOrderByWebCode(context.Background(), utils.ConvertStringToSQL(fmt.Sprint(order_body.OrderNumber)))
	if err != nil {
		log.Println(err)
		return
	}
	if db_order.WebCode.String == fmt.Sprint(order_body.OrderNumber) {
		customer, err := dbconfig.DB.GetCustomerByOrderID(context.Background(), db_order.ID)
		if err != nil {
			log.Println(err)
			return
		}
		err = dbconfig.DB.UpdateCustomer(context.Background(), database.UpdateCustomerParams{
			FirstName: order_body.Customer.FirstName,
			LastName:  order_body.Customer.LastName,
			Email:     utils.ConvertStringToSQL(order_body.Customer.Email),
			Phone:     utils.ConvertStringToSQL(order_body.Customer.Phone),
			UpdatedAt: time.Now().UTC(),
			ID:        customer,
		})
		if err != nil {
			log.Println(err)
			return
		}
		err = dbconfig.DB.UpdateAddressByNameAndCustomer(context.Background(), database.UpdateAddressByNameAndCustomerParams{
			CustomerID:   customer,
			FirstName:    order_body.Customer.DefaultAddress.FirstName,
			LastName:     order_body.Customer.DefaultAddress.LastName,
			Suburb:       utils.ConvertStringToSQL(""),
			City:         utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.City),
			Province:     utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Province),
			PostalCode:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Zip),
			Company:      utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Company),
			UpdatedAt:    time.Now().UTC(),
			Name:         utils.ConvertStringToSQL("default"),
			CustomerID_2: customer,
		})
		if err != nil {
			log.Println(err)
			return
		}
		err = dbconfig.DB.UpdateAddressByNameAndCustomer(context.Background(), database.UpdateAddressByNameAndCustomerParams{
			CustomerID:   customer,
			FirstName:    order_body.Customer.DefaultAddress.FirstName,
			LastName:     order_body.Customer.DefaultAddress.LastName,
			Address1:     utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.FirstName),
			Address2:     utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.FirstName),
			Suburb:       utils.ConvertStringToSQL(""),
			City:         utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.City),
			Province:     utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Province),
			PostalCode:   utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Zip),
			Company:      utils.ConvertStringToSQL(order_body.Customer.DefaultAddress.Company),
			UpdatedAt:    time.Now().UTC(),
			Name:         utils.ConvertStringToSQL("billing"),
			CustomerID_2: customer,
		})
		if err != nil {
			log.Println(err)
			return
		}
		err = dbconfig.DB.UpdateAddressByNameAndCustomer(context.Background(), database.UpdateAddressByNameAndCustomerParams{
			CustomerID:   customer,
			FirstName:    order_body.ShippingAddress.FirstName,
			LastName:     order_body.ShippingAddress.LastName,
			Address1:     utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
			Address2:     utils.ConvertStringToSQL(order_body.ShippingAddress.FirstName),
			Suburb:       utils.ConvertStringToSQL(""),
			City:         utils.ConvertStringToSQL(order_body.ShippingAddress.City),
			Province:     utils.ConvertStringToSQL(order_body.ShippingAddress.Province),
			PostalCode:   utils.ConvertStringToSQL(order_body.ShippingAddress.Zip),
			Company:      utils.ConvertStringToSQL(order_body.ShippingAddress.Company),
			UpdatedAt:    time.Now().UTC(),
			Name:         utils.ConvertStringToSQL("shipping"),
			CustomerID_2: customer,
		})
		if err != nil {
			log.Println(err)
			return
		}
		for _, orderline := range order_body.LineItems {
			if len(orderline.TaxLines) > 0 {
				err := dbconfig.DB.UpdateOrderLineByOrderAndSKU(context.Background(), database.UpdateOrderLineByOrderAndSKUParams{
					LineType:  utils.ConvertStringToSQL("product"),
					Sku:       orderline.Sku,
					Price:     utils.ConvertStringToSQL(orderline.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(orderline.Quantity),
					TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(orderline.TaxLines[0].Rate)),
					TaxTotal:  utils.ConvertStringToSQL(orderline.TaxLines[0].Price),
					UpdatedAt: time.Now().UTC(),
					OrderID:   db_order.ID,
					Sku_2:     orderline.Sku,
				})
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				err := dbconfig.DB.UpdateOrderLineByOrderAndSKU(context.Background(), database.UpdateOrderLineByOrderAndSKUParams{
					LineType:  utils.ConvertStringToSQL("product"),
					Sku:       orderline.Sku,
					Price:     utils.ConvertStringToSQL(orderline.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(orderline.Quantity),
					TaxRate:   sql.NullString{},
					TaxTotal:  sql.NullString{},
					UpdatedAt: time.Now().UTC(),
					OrderID:   db_order.ID,
					Sku_2:     orderline.Sku,
				})
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
		for _, shipping_line := range order_body.ShippingLines {
			if len(shipping_line.TaxLines) > 0 {
				err := dbconfig.DB.UpdateOrderLineByOrderAndSKU(context.Background(), database.UpdateOrderLineByOrderAndSKUParams{
					LineType:  utils.ConvertStringToSQL("shipping"),
					Sku:       shipping_line.Code,
					Price:     utils.ConvertStringToSQL(shipping_line.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(1),
					TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(shipping_line.TaxLines[0].Rate)),
					TaxTotal:  utils.ConvertStringToSQL(shipping_line.TaxLines[0].Price),
					UpdatedAt: time.Now().UTC(),
					OrderID:   db_order.ID,
					Sku_2:     shipping_line.Code,
				})
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				err := dbconfig.DB.UpdateOrderLineByOrderAndSKU(context.Background(), database.UpdateOrderLineByOrderAndSKUParams{
					LineType:  utils.ConvertStringToSQL("shipping"),
					Sku:       shipping_line.Code,
					Price:     utils.ConvertStringToSQL(shipping_line.Price),
					Barcode:   utils.ConvertIntToSQL(0),
					Qty:       utils.ConvertIntToSQL(1),
					TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(shipping_line.TaxLines[0].Rate)),
					TaxTotal:  utils.ConvertStringToSQL(shipping_line.TaxLines[0].Price),
					UpdatedAt: time.Now().UTC(),
					OrderID:   db_order.ID,
					Sku_2:     shipping_line.Code,
				})
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
	customer, err := dbconfig.DB.CreateCustomer(context.Background(), database.CreateCustomerParams{
		ID:        uuid.New(),
		FirstName: order_body.Customer.FirstName,
		LastName:  order_body.Customer.FirstName,
		Email:     utils.ConvertStringToSQL(order_body.Customer.Email),
		Phone:     utils.ConvertStringToSQL(order_body.Customer.Phone),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), CreateDefaultAddress(order_body, customer.ID))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), CreateShippingAddress(order_body, customer.ID))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dbconfig.DB.CreateAddress(context.Background(), CreateBillingAddress(order_body, customer.ID))
	if err != nil {
		log.Println(err)
		return
	}
	order, err := dbconfig.DB.CreateOrder(context.Background(), database.CreateOrderParams{
		ID:            uuid.New(),
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
		log.Println(err)
		return
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
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprintf("%v", value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
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
				log.Println(err)
				return
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
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprintf("%v", value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
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
				log.Println(err)
				return
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
		log.Println(err)
		return
	}
	// TODO it should return something
	// In the parent thread it should update the queue item to completed if there is no error
}
