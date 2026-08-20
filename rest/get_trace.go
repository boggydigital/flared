package rest

import (
	"maps"
	"net/http"
	"slices"

	"github.com/boggydigital/flared/cf_trace"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/sizes"
)

func GetTrace(w http.ResponseWriter, r *http.Request) {

	// GET /trace

	traceMap, err := cf_trace.GetMap(http.DefaultClient)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	root, body := strom.RootBody("flared trace", atoms.FlexCol(sizes.Normal)...)
	for head := range root.GetElementsByTagName("head") {
		head.Append(strom.Stylesheet(tableStyles))
		break
	}

	body.Append(strom.CreateText("h1", "Trace Results"))

	table := createTable("Name", "Content")
	body.Append(table)

	keys := slices.Sorted(maps.Keys(traceMap))
	for _, key := range keys {
		appendRow(table, key, traceMap[key])
	}

	if err = strom.WriteResponse(w, root); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
