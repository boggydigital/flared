package rest

import (
	"github.com/boggydigital/camino"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/redux"
)

var (
	rdx redux.Readable
)

func Init() error {

	amd := camino.GetAbs(data.Metadata)

	var err error
	rdx, err = redux.NewReader(amd, data.AllProperties()...)
	if err != nil {
		return err
	}

	return nil
}
