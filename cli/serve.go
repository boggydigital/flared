package cli

import "net/url"

func ServeHandler(u *url.URL) error {
	return Serve(0)
}

func Serve(port int) error {
	return nil
}
