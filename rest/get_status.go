package rest

import (
	"github.com/boggydigital/compton"
	color "github.com/boggydigital/compton/consts/color"
	"github.com/boggydigital/compton/consts/direction"
	"github.com/boggydigital/compton/consts/size"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"time"
)

const (
	StatusSuccess    = "Success"
	StatusProcessing = "Processing"
	StatusError      = "Error"
)

var statusColors = map[string]color.Color{
	StatusSuccess:    color.Green,
	StatusProcessing: color.Yellow,
	StatusError:      color.Red,
}

type Status struct {
	State          string
	TimestampTitle string
	Timestamp      time.Time
	Names          []string
	LastSetIPs     map[string]string
}

func GetStatus(w http.ResponseWriter, r *http.Request) {

	// GET /status

	var err error
	if rdx, err = rdx.RefreshReader(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	syncStarted := getTime(rdx, data.SyncStarted)
	syncErrored := getTime(rdx, data.SyncErrored)
	syncCompleted := getTime(rdx, data.SyncCompleted)

	state := StatusProcessing
	if syncStarted == syncCompleted || syncCompleted.After(syncStarted) {
		state = StatusSuccess
	} else if syncStarted == syncErrored || syncErrored.After(syncStarted) {
		state = StatusError
	}

	syncNames, _ := rdx.GetAllValues(data.SyncResultsProperty, data.SyncNames)
	lastSetIPs := make(map[string]string)
	for _, name := range syncNames {
		lastSetIPs[name], _ = rdx.GetLastVal(data.LastSetIPsProperty, name)
	}

	var tsTitle string
	var tsTime time.Time

	switch state {
	case StatusSuccess:
		tsTime = syncCompleted
		tsTitle = "Last updated:"
	case StatusProcessing:
		tsTime = syncStarted
		tsTitle = "Started:"
	case StatusError:
		tsTime = syncErrored
		tsTitle = "Last error:"
	}

	p := compton.Page("flared")
	p.SetAttribute("style", "--c-rep:var(--c-background)")

	pageStack := compton.FlexItems(p, direction.Column)
	p.Append(pageStack)

	statusHeading := compton.H1()
	statusText := compton.Fspan(p, state).
		ForegroundColor(statusColors[state])
	statusHeading.Append(statusText)

	pageStack.Append(statusHeading)

	domainIpTable := compton.Table(p)
	domainIpTable.AppendHead("Domain", "Address")

	sortedDomains := slices.Sorted(maps.Keys(lastSetIPs))

	for _, domain := range sortedDomains {
		domainIpTable.AppendRow(domain, lastSetIPs[domain])
	}

	pageStack.Append(domainIpTable)

	pageStack.Append(compton.HeadingText("Debug", 2))

	traceLink := compton.A("/trace")
	traceLink.Append(compton.Fspan(p, "Trace").ForegroundColor(color.Blue))
	pageStack.Append(traceLink)

	cfDashLink := compton.A("https://dash.cloudflare.com/")
	cfDashLink.Append(compton.Fspan(p, "Cloudflare dashboard").ForegroundColor(color.Blue))
	pageStack.Append(cfDashLink)

	tsFspan := compton.Fspan(p, tsTitle+" "+tsTime.Format(time.RFC3339)).
		ForegroundColor(color.Gray).
		FontSize(size.XSmall)

	pageStack.Append(tsFspan)

	if err = p.WriteResponse(w); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

func getTime(rdx redux.Readable, p string) time.Time {
	u := int64(0)
	if str, ok := rdx.GetLastVal(data.SyncResultsProperty, p); ok {
		if su, err := strconv.ParseInt(str, 10, 64); err == nil {
			u = su
		}
	}
	return time.Unix(u, 0)
}
