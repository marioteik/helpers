package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var tc map[string]*template.Template
var pgPath = "./cmd/web/components/pages"
var tmPath = "./cmd/web/components/templates"
var UseCache = true

type TemplateDataObj map[string]any

type TemplateData struct {
	Data      map[string]any
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *TemplateData) error {
	if !UseCache {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return fmt.Errorf("Template %s does not exist", tmpl)
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, data)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func CreateTemplateCache(pagePath *string, tmplPath *string) (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	if pagePath != nil {
		pgPath = *pagePath
	}

	if tmplPath != nil {
		tmPath = *tmplPath
	}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pgPath))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", tmPath))
		if err != nil {
			return nil, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", tmPath))
			if err != nil {
				return nil, err
			}
		}

		tmplCache[name] = ts
	}

	tc = tmplCache

	return tmplCache, nil
}
