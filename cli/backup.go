package cli

import (
	"net/url"

	"github.com/boggydigital/camino"
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/nod"
)

const daysToPreserveFiles = 30

func BackupHandler(_ *url.URL) error {
	return Backup()
}

func Backup() error {
	ea := nod.NewProgress("backing up metadata...")
	defer ea.EndWithResult("done")

	if err := camino.Compress(data.Metadata, data.Backups); err != nil {
		return err
	}

	cba := nod.Begin("cleaning up old backups...")
	defer cba.EndWithResult("done")

	if err := camino.CleanupTimed(data.Backups, true); err != nil {
		return err
	}

	return nil
}
