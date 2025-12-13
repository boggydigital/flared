package main

import (
	"bytes"
	_ "embed"
	"log"
	"net/url"
	"os"

	"github.com/boggydigital/clo"
	"github.com/boggydigital/flared/cli"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
)

var (
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

const (
	dirOverridesFilename = "directories.txt"
	debugParam           = "debug"
)

func main() {

	nod.EnableStdOutPresenter()

	ns := nod.Begin("flared is processing DNS records")
	defer ns.EndWithResult("done")

	if err := pathways.Setup(
		dirOverridesFilename,
		data.DefaultFlaredRootDir,
		nil,
		data.AllAbsDirs...); err != nil {
		log.Fatalln(err)
	}

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		log.Fatalln(err)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"backup":            cli.BackupHandler,
		"create-dns-record": cli.CreateDNSRecordHandler,
		"list-dns-records":  cli.ListDNSRecordsHandler,
		"serve":             cli.ServeHandler,
		"sync":              cli.SyncHandler,
		"trace":             cli.TraceHandler,
		"update-dns-record": cli.UpdateDNSRecordHandler,
		"version":           cli.VersionHandler,
	})

	if err = defs.AssertCommandsHaveHandlers(); err != nil {
		log.Fatalln(err)
	}

	var u *url.URL
	u, err = defs.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	if q := u.Query(); q.Has(debugParam) {
		absLogsDir, err := pathways.GetAbsDir(data.Logs)
		if err != nil {
			log.Fatalln(err)
		}
		logger, err := nod.EnableFileLogger(u.Path, absLogsDir)
		if err != nil {
			log.Fatalln(err)
		}
		defer logger.Close()
	}

	if err = defs.Serve(u); err != nil {
		log.Fatalln(err)
	}
}
