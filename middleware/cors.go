package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	solidManager "github.com/wrhz/solid/manager"
)

func corsMiddleware(w http.ResponseWriter, r *http.Request) bool {
	corsConfig := solidManager.GetCorsConfig()

	allowOrigin := corsConfig.GetAllowOrigin()

	if len(allowOrigin) == 0 {
		return false
	}

	allowCredentials := corsConfig.GetAllowCredentials()

	if allowOrigin[0] == "*" && allowCredentials {
		fmt.Println("The Allow-Credentials should be false when Allow-Origin is \"*\"")
		httpError(w)
		return false
	}

	origin := r.Header.Get("Origin")

	if allowOrigin[0] != "*" {
		if origin == "" {
			return true
		}

		if !slices.Contains(allowOrigin, origin) {
			w.WriteHeader(204)
			return false
		}

		w.Header().Add("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}

	if r.Method == "OPTIONS" {
		allowMethods := corsConfig.GetAllowMethods()
		requestMethod := r.Header.Get("Access-Control-Request-Method")

		if requestMethod != "" && containsIgnoreCase(allowMethods, requestMethod) {
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
		}

		allowHeaders := corsConfig.GetAllowHeaders()
		requestHeaders := r.Header.Get("Access-Control-Request-Headers")

		if requestHeaders != "" && slices.Contains(allowHeaders, requestHeaders) {
			w.Header().Add("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "))
		}

		w.Header().Add("Access-Control-Max-Age", strconv.FormatUint(uint64(corsConfig.GetMaxAge()), 10))

		w.WriteHeader(204)
		return false
	}

	w.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(allowCredentials))

	w.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.GetExposeHeaders(), ", "))

	return true
}