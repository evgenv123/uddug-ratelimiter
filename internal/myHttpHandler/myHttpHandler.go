package myHttpHandler

import (
	"github.com/evgenv123/uddug-ratelimiter/internal/mylimiter"
	"io/ioutil"
	"log"
	"net/http"
)

func MyAPIHandler(w http.ResponseWriter, r *http.Request) {

	if !mylimiter.CheckLimiters(r) {
		log.Println("Rate limit reached!")
		http.Error(w, "Rate limit reached!", http.StatusTooManyRequests)
		return
	}

	// Rate-limiter passed, processing request...
	// Read body
	_, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("<h1>All OK. Request processed!</h1>"))
}
