package cli

import (
	"fmt"
	"github.com/boggydigital/cf_api"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
)

func ListDNSRecordsHandler(u *url.URL) error {

	q := u.Query()
	token := q.Get("token")
	zoneId := q.Get("zone-id")

	return ListDNSRecords(token, zoneId)
}

func ListDNSRecords(token, zoneId string) error {

	ldra := nod.Begin("listing DNS records for zone:%s...", zoneId)
	defer ldra.End()

	client := cf_api.NewClient(http.DefaultClient, token)

	ldrr, err := client.ListDNSRecords(zoneId)
	if err != nil {
		return err
	}

	pdra := nod.Begin("")
	defer pdra.End()

	summary := make(map[string][]string)
	heading := ""

	if ldrr.Success {
		ldra.EndWithResult("success")

		heading = "DNS records:"
		for _, drr := range ldrr.Result {
			for k, v := range drr.Summary() {
				summary[k] = v
			}
		}

	} else {
		ldra.EndWithResult("error")

		for i, e := range ldrr.Errors {
			summary["error "+strconv.Itoa(i)] = []string{fmt.Sprintf("%d:%s", e.Code, e.Message)}
		}
		heading = "Errors:"
	}

	pdra.EndWithSummary(heading, summary)

	return nil
}
