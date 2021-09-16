// Simple HTTP server with rate-limiter and one API endpoint
package main

import (
	"github.com/evgenv123/uddug-ratelimiter/internal/myHttpHandler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api", myHttpHandler.MyAPIHandler)
	log.Println("Starting web server at *:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

