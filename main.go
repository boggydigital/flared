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
	"log"
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
	defer ns.EndWithResult("done")

	if err := pathways.Setup(
		dirOverridesFilename,
		data.DefaultFlaredRootDir,
		nil,
		data.AllAbsDirs...); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	once.Do(func() {
		if err := rest.Init(templates); err != nil {
			log.Println(err.Error())
			os.Exit(2)
		}
	})

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		log.Println(err.Error())
		os.Exit(3)
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
		log.Println(err.Error())
		os.Exit(4)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		log.Println(err.Error())
		os.Exit(5)
	}
}
