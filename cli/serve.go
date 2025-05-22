package cli

import (
	"fmt"
	"github.com/boggydigital/flared/rest"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
)

func ServeHandler(u *url.URL) error {
	portStr := u.Query().Get("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	return Serve(
		port,
		u.Query().Has("stderr"))
}

func Serve(port int, stderr bool) error {

	if stderr {
		nod.EnableStdErrLogger()
		nod.DisableOutput(nod.StdOut)
	}

	sa := nod.Begin("serving at port %d...", port)
	defer sa.EndWithResult("done")

	if err := rest.Init(); err != nil {
		return err
	}

	rest.HandleFuncs()

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
