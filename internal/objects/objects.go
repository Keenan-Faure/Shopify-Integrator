package objects

import "time"

type ResponseString struct {
	Status string
}

// request_validation.go
type RequestBodyUser struct {
	Name string
}

// object_converter.go
type SearchProduct struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	ProductType string `json:"product_type"`
	Vendor      string `json:"vendor"`
}
type Product struct {
	Active         string           `json:"active"`
	Title          string           `json:"title"`
	BodyHTML       string           `json:"body_html"`
	Category       string           `json:"category"`
	Vendor         string           `json:"vendor"`
	ProductType    string           `json:"product_type"`
	Variants       []ProductVariant `json:"variants"`
	ProductOptions []ProductOptions `json:"options"`
	UpdatedAt      string           `json:"updated_at"`
}
type ProductOptions struct {
	Name  string `json:"name"`
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
	UpdatedAt       string         `json:"updated_at"`
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
