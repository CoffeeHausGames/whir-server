package main

import "github.com/lucas-kern/tower-of-babel_server/app/server"

func main() {
	s := &server.Server{
		UseHTTP:  true,
		HTTPPort: 4444,
	}

	s.Start()
}