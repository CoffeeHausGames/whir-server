package handlers

import (
    "context"
    "log"
    "net/http"
    "fmt"
    "encoding/json"

    "github.com/julienschmidt/httprouter"
    "github.com/CoffeeHausGames/whir-server/app/auth"
    "github.com/CoffeeHausGames/whir-server/app/model"
    "go.mongodb.org/mongo-driver/bson"
)

// A middleware that will take a token from the header and ensure this user is valid
// Authentication validates tokens from the Authorization header or cookies
func (env *HandlerEnv) Authentication(n httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        var userCollection model.Collection = env.database.GetUsers()

        // Try to authenticate from cookies first
        jwtCookie, err := r.Cookie("access_token")
        if err == nil {
            clientToken := jwtCookie.Value
            claims, err := auth.ValidateToken(userCollection, clientToken)
            if err == nil {
                // Store the claims in the request context
                ctx := context.WithValue(r.Context(), "claims", claims)

                // Create a new request with the updated context
                r = r.WithContext(ctx)

                // Call the next handler with the updated request
                n(w, r, ps)
                return
            }
        }

        // If cookie authentication failed, try to authenticate from the Authorization header
        clientToken := r.Header.Get("Authorization")
        if clientToken != "" {
            claims, err := auth.ValidateToken(userCollection, clientToken)
            if err == nil {
                log.Printf("Authenticated user from header")
                // Store the claims in the request context
                ctx := context.WithValue(r.Context(), "claims", claims)

                // Create a new request with the updated context
                r = r.WithContext(ctx)

                // Call the next handler with the updated request
                n(w, r, ps)
                return
            }
        }

        log.Printf("Failed to authenticate user")
        // If both methods failed, return an HTTP 401 Unauthorized error
        WriteErrorResponse(w, http.StatusUnauthorized, "Failed to authenticate user")
    }
}

func (env *HandlerEnv) BusinessAuthentication(n httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        var userCollection model.Collection = env.database.GetBusinesses()

        // Try to authenticate from cookies first
        jwtCookie, err := r.Cookie("access_token")
        if err == nil {
            clientToken := jwtCookie.Value
            claims, err := auth.ValidateToken(userCollection, clientToken)
            if err == nil {
                // Store the claims in the request context
                ctx := context.WithValue(r.Context(), "claims", claims)

                // Create a new request with the updated context
                r = r.WithContext(ctx)

                // Call the next handler with the updated request
                n(w, r, ps)
                return
            }
        }

        // If cookie authentication failed, try to authenticate from the Authorization header
        clientToken := r.Header.Get("Authorization")
        if clientToken != "" {
            claims, err := auth.ValidateToken(userCollection, clientToken)
            if err == nil {
                log.Printf("Authenticated user from header")
                // Store the claims in the request context
                ctx := context.WithValue(r.Context(), "claims", claims)

                // Create a new request with the updated context
                r = r.WithContext(ctx)

                // Call the next handler with the updated request
                n(w, r, ps)
                return
            }
        }

        log.Printf("Failed to authenticate user")
        // If both methods failed, return an HTTP 401 Unauthorized error
        WriteErrorResponse(w, http.StatusUnauthorized, "Failed to authenticate user")
    }
}

func performLogin(ctx context.Context, userCollection model.Collection, user model.UserInterface, foundUser model.UserInterface) (model.UserInterface, error) {
    // pull the URL-decoded body from the context (comes from url_decoder middleware)
    decodedData := ctx.Value("body").(string)

    err := json.Unmarshal([]byte(decodedData), &user)
    if err != nil {
        log.Println(err)
        return nil, fmt.Errorf("There was an error with the client request")
    }

    err = validate.Var(user.GetEmail(), "required,email")
    if err != nil {
        log.Println(err)
        return nil, fmt.Errorf("There was an error with user validation")
    }
 
    //TODO sanitize input before Finding in DB to avoid NoSQL injection
    err = userCollection.FindOne(foundUser, ctx, bson.M{"email": *user.GetEmail()})
    if err != nil {
        log.Println(err)
        return nil, fmt.Errorf("There was an error connecting with the server")
    }

    passwordIsValid, msg := VerifyPassword(*user.GetPassword(), *foundUser.GetPassword())
    if passwordIsValid != true {
        log.Println(msg)
        return nil, fmt.Errorf("The username or password is incorrect")
    }

    token, refreshToken, _ := auth.GenerateAllTokens(*foundUser.GetEmail(), *foundUser.GetFirstName(), foundUser.GetLastName(), foundUser.GetID().Hex())

    auth.UpdateAllTokens(userCollection, token, refreshToken, foundUser.GetID().Hex())
    foundUser.SetToken(token)
    foundUser.SetRefreshToken(refreshToken)

    return foundUser, nil
}

func performGoogleLogin(ctx context.Context, email string, userCollection model.Collection, foundUser model.UserInterface) (model.UserInterface, error) {
 
    //TODO sanitize input before Finding in DB to avoid NoSQL injection
    err := userCollection.FindOne(foundUser, ctx, bson.M{"email": email})
    if err != nil {
        log.Println(err)
        return nil, fmt.Errorf("There was an error connecting with the server")
    }

    token, refreshToken, _ := auth.GenerateAllTokens(*foundUser.GetEmail(), *foundUser.GetFirstName(), foundUser.GetLastName(), foundUser.GetID().Hex())

    auth.UpdateAllTokens(userCollection, token, refreshToken, foundUser.GetID().Hex())
    foundUser.SetToken(token)
    foundUser.SetRefreshToken(refreshToken)

    return foundUser, nil
}