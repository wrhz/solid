package middleware

import "net/http"

func httpError(w http.ResponseWriter) {
	http.Error(w, "{ \"error\": \"Internal Server Error\" }", http.StatusInternalServerError)
}