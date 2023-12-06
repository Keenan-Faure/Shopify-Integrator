package objects

import (
	"time"

	"github.com/google/uuid"
)

// invisible endpoints

type FetchAmountResponse struct {
	Amounts []int64  `json:"amounts"`
	Hours   []string `json:"hours"`
}

type OrderAmountResponse struct {
	Count []int64  `json:"count"`
	Days  []string `json:"days"`
}

type RequestWebhookURL struct {
	ForwardingURL string `json:"forwarding_url"`
}

type ResponseWarehouseLocation struct {
	Warehouses       []string `json:"warehouses"`
	ShopifyLocations any      `json:"shopify_locations"`
}

type ShopifyLocations struct {
	Locations []struct {
		ID                    int64     `json:"id"`
		Name                  string    `json:"name"`
		Address1              string    `json:"address1"`
		Address2              string    `json:"address2"`
		City                  string    `json:"city"`
		Zip                   string    `json:"zip"`
		Province              string    `json:"province"`
		Country               string    `json:"country"`
		Phone                 string    `json:"phone"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		CountryCode           string    `json:"country_code"`
		CountryName           string    `json:"country_name"`
		ProvinceCode          string    `json:"province_code"`
		Legacy                bool      `json:"legacy"`
		Active                bool      `json:"active"`
		AdminGraphqlAPIID     string    `json:"admin_graphql_api_id"`
		LocalizedCountryName  string    `json:"localized_country_name"`
		LocalizedProvinceName string    `json:"localized_province_name"`
	} `json:"locations"`
}

// queue.go

type ResponseQueueWorker struct {
	ID          string      `json:"id"`
	QueueType   string      `json:"queue_type"`
	Instruction string      `json:"instruction"`
	Status      string      `json:"status"`
	Object      interface{} `json:"object"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ResponseQueueCount struct {
	AddOrder      int `json:"add_order"`
	AddProduct    int `json:"add_product"`
	AddVariant    int `json:"add_variant"`
	UpdateOrder   int `json:"update_order"`
	UpdateProduct int `json:"update_product"`
	UpdateVariant int `json:"update_variant"`
}

type ResponseQueueItemFilter struct {
	ID          uuid.UUID   `json:"id"`
	QueueType   string      `json:"queue_type"`
	Status      string      `json:"status"`
	Instruction string      `json:"instruction"`
	Object      interface{} `json:"object"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
type RequestQueueHelper struct {
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Instruction string      `json:"instruction"`
	Endpoint    string      `json:"endpoint"`
	ApiKey      string      `json:"api_key"`
	Method      string      `json:"method"`
	Object      interface{} `json:"object"`
}
type ResponseQueueItem struct {
	ID     uuid.UUID   `json:"queue_id"`
	Object interface{} `json:"object"`
}

type RequestQueueItem struct {
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Instruction string      `json:"instruction"`
	Object      interface{} `json:"object"`
}

type RequestQueueItemProducts struct {
	SystemProductID string `json:"system_product_id"`
	SystemVariantID string `json:"system_variant_id"`
	Shopify         struct {
		ProductID string `json:"product_id"`
		VariantID string `json:"variant_id"`
	} `json:"shopify"`
}

// shopify_push.go
type ShopifyIDs struct {
	ProductID string        `json:"product_id"`
	Variants  []ShopifyVIDs `json:"variants"`
}

type ShopifyVIDs struct {
	VariantID string `json:"variant_id"`
}

type AddShopifyCustomCollection struct {
	CustomCollection struct {
		Title string `json:"title"`
	} `json:"custom_collection"`
}

type AddProducToShopifyCollection struct {
	Collect struct {
		ProductID    int `json:"product_id"`
		CollectionID int `json:"collection_id"`
	} `json:"collect"`
}

type AddInventoryItem struct {
	LocationID          int `json:"location_id"`
	InventoryItemID     int `json:"inventory_item_id"`
	AvailableAdjustment int `json:"available_adjustment"`
}

type AddInventoryItemToLocation struct {
	LocationID      int `json:"location_id"`
	InventoryItemID int `json:"inventory_item_id"`
}

type GetShopifyInventoryLevelsList struct {
	InventoryLevels []GetShopifyInventoryLevels `json:"inventory_levels"`
}

type GetShopifyInventoryLevels struct {
	InventoryItemID int `json:"inventory_item_id"`
	Available       int `json:"available"`
	LocationID      int `json:"location_id"`
}

type ResponseShopifyWarehouseLocation struct {
	ID                   uuid.UUID `json:"id"`
	LocationID           string    `json:"location_id"`
	WarehouseName        string    `json:"warehouse_name"`
	ShopifyWarehouseName string    `json:"shopify_warehouse_name"`
}

type ResponseIDs struct {
	ProductID string `json:"id"`
	VariantID string `json:"variant_id"`
}

type ResponseAddInventoryItem struct {
	InventoryLevel struct {
		InventoryItemID   int    `json:"inventory_item_id"`
		LocationID        int    `json:"location_id"`
		Available         int    `json:"available"`
		UpdatedAt         string `json:"updated_at"`
		AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	} `json:"inventory_level"`
}

type ResponseAddInventoryItemLocation struct {
	InventoryLevel struct {
		InventoryItemID   int    `json:"inventory_item_id"`
		LocationID        int    `json:"location_id"`
		Available         int    `json:"available"`
		UpdatedAt         string `json:"updated_at"`
		AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	} `json:"inventory_level"`
}

type ResponseAddProductToShopifyCollection struct {
	Collect struct {
		ID           int    `json:"id"`
		CollectionID int    `json:"collection_id"`
		ProductID    int    `json:"product_id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Position     int    `json:"position"`
		SortValue    string `json:"sort_value"`
	} `json:"collect"`
}

type ShopifyVariantResponse struct {
	Variant struct {
		ID                   int     `json:"id"`
		ProductID            int     `json:"product_id"`
		Title                string  `json:"title"`
		Price                string  `json:"price"`
		Sku                  string  `json:"sku"`
		Position             int     `json:"position"`
		InventoryPolicy      string  `json:"inventory_policy"`
		CompareAtPrice       any     `json:"compare_at_price"`
		FulfillmentService   string  `json:"fulfillment_service"`
		InventoryManagement  string  `json:"inventory_management"`
		Option1              string  `json:"option1"`
		Option2              any     `json:"option2"`
		Option3              any     `json:"option3"`
		CreatedAt            string  `json:"created_at"`
		UpdatedAt            string  `json:"updated_at"`
		Taxable              bool    `json:"taxable"`
		Barcode              any     `json:"barcode"`
		Grams                int     `json:"grams"`
		ImageID              any     `json:"image_id"`
		Weight               float64 `json:"weight"`
		WeightUnit           string  `json:"weight_unit"`
		InventoryItemID      int     `json:"inventory_item_id"`
		InventoryQuantity    int     `json:"inventory_quantity"`
		OldInventoryQuantity int     `json:"old_inventory_quantity"`
		PresentmentPrices    []struct {
			Price struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"price"`
			CompareAtPrice any `json:"compare_at_price"`
		} `json:"presentment_prices"`
		RequiresShipping  bool   `json:"requires_shipping"`
		AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	} `json:"variant"`
}

type ShopifyProductResponse struct {
	Product struct {
		BodyHTML  string `json:"body_html"`
		CreatedAt string `json:"created_at"`
		Handle    string `json:"handle"`
		ID        int    `json:"id"`
		Images    []struct {
			ID         int    `json:"id"`
			ProductID  int    `json:"product_id"`
			Position   int    `json:"position"`
			CreatedAt  string `json:"created_at"`
			UpdatedAt  string `json:"updated_at"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			Src        string `json:"src"`
			VariantIds []struct {
			} `json:"variant_ids"`
		} `json:"images"`
		Options []struct {
			ID        int      `json:"id"`
			ProductID int      `json:"product_id"`
			Name      string   `json:"name"`
			Position  int      `json:"position"`
			Values    []string `json:"values"`
		} `json:"options"`
		ProductType    string `json:"product_type"`
		PublishedAt    string `json:"published_at"`
		PublishedScope string `json:"published_scope"`
		Status         string `json:"status"`
		Tags           string `json:"tags"`
		TemplateSuffix string `json:"template_suffix"`
		Title          string `json:"title"`
		UpdatedAt      string `json:"updated_at"`
		Variants       []struct {
			Barcode             string  `json:"barcode"`
			CompareAtPrice      any     `json:"compare_at_price"`
			CreatedAt           string  `json:"created_at"`
			FulfillmentService  string  `json:"fulfillment_service"`
			Grams               int     `json:"grams"`
			Weight              float64 `json:"weight"`
			WeightUnit          string  `json:"weight_unit"`
			ID                  int     `json:"id"`
			InventoryItemID     int     `json:"inventory_item_id"`
			InventoryManagement string  `json:"inventory_management"`
			InventoryPolicy     string  `json:"inventory_policy"`
			InventoryQuantity   int     `json:"inventory_quantity"`
			Option1             string  `json:"option1"`
			Position            int     `json:"position"`
			Price               string  `json:"price"`
			ProductID           int     `json:"product_id"`
			RequiresShipping    bool    `json:"requires_shipping"`
			Sku                 string  `json:"sku"`
			Taxable             bool    `json:"taxable"`
			Title               string  `json:"title"`
			UpdatedAt           string  `json:"updated_at"`
		} `json:"variants"`
		Vendor string `json:"vendor"`
	} `json:"product"`
}
type ResponseShopifyCustomCollection struct {
	CustomCollection struct {
		ID                int    `json:"id"`
		Handle            string `json:"handle"`
		Title             string `json:"title"`
		UpdatedAt         string `json:"updated_at"`
		BodyHTML          any    `json:"body_html"`
		PublishedAt       string `json:"published_at"`
		SortOrder         string `json:"sort_order"`
		TemplateSuffix    any    `json:"template_suffix"`
		PublishedScope    string `json:"published_scope"`
		AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	} `json:"custom_collection"`
}

type ResponseGetCustomCollections struct {
	CustomCollections []struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	} `json:"custom_collections"`
}

