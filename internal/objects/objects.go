package objects

import (
	"time"

	"github.com/google/uuid"
)

// iocsv.go

type ImportResponse struct {
	ProcessedCounter int `json:"processed_counter"`
	FailCounter      int `json:"fail_counter"`
	SkipCounter      int `json:"skip_counter"`
	ProductsAdded    int `json:"products_added"`
	ProductsUpdated  int `json:"products_updated"`
	VariantsAdded    int `json:"variants_added"`
	VariantsUpdated  int `json:"variants_updated"`
}

type CSVProduct struct {
	ProductCode  string `csv:"product_code"`
	Active       string `csv:"active"`
	Title        string `csv:"title"`
	BodyHTML     string `csv:"body_html"`
	Category     string `csv:"category"`
	Vendor       string `csv:"vendor"`
	ProductType  string `csv:"product_type"`
	SKU          string `csv:"sku"`
	Option1Name  string `csv:"option1_name"`
	Option1Value string `csv:"option1_value"`
	Option2Name  string `csv:"option2_name"`
	Option2Value string `csv:"option2_value"`
	Option3Name  string `csv:"option3_name"`
	Option3Value string `csv:"option3_value"`
	Barcode      string `csv:"barcode"`
	PriceName    string `csv:"price_name"`
	PriceValue   string `csv:"price_value"`
	QtyName      string `csv:"qty_name"`
	QtyValue     int    `csv:"qty_value"`
}

type ResponseString struct {
	Message string `json:"message"`
}

type RequestString struct {
	Message string `json:"message"`
}

// request_validation.go

type RequestBodyUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
type RequestBodyProduct struct {
	ProductCode    string           `json:"product_code"`
	Title          string           `json:"title"`
	BodyHTML       string           `json:"body_html"`
	Category       string           `json:"category"`
	Vendor         string           `json:"vendor"`
	ProductType    string           `json:"product_type"`
	Variants       []ProductVariant `json:"variants"`
	ProductOptions []ProductOptions `json:"options"`
}

type RequestBodyCustomer struct {
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Address   []CustomerAddress `json:"address"`
}

type RequestBodyPreRegister struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RequestBodyValidateToken struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// object_converter.go
type SearchOrder struct {
	Notes         string    `json:"notes"`
	WebCode       string    `json:"web_code"`
	TaxTotal      string    `json:"tax_total"`
	OrderTotal    string    `json:"order_total"`
	ShippingTotal string    `json:"shipping_total"`
	DiscountTotal string    `json:"discount_total"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type Order struct {
	ID                uuid.UUID     `json:"id"`
	Notes             string        `json:"notes"`
	WebCode           string        `json:"web_code"`
	TaxTotal          string        `json:"tax_total"`
	OrderTotal        string        `json:"order_total"`
	ShippingTotal     string        `json:"shipping_total"`
	DiscountTotal     string        `json:"discount_total"`
	UpdatedAt         time.Time     `json:"updated_at"`
	CreatedAt         time.Time     `json:"created_at"`
	OrderCustomer     OrderCustomer `json:"customer"`
	LineItems         []OrderLines  `json:"line_items"`
	ShippingLineItems []OrderLines  `json:"shipping_lines"`
}
type OrderAddress struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address1   string `json:"address_1"`
	Address2   string `json:"address_2"`
	Suburb     string `json:"suburb"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	Company    string `json:"company"`
}

type OrderCustomer struct {
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Address   []CustomerAddress `json:"shipping_address"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Customer struct {
	ID        uuid.UUID         `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Address   []CustomerAddress `json:"shipping_address"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type CustomerAddress struct {
	Type       string `json:"address_type"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address1   string `json:"address_1"`
	Address2   string `json:"address_2"`
	Suburb     string `json:"suburb"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	Company    string `json:"company"`
}

type OrderLines struct {
	SKU      string `json:"sku"`
	Price    string `json:"price"`
	Barcode  int    `json:"barcode"`
	Qty      int    `json:"qty"`
	TaxRate  string `json:"tax_rate"`
	TaxTotal string `json:"tax_total"`
}
type SearchCustomer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type SearchProduct struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	ProductType string    `json:"product_type"`
	Vendor      string    `json:"vendor"`
}
type Product struct {
	ID             uuid.UUID        `json:"id"`
	Active         string           `json:"active"`
	Title          string           `json:"title"`
	BodyHTML       string           `json:"body_html"`
	Category       string           `json:"category"`
	Vendor         string           `json:"vendor"`
	ProductType    string           `json:"product_type"`
	Variants       []ProductVariant `json:"variants"`
	ProductOptions []ProductOptions `json:"options"`
	UpdatedAt      time.Time        `json:"updated_at"`
}
type ProductOptions struct {
	Value string `json:"value"`
}
type ProductVariant struct {
	Sku             string         `json:"sku"`
	Option1         string         `json:"option1"`
	Option2         string         `json:"option2"`
	Option3         string         `json:"option3"`
	Barcode         string         `json:"barcode"`
	VariantPricing  []VariantPrice `json:"variant_price_tiers"`
	VariantQuantity []VariantQty   `json:"variant_quantities"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type VariantPrice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type VariantQty struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Api Endpoints
type Endpoints struct {
	Status      bool             `json:"status"`
	Description string           `json:"description"`
	Routes      map[string]Route `json:"routes"`
	Version     string           `json:"version"`
	Time        time.Time        `json:"time"`
}

type Route struct {
	Description   string            `json:"description"`
	Supports      []string          `json:"supports"`
	Params        map[string]Params `json:"params"`
	AcceptsData   bool              `json:"accepts_data"`
	Format        any               `json:"format"`
	Authorization string            `json:"auth"`
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Product Feed: Shopify

type ShopifyProducts struct {
	Products []ShopifyProduct `json:"products"`
}

type ShopifyProduct struct {
	Title    string           `json:"title"`
	BodyHTML string           `json:"body_html"`
	Vendor   string           `json:"vendor"`
	Type     string           `json:"product_type"`
	Status   string           `json:"status"`
	Variants []ShopifyVariant `json:"variants"`
	Options  []ShopifyOptions `json:"options"`
}

type ShopifyVariant struct {
	Sku            string `json:"sku"`
	Price          string `json:"price"`
	CompareAtPrice string `json:"compare_at_price"`
	Option1        string `json:"option1"`
	Option2        string `json:"option2"`
	Option3        string `json:"option3"`
	Barcode        string `json:"barcode"`
}

type ShopifyOptions struct {
	Name     string `json:"name"`
	Values   string `json:"values"`
	Position string `json:"position"`
}
