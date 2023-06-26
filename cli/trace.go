package cli

import (
	"github.com/boggydigital/cf_api/cf_trace"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
)

func TraceHandler(u *url.URL) error {
	return Trace()
}

func Trace() error {

	ta := nod.Begin("tracing WAN IP address...")
	defer ta.End()

	tm, err := cf_trace.GetMap(http.DefaultClient)
	if err != nil {
		return ta.EndWithError(err)
	}

	summary := make(map[string][]string)
	for k, v := range tm {
		summary[k] = []string{v}
	}

	ta.EndWithSummary("trace results:", summary)

	return nil
}
