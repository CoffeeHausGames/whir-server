package main

import (
	"github.com/CoffeeHausGames/whir-server/app/server"
)

func main() {
	s := &server.Server{
		UseHTTP:  true,
		HTTPPort: 4444,
	}

	s.Start()
}