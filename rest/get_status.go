package rest

import (
	_ "embed"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars/atoms"
	"github.com/boggydigital/strom/vars/colors"
	"github.com/boggydigital/strom/vars/font_sizes"
	"github.com/boggydigital/strom/vars/font_weights"
	"github.com/boggydigital/strom/vars/sizes"
)

const (
	StatusSuccess    = "Success"
	StatusProcessing = "Processing"
	StatusError      = "Error"
)

var statusColors = map[string]colors.Color{
	StatusSuccess:    colors.Green,
	StatusProcessing: colors.Yellow,
	StatusError:      colors.Red,
}

//go:embed "styles/table.css"
var tableStyles []byte

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

	tsText := tsTitle + " " + tsTime.Format(time.DateTime)

	root, body := strom.RootBody("flared", atoms.FlexCol(sizes.Normal)...)

	for head := range root.GetElementsByTagName("head") {
		head.Append(strom.Stylesheet(tableStyles))
		break
	}

	body.Append(strom.CreateText("h1", state).SetStyle("color:" + statusColors[state]))

	dipTable := createTable("Domain", "Address")

	sortedDomains := slices.Sorted(maps.Keys(lastSetIPs))
	for _, domain := range sortedDomains {
		appendRow(dipTable, domain, lastSetIPs[domain])
	}

	body.Append(dipTable)

	body.Append(strom.CreateText("h2", "Debug"))

	body.Append(
		strom.CreateText("a", "Trace").
			SetAttribute("href", "/trace").
			SetStyle("color:"+colors.Blue, "font-weight:"+font_weights.Bold))

	body.Append(
		strom.CreateText("a", "Cloudflare dashboard").
			SetAttribute("href", "https://dash.cloudflare.com/").
			SetStyle("color:"+colors.Blue, "font-weight:"+font_weights.Bold))

	body.Append(strom.CreateText("span", tsText).
		SetStyle("color:"+colors.Gray, "font-size:"+font_sizes.XSmall))

	if err = strom.WriteResponse(w, root); err != nil {
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

func createTable(headings ...string) strom.Element {

	table := strom.Create("table")
	thead := strom.Create("thead")
	tableRow := strom.Create("tr")
	table.Append(thead)
	thead.Append(tableRow)

	for _, heading := range headings {
		tableRow.Append(strom.CreateText("th", heading))
	}

	tbody := strom.Create("tbody")
	table.Append(tbody)

	return table
}

func appendRow(table strom.Element, values ...string) {

	var tbody strom.Element
	for tb := range table.GetElementsByTagName("tbody") {
		tbody = tb
		break
	}

	if tbody == nil {
		return
	}

	row := strom.Create("tr")
	for _, value := range values {
		row.Append(
			strom.CreateText("td", value))
	}
	tbody.Append(row)
}
