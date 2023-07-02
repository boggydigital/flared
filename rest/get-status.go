package rest

import (
	"github.com/boggydigital/cf_ddns/data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	StatusSuccess    = "success"
	StatusProcessing = "processing"
	StatusError      = "error"
)

type Status struct {
	State          string
	TimestampTitle string
	Timestamp      time.Time
	Names          []string
	LastSetIPs     map[string]string
}

func GetStatus(w http.ResponseWriter, r *http.Request) {

	// GET /status

	rdx, err := kvas.ConnectReduxAssets(data.Pwd(), nil, data.SyncResultsProperty, data.LastSetIPsProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	syncStarted := getTime(rdx, data.SyncStarted)
	syncErrored := getTime(rdx, data.SyncErrored)
	syncCompleted := getTime(rdx, data.SyncCompleted)

	state := StatusProcessing
	if syncCompleted.After(syncStarted) {
		state = StatusSuccess
	} else if syncErrored.After(syncStarted) {
		state = StatusError
	}

	syncNames, _ := rdx.GetAllValues(data.SyncResultsProperty, data.SyncNames)
	sort.Strings(syncNames)

	lastSetIPs := make(map[string]string)
	for _, name := range syncNames {
		lastSetIPs[name], _ = rdx.GetFirstVal(data.LastSetIPsProperty, name)
	}

	status := &Status{
		State:      state,
		Names:      syncNames,
		LastSetIPs: lastSetIPs,
	}

	switch status.State {
	case StatusSuccess:
		status.Timestamp = syncCompleted
		status.TimestampTitle = "Last updated:"
	case StatusProcessing:
		status.Timestamp = syncStarted
		status.TimestampTitle = "Started:"
	case StatusError:
		status.Timestamp = syncErrored
		status.TimestampTitle = "Last error:"
	}

	DefaultHeaders(w)

	if err := tmpl.ExecuteTemplate(w, "status-page", status); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}

func getTime(rdx kvas.ReduxAssets, p string) time.Time {
	u := int64(0)
	if str, ok := rdx.GetFirstVal(data.SyncResultsProperty, p); ok {
		if su, err := strconv.ParseInt(str, 10, 64); err == nil {
			u = su
		}
	}
	return time.Unix(u, 0)
}
