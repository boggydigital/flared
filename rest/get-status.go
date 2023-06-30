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
}

func GetStatus(w http.ResponseWriter, r *http.Request) {

	// GET /status

	rdx, err := kvas.ConnectRedux(data.Pwd(), data.SyncResultsProperty)
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

	syncNames, _ := rdx.GetAllValues(data.SyncNames)
	sort.Strings(syncNames)

	status := &Status{
		State: state,
		Names: syncNames,
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

func getTime(rdx kvas.ReduxValues, p string) time.Time {
	u := int64(0)
	if str, ok := rdx.GetFirstVal(p); ok {
		if su, err := strconv.ParseInt(str, 10, 64); err == nil {
			u = su
		}
	}
	return time.Unix(u, 0)
}
