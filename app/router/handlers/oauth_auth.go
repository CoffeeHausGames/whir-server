package handlers

import (
	"golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "io/ioutil"
	"context"
	"google.golang.org/api/idtoken"
	"net/http"
	"os"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"github.com/CoffeeHausGames/whir-server/app/model"
	"github.com/CoffeeHausGames/whir-server/app/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// func (env *HandlerEnv) GoogleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	redirectURL := os.Getenv("REDIRECT_URL")

// 	err := r.ParseForm()
// 	if err != nil {
// 			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
// 			return
// 	}
// 	fmt.Printf("HTTP Request: %+v\n", r)

	// idToken := r.PostFormValue("credential")

	// payload, err := idtoken.Validate(oauth2.NoContext, idToken, "")
	// if err != nil {
	// 		http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
	// 		return
	// }
	
	// email, ok := payload.Claims["email"].(string)
	// if !ok {
	// 		http.Error(w, "Failed to get email from ID token", http.StatusInternalServerError)
	// 		return
	// }



	// Redirect to the React app's homepage
	// http.Redirect(w, r, redirectURL, http.StatusSeeOther)
// }

func (env *HandlerEnv) GoogleAuthentication(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	redirectURL := os.Getenv("REDIRECT_URL")

	var userCollection model.Collection = env.database.GetUsers()
	var user model.User
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.ParseForm()
	if err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
	}

	// Get the ID token from the request
	idToken := r.PostFormValue("credential")

	// Validate the ID token
	payload, err := idtoken.Validate(oauth2.NoContext, idToken, "")
	if err != nil {
			http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
			return
	}

	// Get the user's email from the ID token payload
	email, ok := payload.Claims["email"].(string)
	if !ok {
			http.Error(w, "Failed to get email from ID token", http.StatusInternalServerError)
			return
	}


	// Check if the user is already signed up
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
			http.Error(w, "Failed to check if user is already signed up", http.StatusInternalServerError)
			return
	}

	var userWrapper *model.UserWrapper

	if count > 0 {
			foundUser := new(model.User)
	
			fmt.Println("User is already signed up")
			foundUserInterface, err := performGoogleLogin(r.Context(), email, userCollection, foundUser)
			if err != nil {
					fmt.Println("performLogin failed")
					WriteErrorResponse(w, 401, "There was an error logging in")
					return
			}
	
			foundUser, ok := foundUserInterface.(*model.User)
			if !ok {
					fmt.Println("Wrong user type found")
					WriteErrorResponse(w, 401, "There was an error logging in")
					return
			}
	
			// Convert the found user to a UserWrapper
			userWrapper = model.NewUser(foundUser)
	} else {
			firstName, ok := payload.Claims["given_name"].(string)
			if !ok {
					// Handle the case where the "given_name" field is not a string or is not present
			}
			// The user is not signed up, create a new user
			user.First_name = &firstName
			user.Email = &email
			user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			token, refreshToken, _ := auth.GenerateAllTokens(*user.Email, *user.First_name, user.Last_name, user.ID.Hex())
			user.Token = &token
			user.Refresh_token = &refreshToken
	
			_, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
					http.Error(w, "Failed to create user", http.StatusInternalServerError)
					return
			}
	
			// Convert the new user to a UserWrapper
			userWrapper = model.NewUser(&user)
	}

	// Create cookies for the access token, refresh token, and user ID
	http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    *userWrapper.Token,
			Path:     "/",
			HttpOnly: true,
			// Secure:   true, // Uncomment this line if you're using HTTPS
	})

	http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    *userWrapper.Refresh_token,
			Path:     "/",
			HttpOnly: true,
			// Secure:   true, // Uncomment this line if you're using HTTPS
	})

	WriteRedirectResponse(w, r, redirectURL)
}
