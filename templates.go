package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateHelper struct {
	TemplateCache  map[string]*template.Template
	PagePath       string
	TmplPath       string
	ComponentsPath string
	UseCache       bool
}

type TemplateDataObj map[string]any

type TemplateData struct {
	Data      map[string]any
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

func NewTemplateHelper(pagePath string, templatePath string, components *string) *TemplateHelper {
	return &TemplateHelper{
		PagePath:       pagePath,
		TmplPath:       templatePath,
		ComponentsPath: *components,
		UseCache:       true,
		TemplateCache:  make(map[string]*template.Template),
	}
}

func (tmplHelper *TemplateHelper) RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	if !tmplHelper.UseCache {
		err := tmplHelper.CreateTemplateCache()
		if err != nil {
			return err
		}
	}

	t, ok := tmplHelper.TemplateCache[tmpl]
	if !ok {
		return fmt.Errorf("template %s does not exist", tmpl)
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

func (tmplHelper *TemplateHelper) CreateTemplateCache() error {
	tmplHelper.TemplateCache = map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", tmplHelper.PagePath))
	if err != nil {
		return err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return err
		}

		if tmplHelper.ComponentsPath != "" {
			components, err := filepath.Glob(fmt.Sprintf("%s/**/*.component.gohtml", tmplHelper.ComponentsPath))
			if err != nil {
				return err
			}

			if len(components) > 0 {
				ts, err = ts.ParseGlob(fmt.Sprintf("%s/**/*.component.gohtml", tmplHelper.ComponentsPath))
				if err != nil {
					return err
				}
			}
		}

		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", tmplHelper.TmplPath))
		if err != nil {
			return err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", tmplHelper.TmplPath))
			if err != nil {
				return err
			}
		}

		tmplHelper.TemplateCache[name] = ts
	}

	return nil
}
