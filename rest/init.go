package rest

import (
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/pathways"
	"github.com/boggydigital/redux"
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
	rdx  redux.Readable
)

func Init(templatesFS fs.FS) error {
	tmpl = template.Must(
		template.New("").
			ParseFS(templatesFS, "templates/*.gohtml"))

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
