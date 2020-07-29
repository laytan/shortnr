package ratelimit

import (
	"net/http"

	"github.com/didip/tollbooth"
)

// Define some defaults
const GeneralRateLimit float64 = 1
const AuthorizationRateLimit float64 = 0.5

// Middleware returns the middleware used to ratelimit perSecond requests per second
func Middleware(perSecond float64) func(next http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(perSecond, nil)
	return func(next http.Handler) http.Handler {
		return tollbooth.LimitFuncHandler(lmt, next.ServeHTTP)
	}
}
