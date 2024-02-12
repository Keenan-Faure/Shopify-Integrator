package main

import (
	"errors"
	"integrator/internal/database"
	"log"
	"net/http"
	"objects"
	"strconv"
	"time"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Returns the results of a search query by the customer name and web code of the order

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) CustomerSearchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		search_query := c.Query("q")
		if search_query != "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "Invalid search param")
			return
		}
		customers_by_name, err := dbconfig.DB.GetCustomersByName(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		RespondWithJSON(c, http.StatusOK, customers_by_name)
	}
}

/*
Returns the customer data having the specific id

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) CustomerIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		customer_id := c.Param("id")
		err := IDValidation(customer_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		customer_uuid, err := uuid.Parse(customer_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode customer id: "+customer_id)
			return
		}
		customer, err := CompileCustomerData(dbconfig, customer_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, customer)
	}
}

/*
Returns the respective page of customer data from the database

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) CustomersHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		dbCustomers, err := dbconfig.DB.GetCustomers(c.Request.Context(), database.GetCustomersParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		customers := []objects.Customer{}
		for _, value := range dbCustomers {
			cust, err := CompileCustomerData(dbconfig, value.ID, c.Request.Context(), true)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			customers = append(customers, cust)
		}
		RespondWithJSON(c, http.StatusOK, customers)
	}
}

/*
Returns the results of a search query by the customer name and web code of the order

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OrderSearchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		search_query := c.Query("q")
		if search_query != "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "Invalid search param")
			return
		}
		customer_orders, err := dbconfig.DB.GetOrdersSearchByCustomer(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		webcode_orders, err := dbconfig.DB.GetOrdersSearchWebCode(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		RespondWithJSON(c, http.StatusOK, CompileOrderSearchResult(customer_orders, webcode_orders))
	}
}

/*
Returns the order data having the specific id

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OrderIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		order_id := c.Param("id")
		err := IDValidation(order_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		order_uuid, err := uuid.Parse(order_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode order id: "+order_id)
			return
		}
		order_data, err := CompileOrderData(dbconfig, order_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, order_data)
	}
}

/*
Returns the respective page of order data from the database

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OrdersHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		dbOrders, err := dbconfig.DB.GetOrders(c.Request.Context(), database.GetOrdersParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		orders := []objects.Order{}
		for _, value := range dbOrders {
			ord, err := CompileOrderData(dbconfig, value.ID, c.Request.Context(), true)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			orders = append(orders, ord)
		}
		RespondWithJSON(c, http.StatusOK, orders)
	}
}

/*
Filter Searches for certain products based on their vendor, product type and collection

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ProductFilterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		query_param_type := utils.ConfirmFilters(c.Query("type"))
		query_param_category := utils.ConfirmFilters(c.Query("category"))
		query_param_vendor := utils.ConfirmFilters(c.Query("vendor"))
		response, err := CompileFilterSearch(
			dbconfig,
			c.Request.Context(),
			page,
			utils.ConvertStringToLike(query_param_type),
			utils.ConvertStringToLike(query_param_category),
			utils.ConvertStringToLike(query_param_vendor),
		)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, response)
	}
}

/*
Returns the results of a search query by a product Title and SKU

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ProductSearchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		search_query := c.Query("q")
		if search_query == "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "Invalid search param")
			return
		}
		search, err := dbconfig.DB.GetProductsSearch(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		compiled, err := CompileSearchResult(dbconfig, c.Request.Context(), search)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, compiled)
	}
}

/*
Returns the respective page of product data from the database

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 404, 401, 500
*/
func (dbconfig *DbConfig) ProductsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		dbProducts, err := dbconfig.DB.GetProducts(c.Request.Context(), database.GetProductsParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		products := []objects.Product{}
		for _, value := range dbProducts {
			prod, err := CompileProductData(dbconfig, value.ID, c.Request.Context(), false)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			products = append(products, prod)
		}
		RespondWithJSON(c, http.StatusOK, products)
	}
}

/*
Returns the product data having the specific id

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 404, 401, 500
*/
func (dbconfig *DbConfig) ProductIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode product id '"+product_id+"'")
			return
		}
		product_data, err := CompileProductData(dbconfig, product_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, product_data)
	}
}

