package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Returns the value of the environment variable
func LoadEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load .env")
	}
	value := os.Getenv(strings.ToUpper(key))
	return value
}

// Authorization: ApiKey <key>
func ExtractAPIKey(authString string) (string, error) {
	if authString == "" {
		return "", errors.New("no Authorization found in header")
	}
	if len(authString) <= 7 {
		return "", errors.New("malformed Auth Header")
	}
	if authString[0:6] != "ApiKey" {
		return "", errors.New("malformed second part of authentication")
	}
	return authString[7:], nil
}

// Convert string to LIKE (%) sql format
func ConvertStringToLike(value string) string {
	return "%" + value + "%"
}

// Returns the filter value if valid
func ConfirmFilters(filter string) string {
	if filter != "" || len(filter) > 0 {
		return filter
	}
	return ""
}

// converts a string to a sql.NullString object
func ConvertStringToSQL(description string) sql.NullString {
	if description == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: description,
		Valid:  true,
	}
}

// converts a string to a sql.NullInt32 object
func ConvertIntToSQL(value int) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(value),
		Valid: true,
	}
}

// Checks if a variable is set (string)
func IssetString(variable string) string {
	if variable != "" || len(variable) != 0 {
		return variable
	}
	return ""
}

// Checks if a variable is set (string)
func IssetInt(variable string) int {
	if variable != "" || len(variable) != 0 {
		integer, err := strconv.Atoi(variable)
		if err != nil {
			return 0
		}
		return integer
	}
	return 0
}

// Extracts the PID from the Shopify Response
func ExtractPID(id string) string {
	if id == "" || len(id) == 0 {
		return ""
	}
	if len(id) > 22 {
		return id[22:]
	}
	return ""
}

// Extracts the VID from the Shopify Response
func ExtractVID(id string) string {
	if id == "" || len(id) == 0 {
		return ""
	}
	if len(id) > 29 {
		return id[29:]
	}
	return ""
}

// Returns the next URL in the Shopify Response header
func GetNextURL(next string) string {
	result := strings.Split(next, ", ")
	if len(result) == 0 {
		result = strings.Split(next, "; ")
		if len(result) == 0 {
			return ""
		} else {
			next = strings.TrimSuffix(strings.TrimPrefix(result[0], "<"), ">")
			result = strings.Split(next, "?")
			if len(result) > 0 {
				next = "products.json?" + result[1]
				return next
			}
			return ""
		}
	}
	result = strings.Split(result[len(result)-1], "; ")
	next = strings.TrimSuffix(strings.TrimPrefix(result[0], "<"), ">")
	result = strings.Split(next, "?")
	if len(result) > 0 {
		next = "products.json?" + result[1]
		return next
	}
	return ""
}

// Generates a random password
func RandStringBytes(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
