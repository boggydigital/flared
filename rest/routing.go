package rest

import "net/http"

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/status": http.HandlerFunc(GetStatus),
		"/trace":  http.HandlerFunc(GetTrace),
		"/":       http.RedirectHandler("/status", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