type ResponseShopifyGetLocations struct {
	Locations []struct {
		ID                    int64     `json:"id"`
		Name                  string    `json:"name"`
		Address1              string    `json:"address1"`
		Address2              string    `json:"address2"`
		City                  string    `json:"city"`
		Zip                   string    `json:"zip"`
		Province              string    `json:"province"`
		Country               string    `json:"country"`
		Phone                 string    `json:"phone"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		CountryCode           string    `json:"country_code"`
		CountryName           string    `json:"country_name"`
		ProvinceCode          string    `json:"province_code"`
		Legacy                bool      `json:"legacy"`
		Active                bool      `json:"active"`
		AdminGraphqlAPIID     string    `json:"admin_graphql_api_id"`
		LocalizedCountryName  string    `json:"localized_country_name"`
		LocalizedProvinceName string    `json:"localized_province_name"`
	} `json:"locations"`
}
type ResponseShopifyInventoryLevels struct {
	InventoryLevels []struct {
		InventoryItemID   int    `json:"inventory_item_id"`
		LocationID        int    `json:"location_id"`
		Available         int    `json:"available"`
		UpdatedAt         string `json:"updated_at"`
		AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	} `json:"inventory_levels"`
}

// iocsv.go
type ExportVariant struct {
	Sku     string `json:"sku"`
	Barcode string `json:"barcode"`
}

type ExportProduct struct {
	ID          uuid.UUID `json:"id"`
	ProductCode string    `json:"product_code"`
	Active      string    `json:"active"`
	Title       string    `json:"title"`
	BodyHTML    string    `json:"body_html"`
	Category    string    `json:"category"`
	Vendor      string    `json:"vendor"`
	ProductType string    `json:"product_type"`
}

type ImportResponse struct {
	ProcessedCounter int `json:"processed_counter"`
	FailCounter      int `json:"fail_counter"`
	ProductsAdded    int `json:"products_added"`
	ProductsUpdated  int `json:"products_updated"`
	VariantsAdded    int `json:"variants_added"`
	VariantsUpdated  int `json:"variants_updated"`
}

type CSVProduct struct {
	ProductCode  string        `csv:"product_code"`
	Active       string        `csv:"active"`
	Title        string        `csv:"title"`
	BodyHTML     string        `csv:"body_html"`
	Category     string        `csv:"category"`
	Vendor       string        `csv:"vendor"`
	ProductType  string        `csv:"product_type"`
	SKU          string        `csv:"sku"`
	Option1Name  string        `csv:"option1_name"`
	Option1Value string        `csv:"option1_value"`
	Option2Name  string        `csv:"option2_name"`
	Option2Value string        `csv:"option2_value"`
	Option3Name  string        `csv:"option3_name"`
	Option3Value string        `csv:"option3_value"`
	Barcode      string        `csv:"barcode"`
	Image1       string        `csv:"image_1"`
	Image2       string        `csv:"image_2"`
	Image3       string        `csv:"image_3"`
	Warehouses   []CSVQuantity `csv:"-"`
	Pricing      []CSVPricing  `csv:"-"`
}

type CSVQuantity struct {
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	Value     int    `json:"value"`
}

type CSVPricing struct {
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type ResponseString struct {
	Message string `json:"message"`
}

type RequestString struct {
	Message string `json:"message"`
}

// request_validation.go

type RequestWarehouseLocation struct {
	LocationID           string `json:"location_id"`
	WarehouseName        string `json:"warehouse_name"`
	ShopifyWarehouseName string `json:"shopify_warehouse_name"`
}
type RequestBodyUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
type RequestBodyProduct struct {
	ProductCode    string               `json:"product_code"`
	Title          string               `json:"title"`
	BodyHTML       string               `json:"body_html"`
	Category       string               `json:"category"`
	Vendor         string               `json:"vendor"`
	ProductType    string               `json:"product_type"`
	Variants       []RequestBodyVariant `json:"variants"`
	ProductOptions []ProductOptions     `json:"options"`
}

type RequestBodyVariant struct {
	Sku             string         `json:"sku"`
	Option1         string         `json:"option1"`
	Option2         string         `json:"option2"`
	Option3         string         `json:"option3"`
	Barcode         string         `json:"barcode"`
	VariantPricing  []VariantPrice `json:"variant_price_tiers"`
	VariantQuantity []VariantQty   `json:"variant_quantities"`
	UpdatedAt       time.Time      `json:"updated_at"`
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
	Status        string    `json:"status"`
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
	Status            string        `json:"status"`
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
	Address   []CustomerAddress `json:"addresses"`
	UpdatedAt time.Time         `json:"updated_at"`
}
type Customer struct {
	ID        uuid.UUID         `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Address   []CustomerAddress `json:"addresses"`
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
	ProductCode    string           `json:"product_code"`
	Active         string           `json:"active"`
	Title          string           `json:"title"`
	BodyHTML       string           `json:"body_html"`
	Category       string           `json:"category"`
	Vendor         string           `json:"vendor"`
	ProductType    string           `json:"product_type"`
	Variants       []ProductVariant `json:"variants"`
	ProductOptions []ProductOptions `json:"options"`
	ProductImages  []ProductImages  `json:"product_images"`
	UpdatedAt      time.Time        `json:"updated_at"`
}
type ProductOptions struct {
	Value    string `json:"value"`
	Position int    `json:"position"`
}

type ProductImages struct {
	Src       string    `json:"src"`
	Position  int       `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariant struct {
	ID              uuid.UUID      `json:"id"`
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
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type VariantQty struct {
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	Value     int    `json:"value"`
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
	Products []struct {
		ID                int64     `json:"id"`
		Title             string    `json:"title"`
		BodyHTML          string    `json:"body_html"`
		Vendor            string    `json:"vendor"`
		ProductType       string    `json:"product_type"`
		CreatedAt         time.Time `json:"created_at"`
		Handle            string    `json:"handle"`
		UpdatedAt         time.Time `json:"updated_at"`
		PublishedAt       time.Time `json:"published_at"`
		TemplateSuffix    any       `json:"template_suffix"`
		PublishedScope    string    `json:"published_scope"`
		Tags              string    `json:"tags"`
		Status            string    `json:"status"`
		AdminGraphqlAPIID string    `json:"admin_graphql_api_id"`
		Variants          []struct {
			ID                   int64     `json:"id"`
			ProductID            int64     `json:"product_id"`
			Title                string    `json:"title"`
			Price                string    `json:"price"`
			Sku                  string    `json:"sku"`
			Position             int       `json:"position"`
			InventoryPolicy      string    `json:"inventory_policy"`
			CompareAtPrice       string    `json:"compare_at_price"`
			FulfillmentService   string    `json:"fulfillment_service"`
			InventoryManagement  string    `json:"inventory_management"`
			Option1              string    `json:"option1"`
			Option2              string    `json:"option2"`
			Option3              string    `json:"option3"`
			CreatedAt            time.Time `json:"created_at"`
			UpdatedAt            time.Time `json:"updated_at"`
			Taxable              bool      `json:"taxable"`
			Barcode              string    `json:"barcode"`
			Grams                int       `json:"grams"`
			ImageID              any       `json:"image_id"`
			Weight               float64   `json:"weight"`
			WeightUnit           string    `json:"weight_unit"`
			InventoryItemID      int64     `json:"inventory_item_id"`
			InventoryQuantity    int       `json:"inventory_quantity"`
			OldInventoryQuantity int       `json:"old_inventory_quantity"`
			RequiresShipping     bool      `json:"requires_shipping"`
			AdminGraphqlAPIID    string    `json:"admin_graphql_api_id"`
		} `json:"variants"`
		Options []struct {
			ID        int64    `json:"id"`
			ProductID int64    `json:"product_id"`
			Name      string   `json:"name"`
			Position  int      `json:"position"`
			Values    []string `json:"values"`
		} `json:"options"`
		Images []struct {
			ID                int64     `json:"id"`
			Alt               any       `json:"alt"`
			Position          int       `json:"position"`
			ProductID         int64     `json:"product_id"`
			CreatedAt         time.Time `json:"created_at"`
			UpdatedAt         time.Time `json:"updated_at"`
			AdminGraphqlAPIID string    `json:"admin_graphql_api_id"`
			Width             int       `json:"width"`
			Height            int       `json:"height"`
			Src               string    `json:"src"`
			VariantIds        []any     `json:"variant_ids"`
		} `json:"images"`
		Image struct {
			ID                int64     `json:"id"`
			Alt               any       `json:"alt"`
			Position          int       `json:"position"`
			ProductID         int64     `json:"product_id"`
			CreatedAt         time.Time `json:"created_at"`
			UpdatedAt         time.Time `json:"updated_at"`
			AdminGraphqlAPIID string    `json:"admin_graphql_api_id"`
			Width             int       `json:"width"`
			Height            int       `json:"height"`
			Src               string    `json:"src"`
			VariantIds        []any     `json:"variant_ids"`
		} `json:"image"`
	} `json:"products"`
}

type ShopifyProduct struct {
	ShopifyProd `json:"product"`
}

type ShopifyProd struct {
	Title    string               `json:"title"`
	BodyHTML string               `json:"body_html"`
	Vendor   string               `json:"vendor"`
	Type     string               `json:"product_type"`
	Status   string               `json:"status"`
	Variants []ShopifyProdVariant `json:"variants"`
	Options  []ShopifyOptions     `json:"options"`
}
type ShopifyProdVariant struct {
	ID                   int64  `json:"id"`
	ProductID            int64  `json:"product_id"`
	Title                string `json:"title"`
	Price                string `json:"price"`
	Sku                  string `json:"sku"`
	Position             int    `json:"position"`
	InventoryPolicy      string `json:"inventory_policy"`
	CompareAtPrice       string `json:"compare_at_price"`
	InventoryManagement  string `json:"inventory_management"`
	Option1              string `json:"option1"`
	Option2              string `json:"option2"`
	Option3              string `json:"option3"`
	Barcode              string `json:"barcode"`
	Grams                int    `json:"grams"`
	InventoryItemID      int64  `json:"inventory_item_id"`
	InventoryQuantity    int    `json:"inventory_quantity"`
	OldInventoryQuantity int    `json:"old_inventory_quantity"`
}
type ShopifyVariant struct {
	ShopifyVar `json:"variant"`
}
type ShopifyVar struct {
	Sku                 string `json:"sku"`
	Price               string `json:"price"`
	CompareAtPrice      string `json:"compare_at_price"`
	Option1             string `json:"option1"`
	Option2             string `json:"option2"`
	Option3             string `json:"option3"`
	Barcode             string `json:"barcode"`
	InventoryManagement string `json:"inventory_management"`
}
type ShopifyOptions struct {
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Values   []string `json:"values"`
}

// shopify_settings.go

type ShopifySettings struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

type RequestSettings struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
