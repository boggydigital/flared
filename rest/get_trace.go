package rest

import (
	"github.com/boggydigital/compton"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/flared/cf_trace"
	"github.com/boggydigital/nod"
	"maps"
	"net/http"
	"slices"
)

func GetTrace(w http.ResponseWriter, r *http.Request) {

	// GET /trace

	traceMap, err := cf_trace.GetMap(http.DefaultClient)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	p := compton.Page("flared trace")
	p.SetAttribute("style", "--c-rep: var(--c-background)")

	pageStack := compton.FlexItems(p, direction.Column)
	p.Append(pageStack)

	pageStack.Append(compton.HeadingText("Trace Results", 1))

	table := compton.Table(p)
	table.AppendHead("Name", "Content")

	keys := slices.Sorted(maps.Keys(traceMap))
	for _, key := range keys {
		table.AppendRow(key, traceMap[key])
	}

	pageStack.Append(table)

	if err = p.WriteResponse(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
