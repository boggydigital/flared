package rest

import "net/http"

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/status": http.HandlerFunc(GetStatus),
		"/":       http.RedirectHandler("/status", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
