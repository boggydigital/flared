package data

import "github.com/boggydigital/pathways"

const DefaultFlaredRootDir = "/usr/share/flared"

const (
	Input    pathways.AbsDir = "input"
	Metadata pathways.AbsDir = "metadata"
)

var AllAbsDirs = []pathways.AbsDir{
	Input,
	Metadata,
}
