package cf_api

import (
	"encoding/json"
	"errors"
	"github.com/boggydigital/flared/cf_urls"
	"net/http"
	"net/url"
)

func updateDNSRecordUrl(zoneId, id string) *url.URL {
	return &url.URL{
		Scheme: cf_urls.HTTPS,
		Host:   cf_urls.ApiCloudflareHost,
		Path:   cf_urls.Path(cf_urls.UpdateDNSRecordPathTemplate, zoneId, id),
	}
}

// UpdateDNSRecord
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-update-dns-record
func (c *Client) UpdateDNSRecord(
	zoneId, id string,
	content, name string,
	proxied bool,
	recordType, comment string,
	tags []string,
	ttl int) (*DNSRecordResultResponse, error) {

	dr := &DNSRequest{
		Content: content,
		Name:    name,
		Proxied: proxied,
		Type:    recordType,
		Comment: comment,
		Tags:    tags,
		TTL:     ttl,
	}

	bts, err := json.Marshal(dr)
	if err != nil {
		return nil, err
	}

	udru := updateDNSRecordUrl(zoneId, id)
	req, err := c.newAuthRequest(http.MethodPut, udru, bts...)
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

	var cdr *DNSRecordResultResponse

	err = json.NewDecoder(resp.Body).Decode(&cdr)
	return cdr, err
}

func (c *Client) UpdateDNSARecord(
	zoneId, id string,
	content, name string, proxied bool) (*DNSRecordResultResponse, error) {
	return c.UpdateDNSRecord(
		zoneId, id,
		content, name,
		proxied,
		"A", "",
		nil,
		1) // Setting to 1 means 'automatic'
}
