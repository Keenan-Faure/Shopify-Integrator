package objects

import "time"

type ResponseString struct {
	Status string
}

// request_validation.go
type RequestBodyUser struct {
	Name string
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
	Weight         string `json:"weight"`
	Barcode        string `json:"barcode"`
}

type ShopifyOptions struct {
	Name     string `json:"name"`
	Values   string `json:"values"`
	Position string `json:"position"`
}
