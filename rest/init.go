package rest

import (
	"html/template"
	"io/fs"
)

var (
	tmpl *template.Template
)

func Init(templatesFS fs.FS) {
	tmpl = template.Must(
		template.New("").
			ParseFS(templatesFS, "templates/*.gohtml"))
}
