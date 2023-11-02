package server

import (
	"log"
	"net/http"
	"strconv"
	"github.com/lucas-kern/tower-of-babel_server/app/router"
	"github.com/lucas-kern/tower-of-babel_server/app/server/database"
)

/**TODO:
- Include logger
- Enable HTTPS
*/

// Server represents the server which listens for connections when started
type Server struct {
	Hostname  string `json:"hostname"`  // Server name
	UseHTTP   bool   `json:"UseHTTP"`   // Listen on HTTP
	UseHTTPS  bool   `json:"UseHTTPS"`  // Listen on HTTPS
	HTTPPort  int    `json:"HTTPPort"`  // HTTP port
	HTTPSPort int    `json:"HTTPSPort"` // HTTPS port
	CertFile  string `json:"CertFile"`  // HTTPS certificate
	KeyFile   string `json:"KeyFile"`   // HTTPS private key
	Handler   http.Handler
}

//TODO create server creation from config
// New creates a new server using the [config] map
// func (s *Server) New(config map[string]string) {
// 	// config["o"] = "p"
// 	s.Handler = routes.Load()
// 	// s.HTTPPort = 4444
// 	// s.UseHTTP = true
// }

// Start fires a listener and starts the server on the specified port
// using HTTPS if [Server.UseHTTPS] is true else it uses HTTP
// It creates the database connection
func (s *Server) Start() {
	var db *database.Database
	var err error

	db, err = database.Connect()

	if err != nil {
		log.Fatal(err) //TODO: panic and recover
	}
	defer db.Close()
	
	if s.Handler == nil {
		s.Handler = router.GetRouter(db)
	}

	if s.UseHTTPS {
		s.startHTTPS()
	} else {
		s.startHTTP()
	}
}

func (s *Server) startHTTP() {
	//TODO: Use hostanme too
	log.Printf("Server started on %d\n", s.HTTPPort)
	log.Fatal(http.ListenAndServe(s.address(), s.Handler))
}

func (s *Server) startHTTPS() {
	//TODO: Correctly setup HTTPS connection
	log.Fatal(http.ListenAndServeTLS(s.address(), s.CertFile, s.KeyFile, s.Handler))
}

func (s *Server) address() string {
	if s.UseHTTPS {
		return ":" + strconv.Itoa(s.HTTPSPort)
	}
	return ":" + strconv.Itoa(s.HTTPPort)
}