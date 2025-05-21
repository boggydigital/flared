package cli

import (
	"github.com/boggydigital/flared/cf_trace"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
)

func TraceHandler(u *url.URL) error {
	return Trace()
}

func Trace() error {

	ta := nod.Begin("tracing...")
	defer ta.EndWithResult("done")

	tm, err := cf_trace.GetMap(http.DefaultClient)
	if err != nil {
		return err
	}

	summary := make(map[string][]string)
	for k, v := range tm {
		val := k + "=" + v
		summary[val] = nil
	}

	ta.EndWithSummary("", summary)

	return nil
}
