package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"io"
	"log"
	"net/http"
	"objects"
	"time"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/*
General name of the cookie of the application for google accounts.
If the user logs in with another account the cookie should be the same name,
just updated and overwritten
*/
const cookie_name = "si_googleauth"

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var s = securecookie.New(hashKey, nil)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// // Hash keys should be at least 32 bytes long.
// // Hash key is a constant
var hashKey = []byte("nSTDTVzvNdflcOlclhuaSFJfrkzKdBJjKTeRAhTVVFyiHqrUcNgvmhfXAvlGYpmv")

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/api/google/callback",
	ClientID:     utils.LoadEnv("OAUTH_CLIENT_ID"),
	ClientSecret: utils.LoadEnv("OAUTH_SECRET"),
	// scopes on which to retrieve the userinfo from google api
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"openid",
	},
	Endpoint: google.Endpoint,
}

/*
Callback endpoint for the OAuth2

Route: /api/google/callback

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 303, 307, 400, 404
*/
func (dbconfig *DbConfig) OAuthGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read oauthState from Cookie
		oauthState, _ := c.Cookie("oauthstate")
		if c.Request.FormValue("state") != oauthState {
			log.Println("invalid oauth google state")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		data, err := getUserDataFromGoogle(c.Request.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		// convert to struct
		oauth_data := objects.ResponseOAuthGoogle{}
		err = json.Unmarshal(data, &oauth_data)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// creates the oauth record inside the database
		db_oauth_record, err := dbconfig.DB.GetUserByGoogleID(c.Request.Context(), oauth_data.ID)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		if db_oauth_record.GoogleID == oauth_data.ID {
			// If the user already registers, we create a new cookie
			// and then redirect to the dashboard
			value := map[string]string{
				cookie_name: db_oauth_record.CookieSecret,
			}
			if encoded, err := s.Encode(cookie_name, value); err == nil {
				c.SetCookie(cookie_name, encoded, 0, "/", "", false, false)
			}
			c.Redirect(http.StatusSeeOther, "http://localhost:3000/")
			return
		}
		// user validation
		exists, err := dbconfig.CheckUserEmailType(oauth_data.Email, "google")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if exists {
			RespondWithError(c, http.StatusConflict, "email '"+oauth_data.Email+"' already exists")
			return
		}
		// creates db user
		db_user, err := dbconfig.DB.CreateUser(c.Request.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      oauth_data.GivenName + " " + oauth_data.FamilyName,
			UserType:  "google",
			Email:     oauth_data.Email,
			Password:  utils.RandStringBytes(20),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		oauth_record, err := dbconfig.DB.CreateOAuthRecord(c.Request.Context(), database.CreateOAuthRecordParams{
			ID:        uuid.New(),
			UserID:    db_user.ID,
			GoogleID:  oauth_data.ID,
			Email:     oauth_data.Email,
			Picture:   utils.ConvertStringToSQL(oauth_data.Picture),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// create cookie containing the cookie_secret
		value := map[string]string{
			cookie_name: oauth_record.CookieSecret,
		}
		if encoded, err := s.Encode(cookie_name, value); err == nil {
			c.SetCookie(cookie_name, encoded, 0, "/", "", false, false)
		}
		// redirect back to the application login screen where the user logins in automatically
		// using the new credentials
		c.Redirect(http.StatusSeeOther, "http://localhost:3000/")
	}
}

/*
Logs into Google and redirects the user thereafter.

Route: /api/google/oauth2/login

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 307, 400, 404
*/
func (dbconfig *DbConfig) OAuthGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create oauthState cookie
		oauthState := generateStateOauthCookie(c.Writer)
		/*
			AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
			validate that it matches the the state query parameter on your redirect callback.
		*/
		u := googleOauthConfig.AuthCodeURL(oauthState)
		c.Redirect(http.StatusTemporaryRedirect, u)
	}
}

/*
Initiates the OAuth2 login authorization with google

Route: /api/google/oauth2/login

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OAuthGoogleOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie(cookie_name); err == nil {
			value := make(map[string]string)
			if err = s.Decode(cookie_name, cookie, &value); err == nil {
				cookie_secret := value[cookie_name]
				user, err := dbconfig.DB.GetApiKeyByCookieSecret(c.Request.Context(), cookie_secret)
				if err != nil {
					if err.Error() != "sql: no rows in result set" {
						RespondWithError(c, http.StatusUnauthorized, err.Error())
						return
					}
					RespondWithError(c, http.StatusUnauthorized, err.Error())
					return
				}
				RespondWithJSON(c, http.StatusOK, objects.ResponseLogin{
					Username: user.Name,
					ApiKey:   user.ApiKey,
				})
				return
			} else {
				RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
		} else {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
	}
}

/*
Generates a random state token to be used with the cookie to prevent CSRF attacks
*/
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err.Error())
	}
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

/*
Retrieves the data from google using the code
*/
func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
