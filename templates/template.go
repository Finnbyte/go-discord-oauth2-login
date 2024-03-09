package templates

import (
	"html/template"
	"strings"
	"fmt"
	"os"
	"path/filepath"
)

func processTemplates() *template.Template {
	t := template.New("")
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			fmt.Println(path)
			_, err = t.ParseFiles(path)
			if err != nil {
				fmt.Println(err)
			}
		}
		return err
	})

    if err != nil {
        panic(err)
    }
    return t
}

var templates = processTemplates()

func Get() *template.Template {
	return templates
}

