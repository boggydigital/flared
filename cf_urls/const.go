package cf_urls

const (
	HTTPS = "https"
)

// hosts
const (
	cloudflareHost    = "cloudflare.com"
	WwwCloudflareHost = "www." + cloudflareHost
	ApiCloudflareHost = "api." + cloudflareHost
)

// paths
const (
	TracePath                   = "/cdn-cgi/trace"
	apiBasePath                 = "/client/v4"
	dnsRecordsPathTemplate      = apiBasePath + "/zones/{zone_identifier}/dns_records"
	ListDNSRecordsPathTemplate  = dnsRecordsPathTemplate
	CreateDNSRecordPathTemplate = dnsRecordsPathTemplate
	UpdateDNSRecordPathTemplate = dnsRecordsPathTemplate + "/{identifier}"
)
