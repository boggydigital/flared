package data

import "github.com/boggydigital/pathways"

const DefaultFlaredRootDir = "/usr/share/flared"

const (
	Input    pathways.AbsDir = "input"
	Metadata pathways.AbsDir = "metadata"
	Backups  pathways.AbsDir = "backups"
)

var AllAbsDirs = []pathways.AbsDir{
	Input,
	Metadata,
	Backups,
}
