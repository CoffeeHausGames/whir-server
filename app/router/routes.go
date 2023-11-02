package router

import (
	"github.com/CoffeeHausGames/whir-server/app/router/handlers"
	"github.com/CoffeeHausGames/whir-server/app/server/database"
	"github.com/CoffeeHausGames/whir-server/app/router/handlers/middleware"
	"github.com/julienschmidt/httprouter"
)

// This registers all our routes and can wrap them in middle ware for auth and other items
// Returns the router with paths and handlers
func GetRouter(db *database.Database) *httprouter.Router{
	EnvHandler := handlers.NewHandlerEnv(db)
	router := httprouter.New()
	router.POST("/users/signup", middleware.UrlDecode(EnvHandler.SignUp))
	router.POST("/business/signup", middleware.UrlDecode(EnvHandler.BusinessSignUp))
	router.POST("/users/login", middleware.UrlDecode(EnvHandler.Login))
	router.POST("/business/login", middleware.UrlDecode(EnvHandler.BusinessLogin))
	router.GET("/token", middleware.UrlDecode(EnvHandler.TokenRefresh))
	router.POST("/business", middleware.UrlDecode(EnvHandler.GetBusiness))
	return router
}
