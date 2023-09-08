package main

import (
	"api"
	"database/sql"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"log"
	"net/http"
	"objects"
	"strconv"
	"time"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

// POST /api/customers/
func (dbconfig *DbConfig) PostCustomerHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {

}

// POST /api/orders?token={{token}}&api_key={{key}}
// ngrok exposed url
func (dbconfig *DbConfig) PostOrderHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {

}

// POST /api/products/
func (dbconfig *DbConfig) PostProductHandle(w http.ResponseWriter, r *http.Request) {
	params, err := DecodeProductRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if ProductValidation(params) != nil {
		RespondWithError(w, http.StatusBadRequest, "data validation error")
		return
	}
	err = ProductValidation(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}
	err = ValidateDuplicateOption(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}
	err = ValidateDuplicateSKU(params, dbconfig, r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}
	err = DuplicateOptionValues(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// add product to database
	_, err = dbconfig.DB.CreateProduct(r.Context(), database.CreateProductParams{
		Active:      "1",
		Title:       utils.ConvertStringToSQL(params.Title),
		BodyHtml:    utils.ConvertStringToSQL(params.BodyHTML),
		Category:    utils.ConvertStringToSQL(params.Category),
		Vendor:      utils.ConvertStringToSQL(params.Vendor),
		ProductType: utils.ConvertStringToSQL(params.ProductType),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	fmt.Println(err)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// respond with success message
}

// GET /api/customers/search?q=value
func (dbconfig *DbConfig) CustomerSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query != "" || len(search_query) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid search param")
		return
	}
	customers_by_name, err := dbconfig.DB.GetCustomersByName(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondWithJSON(w, http.StatusOK, customers_by_name)
}

// GET /api/customers/{id}
func (dbconfig *DbConfig) CustomerHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	customer_id := chi.URLParam(r, "id")
	err := IDValidation(customer_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	customer_id_byte, err := uuid.Parse(customer_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode feed_id: "+customer_id)
		return
	}
	customer, err := CompileCustomerData(dbconfig, customer_id_byte, r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, customer)
}

// GET /api/customers?page=1
func (dbconfig *DbConfig) CustomersHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		log.Println("Error decoding page param:", err)
	}
	customers, err := dbconfig.DB.GetCustomers(r.Context(), database.GetCustomersParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, customers)
}

// GET /api/orders/search?q=value
func (dbconfig *DbConfig) OrderSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query != "" || len(search_query) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid search param")
		return
	}
	customer_orders, err := dbconfig.DB.GetOrdersSearchByCustomer(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	webcode_orders, err := dbconfig.DB.GetOrdersSearchWebCode(r.Context(), search_query)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondWithJSON(w, http.StatusOK, CompileOrderSearchResult(customer_orders, webcode_orders))
}

// GET /api/orders/{id}
func (dbconfig *DbConfig) OrderHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	order_id := chi.URLParam(r, "id")
	err := IDValidation(order_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	order_id_byte := []byte(order_id)
	order_data, err := CompileOrderData(dbconfig, order_id_byte, r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, order_data)
}

// GET /api/orders?page=1
func (dbconfig *DbConfig) OrdersHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		log.Println("Error decoding page param:", err)
	}
	dbOrders, err := dbconfig.DB.GetOrders(r.Context(), database.GetOrdersParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, dbOrders)
}

// GET /api/products/filter?data=value&page=1
func (dbconfig *DbConfig) ProductFilterHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		log.Println("Error decoding page param:", err)
	}
	query_param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
	query_param_category := utils.ConfirmFilters(r.URL.Query().Get("category"))
	query_param_vendor := utils.ConfirmFilters(r.URL.Query().Get("vendor"))
	response, err := CompileFilterSearch(dbconfig, r, page, query_param_type, query_param_category, query_param_vendor)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// GET /api/products/search?q=value
func (dbconfig *DbConfig) ProductSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query != "" || len(search_query) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid search param")
		return
	}
	sku_search, err := dbconfig.DB.GetProductsSearchSKU(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	title_search, err := dbconfig.DB.GetProductsSearchTitle(r.Context(), sql.NullString{
		String: utils.ConvertStringToLike(search_query),
		Valid:  true,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, CompileSearchResult(sku_search, title_search))
}

// GET /api/products/{id}
func (dbconfig *DbConfig) ProductHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	product_id := chi.URLParam(r, "id")
	err := IDValidation(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	product_id_byte := []byte(product_id)
	product_data, err := CompileProductData(dbconfig, product_id_byte, r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, product_data)
}

// GET /api/products?page=1
func (dbconfig *DbConfig) ProductsHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		log.Println("Error decoding page param:", err)
	}
	dbProducts, err := dbconfig.DB.GetProducts(r.Context(), database.GetProductsParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, dbProducts)
}

// POST /api/login
func (dbconfig *DbConfig) LoginHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/register
func (dbconfig *DbConfig) RegisterHandle(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeUserRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if UserValidation(body) != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	exists, err := dbconfig.CheckUserExist(body.Name, r)
	if exists {
		RespondWithError(w, http.StatusConflict, err.Error())
		return
	}
	user, err := dbconfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:      body.Name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, user)
}

// GET /api/endpoints
func (dbconfig *DbConfig) EndpointsHandle(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, api.Endpoints())
}

// puppy function for a handler
//func Handle{{Name}} {
// 1.) decode params
// 2.) validate requestBody
// 3.) Check for duplicates in db
// 4.) (fi) create/insert new record
// 4.) (if) return record
// 5.) (else) return error
//}

// GET /api/ready
func (dbconfig *DbConfig) ReadyHandle(w http.ResponseWriter, r *http.Request) {
	if dbconfig.Valid {
		RespondWithJSON(w, 200, objects.ResponseString{
			Status: "OK",
		})
	} else {
		RespondWithJSON(w, 503, objects.ResponseString{
			Status: "Error",
		})
	}
}

// JSON helper functions
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}

// Middleware that determines which headers, http methods and orgins are allowed
func MiddleWare() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
}
