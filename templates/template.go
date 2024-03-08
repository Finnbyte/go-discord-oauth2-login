package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderByFilename(w http.ResponseWriter, filename string, data any) error {
	absolutePath, err := filepath.Abs("./templates/" + filename)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(absolutePath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
