package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lucas-kern/tower-of-babel_server/app/server/database"
	"github.com/lucas-kern/tower-of-babel_server/app/model"
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

// WriteSuccessResponse Write a successful response to a writer 
func WriteSuccessResponse(w http.ResponseWriter, d interface{}){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&model.Response{Data:d}); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	log.Println("Request was a success")
}

// WriteErrorResponse writes an error code and message to a writer
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string){
	log.Println(errorMsg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(model.ErrorResponse{Status: errorCode, Name: errorMsg})
	log.Println("There was an error with the request")
}
