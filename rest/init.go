package rest

import (
	"github.com/boggydigital/flared/data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/pathways"
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
	rdx  kevlar.ReadableRedux
)

func Init(templatesFS fs.FS) error {
	tmpl = template.Must(
		template.New("").
			ParseFS(templatesFS, "templates/*.gohtml"))

	amd, err := pathways.GetAbsDir(data.Metadata)
	if err != nil {
		return err
	}

	rdx, err = kevlar.NewReduxReader(amd, data.AllProperties()...)
	if err != nil {
		return err
	}

	return nil
}
