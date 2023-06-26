package main

import (
	"bytes"
	_ "embed"
	"github.com/boggydigital/cf_ddns/cli"
	"github.com/boggydigital/cf_ddns/data"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"os"
	"path/filepath"
)

var (
	//once = sync.Once{}
	////go:embed "templates/*.gohtml"
	//templates embed.FS
	////go:embed "stencil_app/styles/css.gohtml"
	//stencilAppStyles embed.FS
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

//var tmplFuncs = template.FuncMap{
//	"empty": empty,
//}

const (
	directoriesFilename = "directories.txt"
	settingsFilename    = "settings.txt"
)

var (
	stateDir = "/var/lib/cf_ddns"
	logsDir  = "/var/log/cf_ddns"
)

func main() {

	nod.EnableStdOutPresenter()

	ns := nod.Begin("novus is checking for any news")
	defer ns.End()

	//once.Do(func() {
	//	if err := rest.Init(templates, tmplFuncs, stencilAppStyles); err != nil {
	//		_ = ns.EndWithError(err)
	//		os.Exit(1)
	//	}
	//	cli.Init(templates, tmplFuncs)
	//})

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

	userDefaultsPath := filepath.Join(stateDir, settingsFilename)
	if _, err := os.Stat(userDefaultsPath); err == nil {
		udoFile, err := os.Open(userDefaultsPath)
		if err != nil {
			_ = ns.EndWithError(err)
			os.Exit(1)
		}
		userDefaultsOverrides, err := wits.ReadKeyValues(udoFile)
		if err != nil {
			_ = ns.EndWithError(err)
			os.Exit(1)
		}
		if err := defs.SetUserDefaults(userDefaultsOverrides); err != nil {
			_ = ns.EndWithError(err)
			os.Exit(1)
		}
	}

	if defs.HasUserDefaultsFlag("debug") {
		logger, err := nod.EnableFileLogger(logsDir)
		if err != nil {
			_ = ns.EndWithError(err)
			os.Exit(1)
		}
		defer logger.Close()
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"create-dns-record": cli.CreateDNSRecordHandler,
		"list-dns-records":  cli.ListDNSRecordsHandler,
		"sync":              cli.SyncHandler,
		"trace":             cli.TraceHandler,
		"update-dns-record": cli.UpdateDNSRecordHandler,
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
