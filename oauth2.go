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

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Hash keys should be at least 32 bytes long.
// Hash key is a constant
var hashKey = []byte("nSTDTVzvNdflcOlclhuaSFJfrkzKdBJjKTeRAhTVVFyiHqrUcNgvmhfXAvlGYpmv")

/*
General name of the cookie of the application for google accounts.
If the user logs in with another account the cookie should be the same name,
just updated
*/
const cookie_name = "si_googleauth"

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var s = securecookie.New(hashKey, nil)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

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

// GET /api/google/oauth2/login
func (dbconfig *DbConfig) OAuthGoogleOAuth(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(cookie_name); err == nil {
		value := make(map[string]string)
		if err = s.Decode(cookie_name, cookie.Value, &value); err == nil {
			// retrieve the cookie value from the map and search it's value inside the DB
			// to confirm if the value is correct.
			cookie_secret := value[cookie_name]
			user, err := dbconfig.DB.GetApiKeyByCookieSecret(r.Context(), cookie_secret)
			if err != nil {
				RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}
			// returns the API Key
			RespondWithJSON(w, http.StatusOK, user.ApiKey)
			return
		} else {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
	} else {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
}

// GET /api/google/login
func (dbconfig *DbConfig) OAuthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// GET /api/google/callback
func (dbconfig *DbConfig) OAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// convert to struct
	oauth_data := objects.ResponseOAuthGoogle{}
	err = json.Unmarshal(data, &oauth_data)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// creates the oauth record inside the database
	db_oauth_record, err := dbconfig.DB.GetUserByGoogleID(r.Context(), oauth_data.ID)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
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
			cookie := &http.Cookie{
				Name:     cookie_name,
				Value:    encoded,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, r, "http://localhost:3000/", http.StatusSeeOther)
		return
	}
	// creates db user
	db_user, err := dbconfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      oauth_data.GivenName + " " + oauth_data.FamilyName,
		Email:     oauth_data.Email,
		Password:  utils.RandStringBytes(20), // generates a random password, but user should never login with password though
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	oauth_record, err := dbconfig.DB.CreateOAuthRecord(r.Context(), database.CreateOAuthRecordParams{
		ID:        uuid.New(),
		UserID:    db_user.ID,
		GoogleID:  oauth_data.ID,
		Email:     oauth_data.Email,
		Picture:   utils.ConvertStringToSQL(oauth_data.Picture),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// create cookie containing the cookie_secret
	value := map[string]string{
		cookie_name: oauth_record.CookieSecret,
	}
	if encoded, err := s.Encode(cookie_name, value); err == nil {
		cookie := &http.Cookie{
			Name:     cookie_name,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
	// redirect back to the application login screen where the user logins in automatically
	// using the new credentials
	http.Redirect(w, r, "http://localhost:3000/", http.StatusSeeOther)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

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
