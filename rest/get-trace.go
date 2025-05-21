package rest

import (
	"github.com/boggydigital/flared/cf_trace"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetTrace(w http.ResponseWriter, r *http.Request) {

	// GET /trace

	tm, err := cf_trace.GetMap(http.DefaultClient)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}

	DefaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "trace-page", tm); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

}
