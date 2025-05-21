package cli

import (
	"fmt"
	"github.com/boggydigital/flared/cf_api"
	"github.com/boggydigital/flared/cf_trace"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"github.com/boggydigital/wits"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

func SyncHandler(u *url.URL) error {
	q := u.Query()
	token := q.Get("token")
	filename := q.Get("filename")

	return Sync(token, filename)
}

func Sync(token, filename string) error {

	// 1) read DNS requests from filename
	// 2) Fetch current DNS records for each zone
	// 3) trace WAN IP
	// 4) for reach request:
	// 4.1) if a record doesn't exist - create
	// 4.2) if record exists, get id and update

	sa := nod.Begin("syncing DNS records...")
	defer sa.EndWithResult("done")

	// will always set error, unless it's cleared on success at the end
	syncError := true

	amd, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	rdx, err := redux.NewWriter(amd, data.SyncResultsProperty, data.LastSetIPsProperty)
	if err != nil {
		return err
	}

	defer func() {
		if syncError {
			_ = rdx.ReplaceValues(data.SyncResultsProperty, data.SyncErrored, nts())
		}
	}()

	if err = rdx.ReplaceValues(data.SyncResultsProperty, data.SyncStarted, nts()); err != nil {
		return err
	}

	aid, err := pathways.GetAbsDir(data.Input)
	if err != nil {
		return err
	}

	// ignore everything that's not actual filename
	filename = filepath.Base(filename)

	absFilename := filepath.Join(aid, filename)

	rskva := nod.Begin(" reading %s...", absFilename)
	defer rskva.EndWithResult("done")

	f, err := os.Open(absFilename)
	if err != nil {
		return err
	}

	skv, err := wits.ReadSectionKeyValue(f)
	if err != nil {
		return err
	}

	if err = rdx.ReplaceValues(data.SyncResultsProperty, data.SyncNames, slices.Collect(maps.Keys(skv))...); err != nil {
		return err
	}

	rskva.EndWithResult("done")

	zones := make(map[string]interface{})

	for _, kv := range skv {
		if zoneId, ok := kv["zone-id"]; ok {
			zones[zoneId] = nil
		}
	}

	ta := nod.Begin(" tracing WAN IP address...")
	defer ta.EndWithResult("done")

	ipv4 := ""
	tm, err := cf_trace.GetMap(http.DefaultClient)
	if err != nil {
		return err
	}

	if ip, ok := tm["ip"]; ok {
		ipv4 = ip
	}

	ta.EndWithResult("done")

	// check if we can avoid doing any more work:
	// domains
	if alreadySetLatestContent(ipv4, skv, rdx) {

		if err = rdx.ReplaceValues(data.SyncResultsProperty, data.SyncCompleted, nts()); err != nil {
			return err
		}

		sa.EndWithResult("already set latest content")
		return nil
	}

	ldra := nod.NewProgress(" listing current DNS records...")
	defer ldra.EndWithResult("done")

	c := cf_api.NewClient(http.DefaultClient, token)

	zoneRecords := make(map[string][]cf_api.DNSRecordResult)

	ldra.TotalInt(len(zones))

	for zoneId := range zones {
		ldrr, err := c.ListDNSRecords(zoneId)
		if err != nil {
			return err
		}
		if ldrr.Success {
			zoneRecords[zoneId] = ldrr.Result
		} else {
			for _, e := range ldrr.Errors {
				ldra.Error(fmt.Errorf("%d %s", e.Code, e.Message))
			}
		}
		ldra.Increment()
	}

	cua := nod.NewProgress(" setting DNS records...")
	defer cua.EndWithResult("done")

	cua.TotalInt(len(skv))

	lsips := make(map[string][]string)
	for name := range rdx.Keys(data.LastSetIPsProperty) {
		lsips[name], _ = rdx.GetAllValues(data.LastSetIPsProperty, name)
	}

	for name, rkv := range skv {
		zoneId, ok := rkv["zone-id"]
		if !ok {
			nod.Log("%s %s should specify zone-id", filename, name)
			continue
		}

		content := rkv["content"]
		if content == "" {
			content = ipv4
		}

		proxied := rkv["proxied"] == "true"
		recordType := rkv["type"]
		if recordType == "" {
			recordType = "A"
		}
		comment := rkv["comment"]
		var tags []string
		t := rkv["tags"]
		if t != "" {
			tags = strings.Split(t, ",")
		}

		ttl := 1
		ttlStr := rkv["ttl"]
		if ttl64, err := strconv.ParseInt(ttlStr, 10, 32); err == nil {
			ttl = int(ttl64)
		}

		var drr *cf_api.DNSRecordResultResponse

		if id := recordId(zoneRecords[zoneId], name, recordType); id != "" {

			// current content is the same as the last one - no update needed
			if lsip, ok := lsips[name]; ok && len(lsip) > 0 && lsip[0] == content {
				continue
			}

			drr, err = c.UpdateDNSRecord(zoneId, id, content, name, proxied, recordType, comment, tags, ttl)
		} else {
			drr, err = c.CreateDNSRecord(zoneId, content, name, proxied, recordType, comment, tags, ttl)
		}

		lsips[name] = []string{content}

		if err != nil {
			return err
		}
		nodDNSRecordResult(drr)

		cua.Increment()
	}

	if err = rdx.BatchReplaceValues(data.LastSetIPsProperty, lsips); err != nil {
		return err
	}

	if err = rdx.ReplaceValues(data.SyncResultsProperty, data.SyncCompleted, nts()); err != nil {
		return err
	}

	// clear error state
	syncError = false

	return nil
}

func recordId(dnsRecords []cf_api.DNSRecordResult, name, recordType string) string {
	for _, dr := range dnsRecords {
		if dr.Name == name && dr.Type == recordType {
			return dr.Id
		}
	}
	return ""
}

func nts() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func alreadySetLatestContent(ipv4 string, skv wits.SectionKeyValue, rdx redux.Readable) bool {
	for name, kv := range skv {
		if content, ok := kv["content"]; ok {
			if content != "" && content != ipv4 {
				return false
			}
		}
		if lsip, _ := rdx.GetLastVal(data.LastSetIPsProperty, name); lsip != ipv4 {
			return false
		}
	}
	return true
}
