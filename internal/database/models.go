// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID         uuid.UUID      `json:"id"`
	CustomerID uuid.UUID      `json:"customer_id"`
	Type       sql.NullString `json:"type"`
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

type AppSetting struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	FieldName   string    `json:"field_name"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Customer struct {
	ID        uuid.UUID      `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Customerorder struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	OrderID    uuid.UUID `json:"order_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FetchRestriction struct {
	ID        uuid.UUID `json:"id"`
	Field     string    `json:"field"`
	Flag      string    `json:"flag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FetchStat struct {
	ID               uuid.UUID `json:"id"`
	AmountOfProducts int32     `json:"amount_of_products"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type FetchWorker struct {
	ID                  uuid.UUID `json:"id"`
	Status              string    `json:"status"`
	LocalCount          int32     `json:"local_count"`
	ShopifyProductCount int32     `json:"shopify_product_count"`
	FetchUrl            string    `json:"fetch_url"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type GoogleOauth struct {
	ID           uuid.UUID      `json:"id"`
	UserID       uuid.UUID      `json:"user_id"`
	CookieSecret string         `json:"cookie_secret"`
	CookieToken  string         `json:"cookie_token"`
	GoogleID     string         `json:"google_id"`
	Email        string         `json:"email"`
	Picture      sql.NullString `json:"picture"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type InventoryLocation struct {
	ID                uuid.UUID `json:"id"`
	ShopifyLocationID string    `json:"shopify_location_id"`
	InventoryItemID   string    `json:"inventory_item_id"`
	WarehouseName     string    `json:"warehouse_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Order struct {
	ID            uuid.UUID      `json:"id"`
	Notes         sql.NullString `json:"notes"`
	WebCode       sql.NullString `json:"web_code"`
	TaxTotal      sql.NullString `json:"tax_total"`
	OrderTotal    sql.NullString `json:"order_total"`
	ShippingTotal sql.NullString `json:"shipping_total"`
	DiscountTotal sql.NullString `json:"discount_total"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Status        string         `json:"status"`
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

type OrderStat struct {
	ID         uuid.UUID `json:"id"`
	OrderTotal int32     `json:"order_total"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Product struct {
	ID          uuid.UUID      `json:"id"`
	Active      string         `json:"active"`
	ProductCode string         `json:"product_code"`
	Title       sql.NullString `json:"title"`
	BodyHtml    sql.NullString `json:"body_html"`
	Category    sql.NullString `json:"category"`
	Vendor      sql.NullString `json:"vendor"`
	ProductType sql.NullString `json:"product_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ProductImage struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	Position  int32     `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductOption struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Position  int32     `json:"position"`
}

type PushRestriction struct {
	ID        uuid.UUID `json:"id"`
	Field     string    `json:"field"`
	Flag      string    `json:"flag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type QueueItem struct {
	ID          uuid.UUID       `json:"id"`
	QueueType   string          `json:"queue_type"`
	Instruction string          `json:"instruction"`
	Status      string          `json:"status"`
	Object      json.RawMessage `json:"object"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Description string          `json:"description"`
}

type RegisterToken struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShopifyCollection struct {
	ID                  uuid.UUID      `json:"id"`
	ProductCollection   sql.NullString `json:"product_collection"`
	ShopifyCollectionID string         `json:"shopify_collection_id"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

type ShopifyInventory struct {
	ID                uuid.UUID `json:"id"`
	ShopifyLocationID string    `json:"shopify_location_id"`
	InventoryItemID   string    `json:"inventory_item_id"`
	Available         int32     `json:"available"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ShopifyLocation struct {
	ID                   uuid.UUID `json:"id"`
	ShopifyWarehouseName string    `json:"shopify_warehouse_name"`
	ShopifyLocationID    string    `json:"shopify_location_id"`
	WarehouseName        string    `json:"warehouse_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ShopifyPid struct {
	ID               uuid.UUID `json:"id"`
	ProductCode      string    `json:"product_code"`
	ProductID        uuid.UUID `json:"product_id"`
	ShopifyProductID string    `json:"shopify_product_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ShopifySetting struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	FieldName   string    `json:"field_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
}

type ShopifyVid struct {
	ID                 uuid.UUID `json:"id"`
	Sku                string    `json:"sku"`
	VariantID          uuid.UUID `json:"variant_id"`
	ShopifyVariantID   string    `json:"shopify_variant_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ShopifyInventoryID string    `json:"shopify_inventory_id"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	ApiKey       string    `json:"api_key"`
	WebhookToken string    `json:"webhook_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	Isdefault bool           `json:"isdefault"`
}

type VariantQty struct {
	ID        uuid.UUID     `json:"id"`
	VariantID uuid.UUID     `json:"variant_id"`
	Name      string        `json:"name"`
	Value     sql.NullInt32 `json:"value"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Isdefault bool          `json:"isdefault"`
}

type Warehouse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
