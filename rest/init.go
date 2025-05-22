package rest

import (
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
)

var (
	rdx redux.Readable
)

func Init() error {
	amd, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	rdx, err = redux.NewReader(amd, data.AllProperties()...)
	if err != nil {
		return err
	}

	return nil
}
