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

	version := "/v1" // Define version prefix here

	// User routes
	router.GET(version+"/user", EnvHandler.Authentication(EnvHandler.GetLoggedInUser))
	router.POST(version+"/users/signup", middleware.UrlDecode(EnvHandler.SignUp))
	router.POST(version+"/users/login", middleware.UrlDecode(EnvHandler.UserLogin))
	router.POST(version+"/users/logout", EnvHandler.Logout)


	// Business routes
	router.POST(version+"/business/login", middleware.UrlDecode(EnvHandler.BusinessLogin))
	router.POST(version+"/users/login/google", middleware.UrlDecode(EnvHandler.GoogleAuthentication))
	router.POST(version+"/business/signup", middleware.UrlDecode(EnvHandler.BusinessSignUp))
	router.POST(version+"/business", middleware.UrlDecode(EnvHandler.GetBusiness))
	router.POST(version+"/business/deal", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.AddDeal)))
	router.GET(version+"/business/deal", EnvHandler.BusinessAuthentication(EnvHandler.GetSignedInBusinessDeals))
	router.PUT(version+"/business/deal", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.UpdateDeal)))
	router.DELETE(version+"/business/deal", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.DeleteDeal)))
	router.DELETE(version+"/business/deals", EnvHandler.BusinessAuthentication(middleware.UrlDecode(EnvHandler.DeleteMultipleDeals)))

	// Token routes
	router.GET(version+"/token", EnvHandler.TokenRefresh)

	c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8081", "http://192.168.1.29:4444", "http://10.8.1.245:4444"}, 
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Authorization", "Content-Type", "Cookie-Consent"},
			ExposedHeaders: []string{"X-Auth-Token", "X-Refresh-Token"},
			AllowCredentials: true,
	})

	return c.Handler(router)
}