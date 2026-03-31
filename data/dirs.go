package data

import "github.com/boggydigital/pathways"

const defaultRootDir = "/usr/share/flared"

const (
	Input    pathways.AbsDir = "input"
	Metadata pathways.AbsDir = "metadata"
	Backups  pathways.AbsDir = "backups"
	Logs     pathways.AbsDir = "logs"
)

var AllAbsDirs = []pathways.AbsDir{
	Input,
	Metadata,
	Backups,
	Logs,
}

var Pwd pathways.Pathway

func InitPathways() error {
	var err error
	Pwd, err = pathways.NewRoot(defaultRootDir)
	return err
}
