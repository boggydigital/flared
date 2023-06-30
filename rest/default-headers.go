package rest

import "net/http"

const (
	htmlContentType = "text/html"
	defaultCSP      = "default-src 'self'; " +
		"script-src 'none';" +
		"object-src 'none'; " +
		"img-src 'self' data:; " +
		"style-src 'unsafe-inline';"
)

func DefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", htmlContentType)
	stencilCSP := defaultCSP
	w.Header().Set("Content-Security-Policy", stencilCSP)
}
