package cli

import (
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"net/url"
)

func MigrateHandler(_ *url.URL) error {
	return Migrate()
}

func Migrate() error {
	ma := nod.Begin("migrating data...")
	defer ma.End()

	if err := Backup(); err != nil {
		return ma.EndWithError(err)
	}

	amd, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return ma.EndWithError(err)
	}

	if err := kevlar.Migrate(amd); err != nil {
		return ma.EndWithError(err)
	}

	ma.EndWithResult("done")

	return nil
}
