package data

import (
	"github.com/boggydigital/camino"
)

const flaredRootDir = "/usr/share/flared"

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

	flaredAbsPaths := camino.ResolveAbsPaths(flaredRootDir, absDirNames, nil)

	return camino.Register(flaredAbsPaths, nil, nil)
}
