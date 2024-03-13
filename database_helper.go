package main

import (
	"context"
	"errors"
	"fmt"
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

/* Upserts a Price */
func UpsertPrice(
	dbconfig *DbConfig,
	currentRecord objects.ImportResponse,
	product objects.CSVProduct,
) objects.ImportResponse {
	variantID, exists, err := QueryVariantIDBySKU(dbconfig, product.SKU)
	if err != nil {
		currentRecord.FailCounter++
		return currentRecord
	}
	if exists {
		for _, price := range product.Pricing {
			err := AddPricing(dbconfig, product.SKU, variantID, price.Name, price.Value)
			if err != nil {
				currentRecord.FailCounter++
				continue
			}
		}
		return currentRecord
	} else {
		currentRecord.FailCounter++
		return currentRecord
	}
}

/* Upserts a Warehouse */
func UpsertWarehouse(
	dbconfig *DbConfig,
	currentRecord objects.ImportResponse,
	product objects.CSVProduct,
) objects.ImportResponse {
	variantID, exists, err := QueryVariantIDBySKU(dbconfig, product.SKU)
	if err != nil {
		currentRecord.FailCounter++
		return currentRecord
	}
	if exists {
		for _, warehouse := range product.Warehouses {
			err := AddWarehouse(dbconfig, product.SKU, variantID, warehouse.Name, warehouse.Value)
			if err != nil {
				currentRecord.FailCounter++
				return currentRecord
			}
		}
		return currentRecord
	} else {
		currentRecord.FailCounter++
		return currentRecord
	}
}

/* Upserts a product */
func UpsertProduct(
	dbconfig *DbConfig,
	currentRecord objects.ImportResponse,
	product objects.CSVProduct,
) objects.ImportResponse {
	dbProduct, err := dbconfig.DB.UpsertProduct(context.Background(), database.UpsertProductParams{
		ID:          uuid.New(),
		ProductCode: product.ProductCode,
		Active:      product.Active,
		Title:       utils.ConvertStringToSQL(product.Title),
		BodyHtml:    utils.ConvertStringToSQL(product.BodyHTML),
		Category:    utils.ConvertStringToSQL(product.Category),
		Vendor:      utils.ConvertStringToSQL(product.Vendor),
		ProductType: utils.ConvertStringToSQL(product.ProductType),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		currentRecord.FailCounter++
		return currentRecord
	}
	if dbProduct.Inserted {
		currentRecord.ProductsAdded++
	} else {
		currentRecord.ProductsUpdated++
	}
	option_names := CreateOptionNamesMap(product)
	err = AddProductOptions(dbconfig, dbProduct.ID, product.ProductCode, option_names)
	if err != nil {
		currentRecord.FailCounter++
		return currentRecord
	}
	return currentRecord
}

/* Upsert images */
func UpsertImages(
	dbconfig *DbConfig,
	currentRecord objects.ImportResponse,
	product objects.CSVProduct,
) objects.ImportResponse {
	// overwrite ones with the same position
	images := CreateImageMap(product)
	productID, exists, err := QueryProductByProductCode(dbconfig, product.ProductCode)
	if err != nil {
		currentRecord.FailCounter++
		return currentRecord
	}
	if exists {
		for key := range images {
			if images[key] != "" {
				err := AddImagery(dbconfig, productID, images[key], key+1)
				if err != nil {
					currentRecord.FailCounter++
					return currentRecord
				}
			}
		}
		return currentRecord
	} else {
		// not found so it counts as an error
		currentRecord.FailCounter++
		return currentRecord
	}
}

/* Upserts a variant into the database */
func UpsertVariant(
	dbconfig *DbConfig,
	currentRecord objects.ImportResponse,
	CSVProduct objects.CSVProduct,
) objects.ImportResponse {
	productID, exists, err := QueryProductByProductCode(dbconfig, CSVProduct.ProductCode)
	if err != nil {
		currentRecord.FailCounter++
		return currentRecord
	}
	if exists {
		dbVariant, err := dbconfig.DB.UpsertVariant(context.Background(), database.UpsertVariantParams{
			ID:        uuid.New(),
			ProductID: productID,
			Sku:       CSVProduct.SKU,
			Option1:   utils.ConvertStringToSQL(CSVProduct.Option1Value),
			Option2:   utils.ConvertStringToSQL(CSVProduct.Option2Value),
			Option3:   utils.ConvertStringToSQL(CSVProduct.Option3Value),
			Barcode:   utils.ConvertStringToSQL(CSVProduct.Barcode),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			currentRecord.FailCounter++
			return currentRecord
		}
		if dbVariant.Inserted {
			currentRecord.VariantsAdded++
		} else {
			currentRecord.VariantsUpdated++
		}
		return currentRecord
	} else {
		currentRecord.FailCounter++
		return currentRecord
	}
}

/* Adds a link between a customer and an order */
func AddCustomerOrder(dbconfig *DbConfig, orderID, customerID uuid.UUID) error {
	exists, err := CheckExistsCustomerOrder(dbconfig, context.Background(), customerID, orderID)
	if err != nil {
		return err
	}
	if !exists {
		err = dbconfig.DB.CreateCustomerOrder(context.Background(), database.CreateCustomerOrderParams{
			ID:         uuid.New(),
			CustomerID: customerID,
			OrderID:    orderID,
			UpdatedAt:  time.Now().UTC(),
			CreatedAt:  time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

/* Adds an order to the application */
func AddOrder(dbconfig *DbConfig, orderBody objects.RequestBodyOrder) (uuid.UUID, error) {
	exists, err := CheckExistsOrder(dbconfig, context.Background(), orderBody.Name)
	if err != nil {
		return uuid.Nil, err
	}
	if !exists {
		if err := OrderValidation(orderBody); err != nil {
			return uuid.Nil, err
		}
		dbOrder, err := dbconfig.DB.CreateOrder(context.Background(), database.CreateOrderParams{
			ID:            uuid.New(),
			Status:        orderBody.FinancialStatus,
			Notes:         utils.ConvertStringToSQL(""),
			WebCode:       orderBody.Name,
			TaxTotal:      utils.ConvertStringToSQL(orderBody.TotalTax),
			OrderTotal:    utils.ConvertStringToSQL(orderBody.TotalPrice),
			ShippingTotal: utils.ConvertStringToSQL(orderBody.TotalShippingPriceSet.ShopMoney.Amount),
			DiscountTotal: utils.ConvertStringToSQL(orderBody.TotalDiscounts),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		})
		if err != nil {
			return uuid.Nil, err
		}
		err = AddOrderLines(dbconfig, orderBody, dbOrder.ID)
		if err != nil {
			return uuid.Nil, err
		}
		dbCustomerUUID, err := AddCustomer(
			dbconfig,
			orderBody.Customer,
			orderBody.Customer.FirstName+" "+orderBody.Customer.LastName,
		)
		if err != nil {
			return uuid.Nil, err
		}
		err = AddCustomerOrder(dbconfig, dbOrder.ID, dbCustomerUUID)
		if err != nil {
			return uuid.Nil, err
		}
		return dbOrder.ID, nil
	}
	return uuid.Nil, nil
}

/* Updates an order that already exists inside the application */
func UpdateOrder(dbconfig *DbConfig, orderID uuid.UUID, orderBody objects.RequestBodyOrder) error {
	exists, err := CheckExistsOrder(dbconfig, context.Background(), orderBody.Name)
	if err != nil {
		return err
	}
	if exists {
		_, err = dbconfig.DB.UpdateOrder(context.Background(), database.UpdateOrderParams{
			Notes:         utils.ConvertStringToSQL(orderBody.Note),
			Status:        orderBody.FinancialStatus,
			WebCode:       orderBody.Name,
			TaxTotal:      utils.ConvertStringToSQL(orderBody.TotalTax),
			OrderTotal:    utils.ConvertStringToSQL(orderBody.TotalPrice),
			ShippingTotal: utils.ConvertStringToSQL(orderBody.TotalShippingPriceSet.ShopMoney.Amount),
			DiscountTotal: utils.ConvertStringToSQL(orderBody.TotalDiscounts),
			UpdatedAt:     time.Now().UTC(),
			ID:            orderID,
		})
		if err != nil {
			return err
		}
		// remove previous order line before adding new ones
		err = QueryClearOrderLines(dbconfig, orderID)
		if err != nil {
			return err
		}
		err = AddOrderLines(dbconfig, orderBody, orderID)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

/* Adds an order's line items to the database under the specific orderID */
func AddOrderLines(dbconfig *DbConfig, orderBody objects.RequestBodyOrder, orderID uuid.UUID) error {
	_, err := dbconfig.DB.GetOrderByID(context.Background(), orderID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("invalid order ID provided: " + orderID.String())
		}
	}
	for _, lineItem := range orderBody.LineItems {
		databaseLineItem := database.CreateOrderLineParams{
			ID:        uuid.New(),
			OrderID:   orderID,
			LineType:  utils.ConvertStringToSQL("product"),
			Sku:       lineItem.Sku,
			Price:     utils.ConvertStringToSQL(lineItem.Price),
			Qty:       utils.ConvertIntToSQL(lineItem.Quantity),
			TaxRate:   utils.ConvertStringToSQL(fmt.Sprint(lineItem.TaxLines[0].Rate)), // bad practise
			TaxTotal:  utils.ConvertStringToSQL(lineItem.TaxLines[0].Price),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		if len(lineItem.TaxLines) > 1 {
			databaseLineItem.TaxRate = utils.ConvertStringToSQL(fmt.Sprint(lineItem.TaxLines[0].Rate))
			databaseLineItem.TaxTotal = utils.ConvertStringToSQL(lineItem.TaxLines[0].Price)
		}
		_, err = dbconfig.DB.CreateOrderLine(context.Background(), databaseLineItem)
		if err != nil {
			return err
		}
	}
	for _, shippingLine := range orderBody.ShippingLines {
		databaseShippingLine := database.CreateOrderLineParams{
			ID:       uuid.New(),
			OrderID:  orderID,
			LineType: utils.ConvertStringToSQL("shipping"),
			Sku:      shippingLine.Code,
			Price:    utils.ConvertStringToSQL(shippingLine.Price),
			// TODO will this always remain as 1?
			Qty:       utils.ConvertIntToSQL(1),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		if len(shippingLine.TaxLines) > 1 {
			databaseShippingLine.TaxRate = utils.ConvertStringToSQL(fmt.Sprint(shippingLine.TaxLines[0].Rate))
			databaseShippingLine.TaxTotal = utils.ConvertStringToSQL(shippingLine.TaxLines[0].Price)
		}
		_, err = dbconfig.DB.CreateOrderLine(context.Background(), databaseShippingLine)
		if err != nil {
			return err
		}
	}
	return nil
}

/* Adds an order's line items to the database under the specific orderID */
func UpdateOrderLine(dbconfig *DbConfig, orderBody objects.RequestBodyOrder, orderID uuid.UUID) {
	// TODO is this even necessary
	// since we remove the old order upon receiving the new order?
}

/* Adds a customer to the application */
func AddCustomer(
	dbconfig *DbConfig,
	customer objects.RequestBodyCustomer,
	WebCustomerCode string,
) (uuid.UUID, error) {
	exists, err := CheckExistsCustomer(dbconfig, context.Background(), WebCustomerCode)
	if err != nil {
		return uuid.Nil, err
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
			return uuid.Nil, err
		}
		dbAddressUUID, err := AddAddress(dbconfig, customer.Address, dbCustomer.ID, "default")
		if err != nil {
			return uuid.Nil, err
		}
		err = AddCustomerAddress(dbconfig, customer.Address, dbCustomer.ID, dbAddressUUID, "default")
		if err != nil {
			return uuid.Nil, err
		}
		return dbCustomer.ID, nil
	}
	return uuid.Nil, nil
}

/* Updates a customer inside the application */
func UpdateCustomer(dbconfig *DbConfig, customer objects.RequestBodyCustomer, WebCustomerCode string) error {
	exists, err := CheckExistsCustomer(dbconfig, context.Background(), WebCustomerCode)
	if err != nil {
		return err
	}
	if exists {
		err := dbconfig.DB.UpdateCustomer(context.Background(), database.UpdateCustomerParams{
			ID:        uuid.New(),
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     utils.ConvertStringToSQL(customer.Email),
			Phone:     utils.ConvertStringToSQL(customer.Phone),
			UpdatedAt: time.Now().UTC(),
		})
		return err
	}
	return nil
}

/* Adds a customer address which is a link between a customer and it's address */
func AddCustomerAddress(
	dbconfig *DbConfig,
	addressData objects.CustomerAddress,
	customerID,
	addressID uuid.UUID,
	addressType string,
) error {
	exists, err := CheckExistsCustomerAddress(dbconfig, context.Background(), customerID.String(), addressType)
	if err != nil {
		return err
	}
	if !exists {
		err = dbconfig.DB.CreateCustomerAddress(context.Background(), database.CreateCustomerAddressParams{
			ID:          uuid.New(),
			CustomerID:  customerID,
			AddressID:   addressID,
			AddressType: addressType,
			UpdatedAt:   time.Now().UTC(),
			CreatedAt:   time.Now().UTC(),
		})
		return err
	}
	return nil
}

/* Updates a customer address */
func UpdateCustomerAddress(dbconfig *DbConfig, orderData objects.RequestBodyOrder, customerID uuid.UUID, addressType string) error {
	exists, err := CheckExistsCustomerAddress(dbconfig, context.Background(), customerID.String(), addressType)
	if err != nil {
		return err
	}
	if exists {
		// TODO should we update customer address links?
		// also, does it make sense to update it seeing that it already exists in this block
		return nil
	}
	return nil
}

/* Adds a product to the application */
func AddProduct(dbconfig *DbConfig, productData objects.RequestBodyProduct) (uuid.UUID, int, error) {
	if validation := ProductValidation(dbconfig, productData); validation != nil {
		return uuid.Nil, http.StatusBadRequest, validation
	}
	if err := ValidateDuplicateOption(productData); err != nil {
		return uuid.Nil, http.StatusBadRequest, err
	}
	if err := ValidateDuplicateSKU(productData, dbconfig); err != nil {
		return uuid.Nil, http.StatusConflict, err
	}
	productID, exists, err := QueryProductByProductCode(dbconfig, productData.ProductCode)
	if err != nil {
		return uuid.Nil, 500, err
	}
	if !exists {
		product, err := dbconfig.DB.CreateProduct(context.Background(), database.CreateProductParams{
			ID:          uuid.New(),
			Active:      productData.Active,
			ProductCode: productData.ProductCode,
			Title:       utils.ConvertStringToSQL(productData.Title),
			BodyHtml:    utils.ConvertStringToSQL(productData.BodyHTML),
			Category:    utils.ConvertStringToSQL(productData.Category),
			Vendor:      utils.ConvertStringToSQL(productData.Vendor),
			ProductType: utils.ConvertStringToSQL(productData.ProductType),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})
		if err != nil {
			return uuid.Nil, http.StatusInternalServerError, err
		}
		for key := range productData.ProductOptions {
			_, err := dbconfig.DB.CreateProductOption(context.Background(), database.CreateProductOptionParams{
				ID:        uuid.New(),
				ProductID: product.ID,
				Name:      productData.ProductOptions[key].Value,
				Position:  int32(key + 1),
			})
			if err != nil {
				return uuid.Nil, http.StatusInternalServerError, err
			}
		}
		productID = product.ID
	}
	for _, variant := range productData.Variants {
		if httpCode, err := AddVariant(dbconfig, variant, productID); err != nil {
			return uuid.Nil, httpCode, err
		}
	}
	return productID, 0, nil
}

/* Updates a product to the application */
func UpdateProduct(dbconfig *DbConfig, productData objects.RequestBodyProduct, productID, apiKey string) error {
	productUUID, err := QueryProductByID(dbconfig, productID)
	if err != nil {
		return err
	}
	if productData.Active == "" {
		productData.Active = "0"
	}
	err = dbconfig.DB.UpdateProductByID(context.Background(), database.UpdateProductByIDParams{
		Active:      productData.Active,
		Title:       utils.ConvertStringToSQL(productData.Title),
		BodyHtml:    utils.ConvertStringToSQL(productData.BodyHTML),
		Category:    utils.ConvertStringToSQL(productData.Category),
		Vendor:      utils.ConvertStringToSQL(productData.Vendor),
		ProductType: utils.ConvertStringToSQL(productData.ProductCode),
		UpdatedAt:   time.Now().UTC(),
		ID:          productUUID,
	})
	if err != nil {
		return err
	}
	for productOptionKey := range productData.ProductOptions {
		_, err = dbconfig.DB.UpdateProductOption(context.Background(), database.UpdateProductOptionParams{
			Name:       productData.ProductOptions[productOptionKey].Value,
			Position:   int32(productOptionKey + 1),
			ProductID:  productUUID,
			Position_2: int32(productOptionKey + 1),
		})
		if err != nil {
			return err
		}
	}
	for _, variantData := range productData.Variants {
		err = UpdateVariant(dbconfig, variantData, productID)
		if err != nil {
			return err
		}
	}
	if productData.Active == "1" {
		err = UpdateShopifyProduct(dbconfig, productUUID, apiKey)
		if err != nil {
			return err
		}
	}
	return nil
}

/* Pushes an update for a product and it's variants to Shopify if the product is Active */
func UpdateShopifyProduct(dbconfig *DbConfig, productID uuid.UUID, apiKey string) error {
	productData, err := CompileProduct(dbconfig, productID, context.Background(), false)
	if err != nil {
		return err
	}
	err = CompileInstructionProduct(dbconfig, productData, apiKey)
	if err != nil {
		return err
	}
	for _, variant := range productData.Variants {
		err = CompileInstructionVariant(dbconfig, variant, productData, apiKey)
		if err != nil {
			return err
		}
	}
	return nil
}

/* Adds a product variant to the application. The productID needs to point to a valid product*/
func AddVariant(dbconfig *DbConfig, variantData objects.RequestBodyVariant, productID uuid.UUID) (int, error) {
	if err := DuplicateOptionValues(dbconfig, variantData, productID); err != nil {
		return http.StatusBadRequest, err
	}
	variant, err := dbconfig.DB.CreateVariant(context.Background(), database.CreateVariantParams{
		ID:        uuid.New(),
		ProductID: productID,
		Sku:       variantData.Sku,
		Option1:   utils.ConvertStringToSQL(variantData.Option1),
		Option2:   utils.ConvertStringToSQL(variantData.Option2),
		Option3:   utils.ConvertStringToSQL(variantData.Option3),
		Barcode:   utils.ConvertStringToSQL(variantData.Barcode),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for key_pricing := range variantData.VariantPricing {
		err = AddPricing(
			dbconfig,
			variant.Sku,
			variant.ID,
			variantData.VariantPricing[key_pricing].Name,
			variantData.VariantPricing[key_pricing].Value,
		)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	for key_qty := range variantData.VariantQuantity {
		err = AddWarehouse(
			dbconfig,
			variant.Sku,
			variant.ID,
			variantData.VariantQuantity[key_qty].Name,
			variantData.VariantQuantity[key_qty].Value,
		)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}

/* Updates a product variant inside the application. The productID needs to point to a valid product */
func UpdateVariant(
	dbconfig *DbConfig,
	variantData objects.RequestBodyVariant,
	productID string,
) error {
	dbVariantID, exists, err := QueryVariantIDBySKU(dbconfig, variantData.Sku)
	if err != nil {
		return err
	}
	if exists {
		err = dbconfig.DB.UpdateVariant(context.Background(), database.UpdateVariantParams{
			Option1:   utils.ConvertStringToSQL(variantData.Option1),
			Option2:   utils.ConvertStringToSQL(variantData.Option2),
			Option3:   utils.ConvertStringToSQL(variantData.Option3),
			Barcode:   utils.ConvertStringToSQL(variantData.Barcode),
			UpdatedAt: time.Now().UTC(),
			Sku:       variantData.Sku,
		})
		if err != nil {
			return err
		}
		for key_pricing := range variantData.VariantPricing {
			err = AddPricing(
				dbconfig,
				variantData.Sku,
				dbVariantID,
				variantData.VariantPricing[key_pricing].Name,
				variantData.VariantPricing[key_pricing].Value,
			)
			if err != nil {
				return err
			}
		}
		for key_qty := range variantData.VariantQuantity {
			err = AddWarehouse(
				dbconfig,
				variantData.Sku,
				dbVariantID,
				variantData.VariantQuantity[key_qty].Name,
				variantData.VariantQuantity[key_qty].Value,
			)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("'" + variantData.Sku + "' do not exist")
	}
	return nil
}

/* Adds a new User to the application */
func AddUser(dbconfig *DbConfig, userData objects.RequestBodyRegister) (database.User, error) {
	user, err := dbconfig.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      userData.Name,
		UserType:  "app",
		Email:     userData.Email,
		Password:  userData.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}

/* Adds a new User to the application */
func AddUserRegistration(
	dbconfig *DbConfig,
	preRegisterDetails objects.RequestBodyPreRegister,
) (database.RegisterToken, error) {
	token, err := dbconfig.DB.CreateToken(context.Background(), database.CreateTokenParams{
		ID:        uuid.New(),
		Name:      preRegisterDetails.Name,
		Email:     preRegisterDetails.Email,
		Token:     uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.RegisterToken{}, err
	}
	return token, nil
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
	if pricing_name != "Selling Price" && pricing_name != "Compare At Price" {
		return errors.New("invalid price '" + pricing_name + "'")
	}
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
func AddAddress(dbconfig *DbConfig, addressData objects.CustomerAddress, customerID uuid.UUID, addressType string) (uuid.UUID, error) {
	exists, err := CheckExistsCustomerAddress(dbconfig, context.Background(), customerID.String(), addressType)
	if err != nil {
		return uuid.Nil, err
	}
	if !exists {
		dbCustomerAddress, err := dbconfig.DB.CreateAddress(context.Background(), database.CreateAddressParams{
			ID:           uuid.New(),
			CustomerID:   customerID,
			Type:         addressType,
			FirstName:    addressData.FirstName,
			LastName:     addressData.LastName,
			Address1:     utils.ConvertStringToSQL(addressData.Address1),
			Address2:     utils.ConvertStringToSQL(addressData.Address2),
			City:         utils.ConvertStringToSQL(addressData.City),
			Province:     utils.ConvertStringToSQL(addressData.Province),
			ProvinceCode: utils.ConvertStringToSQL(addressData.ProvinceCode),
			Company:      utils.ConvertStringToSQL(addressData.Company),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		})
		if err != nil {
			return uuid.Nil, err
		}
		return dbCustomerAddress.ID, err
	}
	return uuid.Nil, nil
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
func CheckVID(dbconfig *DbConfig, sku string, r *http.Request) (string, error) {
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
func CheckPID(dbconfig *DbConfig, product_code string, r *http.Request) (string, error) {
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

/* Checks if an order already exists inside the database using it's web code */
func CheckExistsOrderByID(dbconfig *DbConfig, ctx context.Context, orderID uuid.UUID) (bool, error) {
	dbOrder, err := dbconfig.DB.GetOrderByID(ctx, orderID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	if dbOrder.ID == orderID {
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

/* Checks if the customer-order link already exists inside the database */
func CheckExistsCustomerOrder(dbconfig *DbConfig, ctx context.Context, customerID, orderID uuid.UUID) (bool, error) {
	_, err := dbconfig.DB.GetOrderIDByCustomerID(ctx, database.GetOrderIDByCustomerIDParams{
		CustomerID: customerID,
		OrderID:    orderID,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
func CheckUExistsUser(dbconfig *DbConfig, name string, r *http.Request) (bool, error) {
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
func CheckUserCredentials(
	dbconfig *DbConfig,
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
func CheckExistsToken(dbconfig *DbConfig, email string, r *http.Request) (uuid.UUID, bool, error) {
	token, err := dbconfig.DB.GetToken(r.Context(), email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return uuid.Nil, false, err
		}
	}
	if token.Email == email {
		return token.Token, true, nil
	}
	return uuid.Nil, false, nil
}

/* Checks if a token already exists in the database */
func CheckUserEmailType(dbconfig *DbConfig, email, user_type string) (bool, error) {
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
