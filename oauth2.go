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

// Hash keys should be at least 32 bytes long
var hashKey = []byte(securecookie.GenerateRandomKey(64))

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
		"openid"},
	Endpoint: google.Endpoint,
}

// GET /api/google/cookie/login
func (dbconfig *DbConfig) OAuthGoogleCookie(w http.ResponseWriter, r *http.Request) {
	// retrieve the cookie and check if the user record exists in the database...
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

	// creates db user
	db_user, err := dbconfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      oauth_data.GivenName + " " + oauth_data.FamilyName,
		Email:     oauth_data.Email,
		Password:  utils.RandStringBytes(10), // generates a random password, but user should never login with password though
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// creates the oauth record inside the database
	// TODO should probably check if the record already exists inside the db
	db_oauth_record, err := dbconfig.DB.GetUserByGoogleID(r.Context(), oauth_data.ID)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if db_oauth_record.GoogleID == oauth_data.ID {
		RespondWithError(w, http.StatusConflict, "user already registered")
		return
	}
	oauth_record, err := dbconfig.DB.CreateOAuthRecord(r.Context(), database.CreateOAuthRecordParams{
		ID:          uuid.New(),
		UserID:      db_user.ID,
		CookieToken: hashKey,
		GoogleID:    oauth_data.ID,
		Email:       oauth_data.Email,
		Picture:     utils.ConvertStringToSQL(oauth_data.Picture),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// create cookie containing the cookie_secret
	value := map[string]string{
		"si_googleauth": oauth_record.CookieSecret,
	}
	if encoded, err := s.Encode("si_googleauth", value); err == nil {
		cookie := &http.Cookie{
			Name:     "si_googleauth",
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
	// create an oauth table record with a users record
	// when logging in we inner join the two tables based on the user_id
	// this way the user still uses's an API Key to access the api's resources

	// redirect back to the application login screen where the user logins in automatically
	// using the new credentials
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
