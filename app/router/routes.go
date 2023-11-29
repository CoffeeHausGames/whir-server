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

	// User routes
	router.GET("/users", EnvHandler.Authentication(EnvHandler.GetUser))
	router.POST("/users/signup", middleware.UrlDecode(EnvHandler.SignUp))
	router.POST("/users/login", middleware.UrlDecode(EnvHandler.Login))

	// Business routes
	router.POST("/business/login", middleware.UrlDecode(EnvHandler.BusinessLogin))
	router.POST("/business/signup", middleware.UrlDecode(EnvHandler.BusinessSignUp))
	router.POST("/business", middleware.UrlDecode(EnvHandler.GetBusiness))
	router.POST("/business/deal", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.AddDeal)))
	router.GET("/business/deal", EnvHandler.BusinessAuthentication(EnvHandler.GetSignedInBusinessDeals))
	router.PUT("/business/deal", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.UpdateDeal)))
	router.DELETE("/business/deal", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.DeleteDeal)))
	router.DELETE("/business/deals", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.DeleteMultipleDeals)))
	// make route to delete multiple deals

	
	// Token routes
	router.GET("/token", EnvHandler.TokenRefresh)
	



	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8081", "http://192.168.1.29:4444", "http://10.8.1.245:4444"}, // your origin here
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	return c.Handler(router)
}