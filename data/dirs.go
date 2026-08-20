package data

import (
	"os"

	"github.com/boggydigital/camino"
)

const (
	flaredRootDir       = "/usr/share/flared"
	directoriesFilename = "directories.txt"
)

const (
	Input camino.AbsDir = iota
	Metadata
	Backups
	Logs
)

var absDirNames = map[camino.AbsDir]string{
	Input:    "input",
	Metadata: "metadata",
	Backups:  "backups",
	Logs:     "logs",
}

func InitFlaredCamino() error {

	var overrides map[string]string

	if _, err := os.Stat(directoriesFilename); err == nil {
		if overrides, err = camino.ReadOverrides(directoriesFilename); err != nil {
			return err
		}
	}

	flaredAbsPaths := camino.ResolveAbsPaths(flaredRootDir, absDirNames, overrides)

	return camino.Register(flaredAbsPaths, nil, nil)
}
