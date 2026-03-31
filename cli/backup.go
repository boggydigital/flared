package cli

import (
	"net/url"

	"github.com/boggydigital/backups"
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

	amp := data.Pwd.AbsDirPath(data.Metadata)
	abp := data.Pwd.AbsDirPath(data.Backups)

	if err := backups.Compress(amp, abp); err != nil {
		return err
	}

	cba := nod.NewProgress("cleaning up old backups...")
	defer cba.EndWithResult("done")

	if err := backups.Cleanup(abp, true, cba); err != nil {
		return err
	}

	return nil
}
