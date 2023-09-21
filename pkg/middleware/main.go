package middleware

import (
	"log"
	"net/http"
)

type Chain func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, c []Chain) http.Handler{
	var res http.Handler
	if len(c) == 0 {
		res = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})

		return res
	}
	
	for i := len(c); i >= 0; i-- {
		res = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}

	return res
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start log")
		h.ServeHTTP(w, r)
		log.Println("end log")
	})
}

func HelloMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start hello")
		h.ServeHTTP(w, r)
		log.Println("end log")
	})
}
