package cf_api

import (
	"fmt"
	"time"
)

type CodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DNSRecordResult struct {
	Content   string    `json:"content"`
	Name      string    `json:"name"`
	Proxied   bool      `json:"proxied"`
	Type      string    `json:"type"`
	Comment   string    `json:"comment"`
	CreatedOn time.Time `json:"created_on"`
	Id        string    `json:"id"`
	Locked    bool      `json:"locked"`
	Meta      struct {
		AutoAdded bool   `json:"auto_added"`
		Source    string `json:"source"`
	} `json:"meta"`
	ModifiedOn time.Time `json:"modified_on"`
	Proxiable  bool      `json:"proxiable"`
	Tags       []string  `json:"tags"`
	TTL        int       `json:"ttl"`
	ZoneId     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
}

func (drr *DNSRecordResult) Summary() map[string][]string {
	result := make(map[string][]string)

	result[drr.Id] = []string{
		fmt.Sprintf("%s=%s", "type", drr.Type),
		fmt.Sprintf("%s=%s", "name", drr.Name),
		fmt.Sprintf("%s=%s", "content", drr.Content),
		fmt.Sprintf("%s=%t", "proxied", drr.Proxied),
	}

	return result
}

type ResultInfo struct {
	Count      int `json:"count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
}

type DNSRecordsResultsResponse struct {
	Errors     []CodeMessage     `json:"errors"`
	Messages   []CodeMessage     `json:"messages"`
	Result     []DNSRecordResult `json:"result"`
	Success    bool              `json:"success"`
	ResultInfo ResultInfo        `json:"result_info"`
}

type DNSRecordResultResponse struct {
	Errors     []CodeMessage   `json:"errors"`
	Messages   []CodeMessage   `json:"messages"`
	Result     DNSRecordResult `json:"result"`
	Success    bool            `json:"success"`
	ResultInfo ResultInfo      `json:"result_info"`
}

type DNSRequest struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	Type    string   `json:"type"`
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	TTL     int      `json:"ttl"`
}
