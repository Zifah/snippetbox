package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.hafiz.com.ng/internal/models"
	"snippetbox.hafiz.com.ng/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	LatestSnippets  []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		name := filepath.Base(p)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			p,
		}

		templateSet, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
