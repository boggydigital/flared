package cf_trace

import (
	"github.com/boggydigital/flared/cf_urls"
	"net/url"
)

func TraceUrl() *url.URL {
	return &url.URL{
		Scheme: cf_urls.HTTPS,
		Host:   cf_urls.WwwCloudflareHost,
		Path:   cf_urls.TracePath,
	}
}
