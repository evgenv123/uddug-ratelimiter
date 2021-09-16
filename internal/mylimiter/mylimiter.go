// Package mylimiter implements token bucket algorithm for rate limiting
// https://en.wikipedia.org/wiki/Token_bucket
// could use "golang.org/x/time/rate"
package mylimiter

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type limiter struct {
	mu             sync.Mutex
	Period         time.Duration
	MaxEvents      int
	TokensInBucket int
}

// We will have 3 types of limiters: based on ip, user (auth token) and http endpoint (uri)
var (
	ipLimiters   = make(map[string]*limiter)
	userLimiters = make(map[string]*limiter)
	uriLimiters  = make(map[string]*limiter)
)

// limiterTicker refills token bucket
// could use time subtraction without goroutines
func limiterTicker(l *limiter) {
	l.mu.Lock()
	tick := l.Period
	l.mu.Unlock()
	for {
		time.Sleep(tick)
		l.mu.Lock()
		l.TokensInBucket = l.MaxEvents
		l.mu.Unlock()
		// TODO: Stop goroutine to save memory after some time of inactivity
	}
}

// limiterEvent checks if 1 event can occur.
// Returns true if can occur (and decreases number of tokens in the bucket),
// returns false if rate limit is exceeded (doesn't touch number of tokens in this case)
func limiterEvent(l *limiter) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.TokensInBucket >= 1 {
		l.TokensInBucket -= 1
		return true
	}
	return false
}

// createLimiter Creates new Limiter (token bucket), initiating it full of tokens.
// Also, we start goroutine to refill our bucket
func createLimiter(period time.Duration, maxEvents int) *limiter {
	newLimiter := limiter{
		Period:         period,
		MaxEvents:      maxEvents,
		TokensInBucket: maxEvents,
	}
	go limiterTicker(&newLimiter)
	return &newLimiter
}

// CheckLimiters initializes 3 types of Limiters and checks if events can occur
func CheckLimiters(r *http.Request) bool {
	// Using net.ParseIP to get rid of port (:xxxx) in RemoteAddr
	clientIPAddr := string(net.ParseIP(r.RemoteAddr))
	userAuthToken := r.Header.Get("Authorization")
	uri := r.RequestURI

	// Create limiters if we don't have it for our request
	if _, ok := ipLimiters[clientIPAddr]; !ok {
		ipLimiters[clientIPAddr] = createLimiter(1*time.Minute, 11)
	}
	if _, ok := userLimiters[userAuthToken]; !ok {
		userLimiters[userAuthToken] = createLimiter(1*time.Second, 2)
	}
	if _, ok := uriLimiters[uri]; !ok {
		uriLimiters[uri] = createLimiter(1*time.Second, 10)
	}

	// If any of limiters can't make event we return "false"
	return limiterEvent(ipLimiters[clientIPAddr]) && limiterEvent(userLimiters[userAuthToken]) && limiterEvent(uriLimiters[uri])
}

