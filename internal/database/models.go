// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID         uuid.UUID      `json:"id"`
	CustomerID uuid.UUID      `json:"customer_id"`
	Name       sql.NullString `json:"name"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Address1   sql.NullString `json:"address1"`
	Address2   sql.NullString `json:"address2"`
	Suburb     sql.NullString `json:"suburb"`
	City       sql.NullString `json:"city"`
	Province   sql.NullString `json:"province"`
	PostalCode sql.NullString `json:"postal_code"`
	Company    sql.NullString `json:"company"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type Customer struct {
	ID        uuid.UUID      `json:"id"`
	OrderID   uuid.UUID      `json:"order_id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Order struct {
	ID            uuid.UUID      `json:"id"`
	CustomerID    uuid.UUID      `json:"customer_id"`
	Notes         sql.NullString `json:"notes"`
	WebCode       sql.NullString `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type OrderLine struct {
	ID        uuid.UUID      `json:"id"`
	OrderID   uuid.UUID      `json:"order_id"`
	LineType  sql.NullString `json:"line_type"`
	Sku       string         `json:"sku"`
	Price     sql.NullString `json:"price"`
	Barcode   sql.NullInt32  `json:"barcode"`
	Qty       sql.NullInt32  `json:"qty"`
	TaxTotal  sql.NullString `json:"tax_total"`
	TaxRate   sql.NullString `json:"tax_rate"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Product struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ProductOption struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
}

type RegisterToken struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	WebhookToken string    `json:"webhook_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	ApiKey       string    `json:"api_key"`
}

type Variant struct {
	ID        uuid.UUID      `json:"id"`
	ProductID uuid.UUID      `json:"product_id"`
	Sku       string         `json:"sku"`
	Option1   sql.NullString `json:"option1"`
	Option2   sql.NullString `json:"option2"`
	Option3   sql.NullString `json:"option3"`
	Barcode   sql.NullString `json:"barcode"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type VariantPricing struct {
	ID        uuid.UUID      `json:"id"`
	VariantID uuid.UUID      `json:"variant_id"`
	Name      string         `json:"name"`
	Value     sql.NullString `json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type VariantQty struct {
	ID        uuid.UUID     `json:"id"`
	VariantID uuid.UUID     `json:"variant_id"`
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
