package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"errors"

	"github.com/CoffeeHausGames/whir-server/app/server/database"
	"github.com/CoffeeHausGames/whir-server/app/model"
)

// HandlerEnv is a wrapper for the genral request handling and contains a database instance
type HandlerEnv struct {
	database *database.Database
}

// NewHandlerEnv returns a new [HandlerEnv] with the specified database
func NewHandlerEnv(db *database.Database) *HandlerEnv {
	return &HandlerEnv{
		database: db,
	}
}

// WriteSuccessResponse writes a successful response to a writer. 
// It sets the HTTP status to 200 and sends a JSON-encoded response.
//
// Parameters:
// w: The http.ResponseWriter to write the response to.
// d: The data to be included in the response. This can be any type that can be JSON-encoded.
// user: The user for whom the response is being written. If sendTokens is true, the user's auth token and refresh token will be included in the response.
// sendTokens: A boolean indicating whether to include the user's auth token and refresh token in the response.
func WriteSuccessResponse(w http.ResponseWriter, r *http.Request, d interface{}, user model.UserInterface, sendTokens bool){

	if sendTokens {
			err := WriteTokenResponse(w, r, user)
			if err != nil {
					WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
					return
			}   
	}
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&model.Response{Data:d}); err != nil {
					WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
					return
	}

	log.Println("Request was a success")
}

// WriteTokenResponse writes the auth token and refresh token to a writer, either in cookies or in the headers.
// It first checks if the auth token and refresh token are not nil.
// If either of them is nil, it returns an error.
// If the user has consented to cookies (user.Cookie_consent is true), it sets the tokens in cookies.
// If the user has not consented to cookies (user.Cookie_consent is false or nil), it sets the tokens in the headers.
//
// Parameters:
// w: The http.ResponseWriter to write the response to.
// user: The user for whom the response is being written. The user's auth token and refresh token will be included in the response.
//
// Returns:
// An error if the auth token or refresh token is nil.
func WriteTokenResponse(w http.ResponseWriter, r *http.Request, user model.UserInterface) error {
	authToken := user.GetToken()
	refreshToken := user.GetRefreshToken()

	if authToken == nil || refreshToken == nil {
			return errors.New("Token or refresh token was nil")
	}

	// Check for a Cookie-Consent header in the request
	cookieConsent := r.Header.Get("Cookie-Consent")

	if cookieConsent == "true" {
			log.Println("User has consented to cookies")
			// If the user has consented to cookies, set the tokens in cookies
			http.SetCookie(w, &http.Cookie{
					Name:     "access_token",
					Value:    *authToken,
					Path:     "/",
					// Secure:  true, // Uncomment this line if you're using HTTPS
					HttpOnly: true,
					MaxAge:   3600, // 1 hour
			})
			http.SetCookie(w, &http.Cookie{
					Name:     "refresh_token",
					Value:    *refreshToken,
					Path:     "/",
					// Secure:  true, // Uncomment this line if you're using HTTPS
					HttpOnly: true,
					MaxAge:   86400, // 24 hours
			})
	} else {
			// If the user has not consented to cookies, set the tokens in the headers
			w.Header().Set("X-Auth-Token", *authToken)
			w.Header().Set("X-Refresh-Token", *refreshToken)
	}

	log.Println("Auth token and refresh token sent")
	return nil
}

// WriteRedirectResponse writes a redirect response to a writer
func WriteRedirectResponse(w http.ResponseWriter, r *http.Request, redirectURL string) {
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	log.Println("Redirecting to:", redirectURL)
}

// WriteErrorResponse writes an error code and message to a writer
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string){
	log.Println(errorMsg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(model.ErrorResponse{Status: errorCode, Name: errorMsg})
	log.Println("There was an error with the request")
}
