package rest

import (
	"github.com/boggydigital/middleware"
	"net/http"
)

var (
	GetOnly = middleware.GetMethodOnly
	BrGzip  = middleware.BrGzip
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/status": BrGzip(GetOnly(http.HandlerFunc(GetStatus))),
		"/trace":  BrGzip(GetOnly(http.HandlerFunc(GetTrace))),
		"/":       http.RedirectHandler("/status", http.StatusPermanentRedirect),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
