package cf_api

import (
	"encoding/json"
	"errors"
	"github.com/boggydigital/flared/cf_urls"
	"net/http"
	"net/url"
)

func listDNSRecordsUrl(zoneId string) *url.URL {
	return &url.URL{
		Scheme: cf_urls.HTTPS,
		Host:   cf_urls.ApiCloudflareHost,
		Path:   cf_urls.Path(cf_urls.ListDNSRecordsPathTemplate, zoneId),
	}
}

// ListDNSRecords
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records
func (c *Client) ListDNSRecords(zoneId string) (*DNSRecordsResultsResponse, error) {

	ldru := listDNSRecordsUrl(zoneId)
	req, err := c.newAuthRequest(http.MethodGet, ldru)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(resp.Status)
	}

	var ldr *DNSRecordsResultsResponse

	err = json.NewDecoder(resp.Body).Decode(&ldr)
	return ldr, err
}
