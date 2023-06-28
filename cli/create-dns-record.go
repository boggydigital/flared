package cli

import (
	"fmt"
	"github.com/boggydigital/cf_api"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CreateDNSRecordHandler(u *url.URL) error {

	q := u.Query()
	token := q.Get("token")
	zoneId := q.Get("zone-id")
	content := q.Get("content")
	name := q.Get("name")
	proxied := q.Get("proxied") == "true"
	recordType := q.Get("record-type")
	comment := q.Get("comment")

	tags := make([]string, 0)
	if q.Has("tags") {
		tags = strings.Split(q.Get("tags"), ",")
	}

	ttl := 1
	ttlStr := q.Get("ttl")
	if ttl64, err := strconv.ParseInt(ttlStr, 10, 32); err == nil {
		ttl = int(ttl64)
	}

	return CreateDNSRecord(token, zoneId, content, name, proxied, recordType, comment, tags, ttl)
}

func CreateDNSRecord(
	token string,
	zoneId string,
	content, name string,
	proxied bool,
	recordType, comment string,
	tags []string,
	ttl int) error {

	cdra := nod.Begin("creating DNS record...")
	defer cdra.End()

	c := cf_api.NewClient(http.DefaultClient, token)

	cdrr, err := c.CreateDNSRecord(zoneId, content, name, proxied, recordType, comment, tags, ttl)
	if err != nil {
		return cdra.EndWithError(err)
	}

	cdra.EndWithResult(nodDNSRecordResult(cdrr))

	return nil
}

func nodDNSRecordResult(drrr *cf_api.DNSRecordResultResponse) string {
	pdra := nod.Begin("")
	defer pdra.End()

	summary := make(map[string][]string)
	heading := ""

	result := "success"

	if drrr.Success {

		heading = "DNS record:"
		summary = drrr.Result.Summary()

	} else {
		result = "error"

		for i, e := range drrr.Errors {
			summary["error "+strconv.Itoa(i)] = []string{fmt.Sprintf("%d:%s", e.Code, e.Message)}
		}
		heading = "Errors:"
	}

	pdra.EndWithSummary(heading, summary)

	return result
}