/*
Logs a user into the application. This does not set any cookies

Authorization: Required

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) LoginHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := DecodeLoginRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = UserValidation(body.Username, body.Password)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		db_user, exists, err := dbconfig.CheckUserCredentials(body, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			RespondWithError(c, http.StatusNotFound, "invalid username and password combination")
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseLogin{
			Username: db_user.Name,
			ApiKey:   db_user.ApiKey,
		})
	}
}

/*
Logs a user out of the application. If cookies are set, they will be  set to be expired

Authorization: Required

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) LogoutHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie(cookie_name); err == nil {
			value := make(map[string]string)
			if err = s.Decode(cookie_name, cookie.Value, &value); err == nil {
				// removes the cookie
				cookie := &http.Cookie{
					Name:   cookie_name,
					Value:  "",
					Secure: false,
					Path:   "/",
					MaxAge: -1,
				}
				http.SetCookie(c.Writer, cookie)
			}
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Preregisters a new user. A token is sent to the email that the user provides
Which is then used in the registration.

Authorization: None

Response-Type: application/json

Possible HTTP Codes:  200, 400, 401, 409, 500
*/
func (dbconfig *DbConfig) PreRegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := utils.LoadEnv("email")
		email_psw := utils.LoadEnv("email_psw")
		if email == "" || email_psw == "" {
			RespondWithError(c, http.StatusInternalServerError, "invalid email or email password")
			return
		}
		request_body, err := DecodePreRegisterRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = PreRegisterValidation(request_body)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		// user validation
		exists, err := dbconfig.CheckUserEmailType(email, "app")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if exists {
			RespondWithError(c, http.StatusConflict, "email '"+email+"' already exists")
			return
		}
		token_value := uuid.UUID{}
		token_value, exists, err = dbconfig.CheckTokenExists(email, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			token, err := dbconfig.DB.CreateToken(c.Request.Context(), database.CreateTokenParams{
				ID:        uuid.New(),
				Name:      request_body.Name,
				Email:     request_body.Email,
				Token:     uuid.New(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			token_value = token.Token
		}
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		err = SendEmail(token_value, request_body.Email, request_body.Name)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, objects.ResponseString{
			Message: "email sent",
		})

	}
}

/*
Registers a new user. It expects an email and a token to be passed
into the body of the request. The token will be verified to confirm if it exists internally

Authorization: None

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 409, 500
*/
func (dbconfig *DbConfig) RegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := DecodeUserRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = ValidateTokenValidation(body)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		token, err := dbconfig.DB.GetTokenValidation(c.Request.Context(), body.Email)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		request_token, err := uuid.Parse(body.Token)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode token: "+body.Token)
			return
		}
		if token.Token != request_token {
			RespondWithError(c, http.StatusNotFound, "invalid token for user")
			return
		}
		err = UserValidation(body.Name, body.Password)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		exists, err := dbconfig.CheckUserExist(body.Name, c.Request)
		if exists {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		user, err := dbconfig.DB.CreateUser(c.Request.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      body.Name,
			UserType:  "app",
			Email:     body.Email,
			Password:  body.Password,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, ConvertDatabaseToRegister(user))
	}
}

/*
Confirms if the API is ready to start accepting requests.

Authorization: None

Response-Type: application/json

Possible HTTP Codes: 200, 503
*/
func (dbconfig *DbConfig) ReadyHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if dbconfig.Valid {
			RespondWithJSON(c, http.StatusOK, gin.H{"message": "OK"})
		} else {
			RespondWithError(c, http.StatusServiceUnavailable, "Unavailable")
		}
	}
}

// Helper function
// logs all error messages in current context to stdout
// the error message in the parameters is returned over the API
// after the chain has been aborted.
func RespondWithError(c *gin.Context, http_code int, err_message string) {
	for _, err := range c.Errors {
		// TODO log previous errors from the authentication middlewares inside database table
		log.Println(err.Err.Error())
		break
	}
	c.AbortWithStatusJSON(http_code, gin.H{
		"message": err_message,
	})
}

// Helper function
// responds with a payload and http code over the API
// after sucessfully processing the request.
func RespondWithJSON(c *gin.Context, http_code int, payload any) {
	c.JSON(http_code, payload)
}

// Helper function
// appends the error to the current context and passes on to the next Middleware
// only used in the Authentication phase when one middleware auth is not met
// it uses another in the format below
// query_param -> api_key as header -> basic
// in the case that the last auth fails, it raises and error with RespondWithError
func AppendErrorNext(c *gin.Context, http_code int, err_message string) {
	c.Error(errors.New(err_message))
	c.Next()
}
