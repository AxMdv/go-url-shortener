// Package auth provides managing http.Cookies.
// It contains methods, such as: creating id to cookies, getting ID from Cookie, validation of cookie, etc.
package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// secretKey is key, that is used to create signature of a uuid string.
const secretKey = "abwLpqp4uKfxiQJIbfYIudou7K7qbtXE"

// struct to pass user id through request context.
type requestContextUserIDValue struct{}

// CreateIDToCookie returns id of a current user, and cookieValue - concatenation of user ID and signature.
func CreateIDToCookie() (id string, cookieValue string, err error) {

	//create uuid
	newID, err := createUUID()
	if err != nil {
		log.Print(err)
		return "", "", err
	}
	id = newID.String()
	byteID := []byte(id)

	//create signature
	sign := createSignature(byteID)

	//create final hex string for value in cookie
	cookieValue = concatIDAndSignature(byteID, sign)
	return id, cookieValue, err

}

// SetUUIDToRequestContext returns http.Request with id value in request context.
func SetUUIDToRequestContext(r *http.Request, id string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), requestContextUserIDValue{}, id))

}

// GetIDFromCookie returns user UUID from value of http.Cookie that is sent by a client.
func GetIDFromCookie(cookieValue string) (id string) {
	decodedCookieVal, _ := hex.DecodeString(cookieValue)
	id = getIDFromCookieVal(decodedCookieVal)
	return id
}

// GetUUIDFromContext returns user UUID from Context.
func GetUUIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(requestContextUserIDValue{}).(string)
	return userID
}

// ValidateCookie returns true if Cokie is valid.
func ValidateCookie(cookie *http.Cookie) (bool, error) {
	return validateID(cookie.Value)
}

// validateID returns true if cookieValue is concatenation of user UUID and valid signaure.
func validateID(cookieValue string) (bool, error) {
	if cookieValue == "" {
		return false, nil
	}
	decodedCookieVal, err := hex.DecodeString(cookieValue)
	if err != nil {
		return false, err
	}
	id := getIDFromCookieVal(decodedCookieVal)
	valid := validateIDWithSignature(id, decodedCookieVal)
	if !valid {
		return false, nil
	}
	return true, nil
}

// getIDFromCookieVal returns user ID that is the first 36 bytes of decoded cookie value.
func getIDFromCookieVal(decodedCookieVal []byte) (id string) {
	return string(decodedCookieVal[:36])
}

// validateIDWithSignature returns if the id from cookie value is valid.
// Note: we should create a signature from id and compare it with signature in request to validate ID.
func validateIDWithSignature(id string, decodedCookieVal []byte) (valid bool) {
	validSignature := string(createSignature([]byte(id)))
	requestSignature := string(decodedCookieVal[36:])
	return validSignature == requestSignature
}

// createUUID() returns new uuid
func createUUID() (uuid.UUID, error) {
	newID, err := uuid.NewV6()
	if err != nil {
		return newID, err
	}
	return newID, nil
}

// createSignature encrypts src with a secretKey
func createSignature(src []byte) (sign []byte) {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(src)
	sign = h.Sum(nil)
	return sign
}

// concatIDAndSignature returns concatenation of id and signature. Note: It should be two hex strings, NOT slices of bytes!
func concatIDAndSignature(byteID []byte, sign []byte) (resultStr string) {
	strEncodedID := hex.EncodeToString(byteID)
	strSign := hex.EncodeToString(sign)
	resultStr = strEncodedID + strSign
	return resultStr
}
