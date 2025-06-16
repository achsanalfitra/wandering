package main

import (
	"net/http"

	"github.com/achsanalfitra/wandering/backend/internal/router"
)

func main() {
	r := router.NewRouter()
	// r.Register("GET", "/hello", http.HandlerFunc(yourfunction)) // register your handler functions

	s := &http.Server{
		Addr:    ":6969", // main app runs on 6969
		Handler: r,
	}

	s.ListenAndServe()
}
