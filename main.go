package main

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/flared/cli"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/flared/rest"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"os"
	"sync"
)

var (
	once = sync.Once{}
	//go:embed "templates/*.gohtml"
	templates embed.FS
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

const (
	dirOverridesFilename = "directories.txt"
)

var (
	stateDir = "/var/lib/flared"
	logsDir  = "/var/log/flared"
)

func main() {

	nod.EnableStdOutPresenter()

	ns := nod.Begin("flared is processing DNS records")
	defer ns.End()

	if err := pathways.Setup(
		dirOverridesFilename,
		data.DefaultFlaredRootDir,
		nil,
		data.AllAbsDirs...); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	once.Do(func() {
		if err := rest.Init(templates); err != nil {
			_ = ns.EndWithError(err)
			os.Exit(1)
		}
	})

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
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

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}
}
