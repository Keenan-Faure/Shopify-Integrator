package objects

import "time"

type RequestBodyOrder struct {
	ID                    int64  `json:"id"`
	AdminGraphqlAPIID     string `json:"admin_graphql_api_id"`
	AppID                 int    `json:"app_id"`
	BrowserIP             string `json:"browser_ip"`
	BuyerAcceptsMarketing bool   `json:"buyer_accepts_marketing"`
	CancelReason          any    `json:"cancel_reason"`
	CancelledAt           any    `json:"cancelled_at"`
	CartToken             any    `json:"cart_token"`
	CheckoutID            int64  `json:"checkout_id"`
	CheckoutToken         string `json:"checkout_token"`
	ClientDetails         struct {
		AcceptLanguage any    `json:"accept_language"`
		BrowserHeight  any    `json:"browser_height"`
		BrowserIP      string `json:"browser_ip"`
		BrowserWidth   any    `json:"browser_width"`
		SessionHash    any    `json:"session_hash"`
		UserAgent      string `json:"user_agent"`
	} `json:"client_details"`
	ClosedAt                any       `json:"closed_at"`
	Confirmed               bool      `json:"confirmed"`
	ContactEmail            string    `json:"contact_email"`
	CreatedAt               time.Time `json:"created_at"`
	Currency                string    `json:"currency"`
	CurrentSubtotalPrice    string    `json:"current_subtotal_price"`
	CurrentSubtotalPriceSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"current_subtotal_price_set"`
	CurrentTotalDiscounts    string `json:"current_total_discounts"`
	CurrentTotalDiscountsSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"current_total_discounts_set"`
	CurrentTotalDutiesSet any    `json:"current_total_duties_set"`
	CurrentTotalPrice     string `json:"current_total_price"`
	CurrentTotalPriceSet  struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"current_total_price_set"`
	CurrentTotalTax    string `json:"current_total_tax"`
	CurrentTotalTaxSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"current_total_tax_set"`
	CustomerLocale string `json:"customer_locale"`
	DeviceID       any    `json:"device_id"`
	DiscountCodes  []struct {
		Code   string `json:"code"`
		Amount string `json:"amount"`
		Type   string `json:"type"`
	} `json:"discount_codes"`
	Email                  string    `json:"email"`
	EstimatedTaxes         bool      `json:"estimated_taxes"`
	FinancialStatus        string    `json:"financial_status"`
	FulfillmentStatus      any       `json:"fulfillment_status"`
	Gateway                string    `json:"gateway"`
	LandingSite            any       `json:"landing_site"`
	LandingSiteRef         any       `json:"landing_site_ref"`
	LocationID             any       `json:"location_id"`
	MerchantOfRecordAppID  any       `json:"merchant_of_record_app_id"`
	Name                   string    `json:"name"`
	Note                   any       `json:"note"`
	NoteAttributes         []any     `json:"note_attributes"`
	Number                 int       `json:"number"`
	OrderNumber            int       `json:"order_number"`
	OrderStatusURL         string    `json:"order_status_url"`
	OriginalTotalDutiesSet any       `json:"original_total_duties_set"`
	PaymentGatewayNames    []string  `json:"payment_gateway_names"`
	Phone                  string    `json:"phone"`
	PresentmentCurrency    string    `json:"presentment_currency"`
	ProcessedAt            time.Time `json:"processed_at"`
	ProcessingMethod       string    `json:"processing_method"`
	Reference              string    `json:"reference"`
	ReferringSite          any       `json:"referring_site"`
	SourceIdentifier       string    `json:"source_identifier"`
	SourceName             string    `json:"source_name"`
	SourceURL              any       `json:"source_url"`
	SubtotalPrice          string    `json:"subtotal_price"`
	SubtotalPriceSet       struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"subtotal_price_set"`
	Tags     string `json:"tags"`
	TaxLines []struct {
		Price    string  `json:"price"`
		Rate     float64 `json:"rate"`
		Title    string  `json:"title"`
		PriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"price_set"`
		ChannelLiable bool `json:"channel_liable"`
	} `json:"tax_lines"`
	TaxesIncluded     bool   `json:"taxes_included"`
	Test              bool   `json:"test"`
	Token             string `json:"token"`
	TotalDiscounts    string `json:"total_discounts"`
	TotalDiscountsSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"total_discounts_set"`
	TotalLineItemsPrice    string `json:"total_line_items_price"`
	TotalLineItemsPriceSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"total_line_items_price_set"`
	TotalOutstanding string `json:"total_outstanding"`
	TotalPrice       string `json:"total_price"`
	TotalPriceSet    struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"total_price_set"`
	TotalShippingPriceSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"total_shipping_price_set"`
	TotalTax    string `json:"total_tax"`
	TotalTaxSet struct {
		ShopMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"shop_money"`
		PresentmentMoney struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currency_code"`
		} `json:"presentment_money"`
	} `json:"total_tax_set"`
	TotalTipReceived string    `json:"total_tip_received"`
	TotalWeight      int       `json:"total_weight"`
	UpdatedAt        time.Time `json:"updated_at"`
	UserID           int64     `json:"user_id"`
	BillingAddress   struct {
		FirstName    string  `json:"first_name"`
		Address1     string  `json:"address1"`
		Phone        any     `json:"phone"`
		City         string  `json:"city"`
		Zip          string  `json:"zip"`
		Province     string  `json:"province"`
		Country      string  `json:"country"`
		LastName     string  `json:"last_name"`
		Address2     string  `json:"address2"`
		Company      string  `json:"company"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Name         string  `json:"name"`
		CountryCode  string  `json:"country_code"`
		ProvinceCode string  `json:"province_code"`
	} `json:"billing_address"`
	Customer struct {
		ID                    int64     `json:"id"`
		Email                 string    `json:"email"`
		AcceptsMarketing      bool      `json:"accepts_marketing"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		FirstName             string    `json:"first_name"`
		LastName              string    `json:"last_name"`
		State                 string    `json:"state"`
		Note                  any       `json:"note"`
		VerifiedEmail         bool      `json:"verified_email"`
		MultipassIdentifier   any       `json:"multipass_identifier"`
		TaxExempt             bool      `json:"tax_exempt"`
		Phone                 string    `json:"phone"`
		EmailMarketingConsent struct {
			State            string `json:"state"`
			OptInLevel       string `json:"opt_in_level"`
			ConsentUpdatedAt any    `json:"consent_updated_at"`
		} `json:"email_marketing_consent"`
		SmsMarketingConsent struct {
			State                string `json:"state"`
			OptInLevel           string `json:"opt_in_level"`
			ConsentUpdatedAt     any    `json:"consent_updated_at"`
			ConsentCollectedFrom string `json:"consent_collected_from"`
		} `json:"sms_marketing_consent"`
		Tags                      string    `json:"tags"`
		Currency                  string    `json:"currency"`
		AcceptsMarketingUpdatedAt time.Time `json:"accepts_marketing_updated_at"`
		MarketingOptInLevel       any       `json:"marketing_opt_in_level"`
		TaxExemptions             []any     `json:"tax_exemptions"`
		AdminGraphqlAPIID         string    `json:"admin_graphql_api_id"`
		DefaultAddress            struct {
			ID           int64  `json:"id"`
			CustomerID   int64  `json:"customer_id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Company      string `json:"company"`
			Address1     string `json:"address1"`
			Address2     string `json:"address2"`
			City         string `json:"city"`
			Province     string `json:"province"`
			Country      string `json:"country"`
			Zip          string `json:"zip"`
			Phone        string `json:"phone"`
			Name         string `json:"name"`
			ProvinceCode string `json:"province_code"`
			CountryCode  string `json:"country_code"`
			CountryName  string `json:"country_name"`
			Default      bool   `json:"default"`
		} `json:"default_address"`
	} `json:"customer"`
	DiscountApplications []struct {
		TargetType       string `json:"target_type"`
		Type             string `json:"type"`
		Value            string `json:"value"`
		ValueType        string `json:"value_type"`
		AllocationMethod string `json:"allocation_method"`
		TargetSelection  string `json:"target_selection"`
		Title            string `json:"title"`
		Description      any    `json:"description"`
	} `json:"discount_applications"`
	Fulfillments []any `json:"fulfillments"`
	LineItems    []struct {
		ID                  int64  `json:"id"`
		AdminGraphqlAPIID   string `json:"admin_graphql_api_id"`
		FulfillableQuantity int    `json:"fulfillable_quantity"`
		FulfillmentService  string `json:"fulfillment_service"`
		FulfillmentStatus   any    `json:"fulfillment_status"`
		GiftCard            bool   `json:"gift_card"`
		Grams               int    `json:"grams"`
		Name                string `json:"name"`
		Price               string `json:"price"`
		PriceSet            struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"price_set"`
		ProductExists    bool   `json:"product_exists"`
		ProductID        any    `json:"product_id"`
		Properties       []any  `json:"properties"`
		Quantity         int    `json:"quantity"`
		RequiresShipping bool   `json:"requires_shipping"`
		Sku              string `json:"sku"`
		Taxable          bool   `json:"taxable"`
		Title            string `json:"title"`
		TotalDiscount    string `json:"total_discount"`
		TotalDiscountSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"total_discount_set"`
		VariantID                  any    `json:"variant_id"`
		VariantInventoryManagement any    `json:"variant_inventory_management"`
		VariantTitle               string `json:"variant_title"`
		Vendor                     string `json:"vendor"`
		TaxLines                   []struct {
			ChannelLiable bool   `json:"channel_liable"`
			Price         string `json:"price"`
			PriceSet      struct {
				ShopMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"shop_money"`
				PresentmentMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"presentment_money"`
			} `json:"price_set"`
			Rate  float64 `json:"rate"`
			Title string  `json:"title"`
		} `json:"tax_lines"`
		Duties              []any `json:"duties"`
		DiscountAllocations []struct {
			Amount    string `json:"amount"`
			AmountSet struct {
				ShopMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"shop_money"`
				PresentmentMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"presentment_money"`
			} `json:"amount_set"`
			DiscountApplicationIndex int `json:"discount_application_index"`
		} `json:"discount_allocations"`
	} `json:"line_items"`
	PaymentTerms    any   `json:"payment_terms"`
	Refunds         []any `json:"refunds"`
	ShippingAddress struct {
		FirstName    string `json:"first_name"`
		Address1     string `json:"address1"`
		Phone        any    `json:"phone"`
		City         string `json:"city"`
		Zip          string `json:"zip"`
		Province     string `json:"province"`
		Country      string `json:"country"`
		LastName     string `json:"last_name"`
		Address2     string `json:"address2"`
		Company      string `json:"company"`
		Latitude     any    `json:"latitude"`
		Longitude    any    `json:"longitude"`
		Name         string `json:"name"`
		CountryCode  string `json:"country_code"`
		ProvinceCode string `json:"province_code"`
	} `json:"shipping_address"`
	ShippingLines []struct {
		ID                 int64  `json:"id"`
		CarrierIdentifier  string `json:"carrier_identifier"`
		Code               string `json:"code"`
		DeliveryCategory   any    `json:"delivery_category"`
		DiscountedPrice    string `json:"discounted_price"`
		DiscountedPriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"discounted_price_set"`
		Phone    any    `json:"phone"`
		Price    string `json:"price"`
		PriceSet struct {
			ShopMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"shop_money"`
			PresentmentMoney struct {
				Amount       string `json:"amount"`
				CurrencyCode string `json:"currency_code"`
			} `json:"presentment_money"`
		} `json:"price_set"`
		RequestedFulfillmentServiceID any    `json:"requested_fulfillment_service_id"`
		Source                        string `json:"source"`
		Title                         string `json:"title"`
		TaxLines                      []struct {
			ChannelLiable bool   `json:"channel_liable"`
			Price         string `json:"price"`
			PriceSet      struct {
				ShopMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"shop_money"`
				PresentmentMoney struct {
					Amount       string `json:"amount"`
					CurrencyCode string `json:"currency_code"`
				} `json:"presentment_money"`
			} `json:"price_set"`
			Rate  float64 `json:"rate"`
			Title string  `json:"title"`
		} `json:"tax_lines"`
		DiscountAllocations []any `json:"discount_allocations"`
	} `json:"shipping_lines"`
}
