package cf_urls

import "strings"

func Path(template string, params ...string) string {
	s := template

	var zoneId, id string
	if len(params) > 0 {
		zoneId = params[0]
		if len(params) > 1 {
			id = params[1]
		}
	}

	if zoneId != "" {
		s = strings.Replace(s, "{zone_identifier}", zoneId, 1)
	}

	if id != "" {
		s = strings.Replace(s, "{identifier}", id, 1)
	}

	return s
}
