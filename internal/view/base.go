package view

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type TemplateEngine struct {
	templates *template.Template
}

func NewTemplateEngine() (TemplateEngine, error) {
	return TemplateEngine{}, nil
}

func (t *TemplateEngine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *TemplateEngine) Load() error {
	templ := template.New("")
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		return err
	}
	t.templates = templ
	return nil
}
