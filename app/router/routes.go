package router

import (
	"github.com/CoffeeHausGames/whir-server/app/router/handlers"
	"github.com/CoffeeHausGames/whir-server/app/server/database"
	"github.com/CoffeeHausGames/whir-server/app/router/handlers/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
)

func GetRouter(db *database.Database) http.Handler {
	EnvHandler := handlers.NewHandlerEnv(db)
	router := httprouter.New()
	router.POST("/users/signup", middleware.UrlDecode(EnvHandler.SignUp))
	router.POST("/business/signup", middleware.UrlDecode(EnvHandler.BusinessSignUp))
	router.POST("/users/login", middleware.UrlDecode(EnvHandler.Login))
	router.POST("/business/login", middleware.UrlDecode(EnvHandler.BusinessLogin))
	router.GET("/token", middleware.UrlDecode(EnvHandler.TokenRefresh))
	router.POST("/business", middleware.UrlDecode(EnvHandler.GetBusiness))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // your origin here
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	return c.Handler(router)
}