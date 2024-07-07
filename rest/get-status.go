package rest

import (
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/kevlar"
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
	sort.Strings(syncNames)

	lastSetIPs := make(map[string]string)
	for _, name := range syncNames {
		lastSetIPs[name], _ = rdx.GetLastVal(data.LastSetIPsProperty, name)
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

func getTime(rdx kevlar.ReadableRedux, p string) time.Time {
	u := int64(0)
	if str, ok := rdx.GetLastVal(data.SyncResultsProperty, p); ok {
		if su, err := strconv.ParseInt(str, 10, 64); err == nil {
			u = su
		}
	}
	return time.Unix(u, 0)
}
