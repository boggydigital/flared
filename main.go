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
	"github.com/boggydigital/wits"
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
	directoriesFilename = "directories.txt"
)

var (
	stateDir = "/var/lib/flared"
	logsDir  = "/var/log/flared"
)

func main() {

	nod.EnableStdOutPresenter()

	ns := nod.Begin("flared is processing DNS records")
	defer ns.End()

	once.Do(func() {
		rest.Init(templates)
	})

	if err := readUserDirectories(); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	data.ChRoot(stateDir)

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"create-dns-record": cli.CreateDNSRecordHandler,
		"list-dns-records":  cli.ListDNSRecordsHandler,
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

func readUserDirectories() error {
	if _, err := os.Stat(directoriesFilename); os.IsNotExist(err) {
		return nil
	}

	udFile, err := os.Open(directoriesFilename)
	if err != nil {
		return err
	}

	dirs, err := wits.ReadKeyValue(udFile)
	if err != nil {
		return err
	}

	if sd, ok := dirs["state"]; ok {
		stateDir = sd
	}
	if ld, ok := dirs["logs"]; ok {
		logsDir = ld
	}
	//validate that directories actually exist
	if _, err := os.Stat(stateDir); err != nil {
		return err
	}
	if _, err := os.Stat(logsDir); err != nil {
		return err
	}

	return nil
}
