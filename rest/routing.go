package rest

import (
	"net/http"
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"GET /status": http.HandlerFunc(GetStatus),
		"GET /trace":  http.HandlerFunc(GetTrace),
		"GET /":       http.RedirectHandler("/status", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
