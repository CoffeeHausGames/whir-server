package router

import (
	"github.com/lucas-kern/tower-of-babel_server/app/router/handlers"
	"github.com/lucas-kern/tower-of-babel_server/app/router/handlers/middleware"
	"github.com/lucas-kern/tower-of-babel_server/app/server/database"
	"github.com/julienschmidt/httprouter"
)

// This registers all our routes and can wrap them in middle ware for auth and other items
// Returns the router with paths and handlers
func GetRouter(db *database.Database) *httprouter.Router{
	EnvHandler := handlers.NewHandlerEnv(db)
	router := httprouter.New()
	router.POST("/bases/place", middleware.ParseFormData(EnvHandler.Authentication(EnvHandler.PlaceBuilding)))
	router.POST("/users/signup", middleware.ParseFormData(EnvHandler.SignUp))
	router.POST("/users/login", middleware.ParseFormData(EnvHandler.Login))
	router.GET("/token", EnvHandler.TokenRefresh)
	return router
}
