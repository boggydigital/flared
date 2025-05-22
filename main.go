package main

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/flared/cli"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"log"
	"os"
)

var (
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

const (
	dirOverridesFilename = "directories.txt"
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
		"migrate":           cli.MigrateHandler,
		"serve":             cli.ServeHandler,
		"sync":              cli.SyncHandler,
		"trace":             cli.TraceHandler,
		"update-dns-record": cli.UpdateDNSRecordHandler,
		"version":           cli.VersionHandler,
	})

	if err = defs.AssertCommandsHaveHandlers(); err != nil {
		log.Fatalln(err)
	}

	if err = defs.Serve(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}
