package cli

import (
	"github.com/boggydigital/flared/cf_api"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func UpdateDNSRecordHandler(u *url.URL) error {
	q := u.Query()
	token := q.Get("token")
	zoneId := q.Get("zone-id")
	id := q.Get("id")
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

	return UpdateDNSRecord(token, zoneId, id, content, name, proxied, recordType, comment, tags, ttl)
}

func UpdateDNSRecord(
	token string,
	zoneId, id string,
	content, name string,
	proxied bool,
	recordType, comment string,
	tags []string,
	ttl int) error {

	udra := nod.Begin("updating DNS record...")
	defer udra.EndWithResult("done")

	c := cf_api.NewClient(http.DefaultClient, token)

	cdrr, err := c.UpdateDNSRecord(zoneId, id, content, name, proxied, recordType, comment, tags, ttl)
	if err != nil {
		return err
	}

	udra.EndWithResult(nodDNSRecordResult(cdrr))

	return nil
}
