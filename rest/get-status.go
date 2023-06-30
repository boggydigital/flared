package rest

import (
	"github.com/boggydigital/cf_ddns/data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/http"
	"strconv"
	"time"
)

type Status struct {
	Success   bool
	Completed time.Time
	Errored   time.Time
	Names     []string
}

func GetStatus(w http.ResponseWriter, r *http.Request) {

	// GET /status

	rdx, err := kvas.ConnectRedux(data.Pwd(), data.SyncResultsProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	syncNames, _ := rdx.GetAllValues(data.SyncNames)
	syncStarted := getTime(rdx, data.SyncStarted)
	syncErrored := getTime(rdx, data.SyncErrored)
	syncCompleted := getTime(rdx, data.SyncCompleted)

	success := false
	if syncCompleted.After(syncStarted) {
		success = true
	}

	status := &Status{
		Success:   success,
		Completed: syncCompleted,
		Errored:   syncErrored,
		Names:     syncNames,
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
