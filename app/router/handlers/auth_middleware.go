package handlers

import (
    "context"
		"log"
    "net/http"

		"github.com/julienschmidt/httprouter"
		"github.com/lucas-kern/tower-of-babel_server/app/auth"
		"github.com/lucas-kern/tower-of-babel_server/app/model"
)

// A middleware that will take a token from the header and ensure this user is valid
// Authentication validates token
func (env *HandlerEnv) Authentication(n httprouter.Handle) httprouter.Handle {
	var userCollection model.Collection = env.database.GetUsers()
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        clientToken := r.Header.Get("Authorization")
        if clientToken == "" {
						log.Printf("There is no authorization token")

            return
				}

        claims, err := auth.ValidateToken(userCollection, clientToken)
        if err != "" {
						log.Panic(err)
            return
				}
				
        // Store the claims in the request context
        ctx := context.WithValue(r.Context(), "claims", claims)

        // Create a new request with the updated context
        r = r.WithContext(ctx)

        // Call the next handler with the updated request
        n(w, r, ps)
    }
}