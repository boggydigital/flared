package cf_trace

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
)

func GetMap(hc *http.Client) (map[string]string, error) {

	resp, err := hc.Get(TraceUrl().String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(resp.Status)
	}

	tm := make(map[string]string)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if k, v, ok := strings.Cut(scanner.Text(), "="); ok {
			tm[k] = v
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tm, nil
}
