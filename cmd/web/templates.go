package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.hafiz.com.ng/internal/models"
)

type templateData struct {
	CurrentYear    int
	Snippet        *models.Snippet
	LatestSnippets []*models.Snippet
	Form           any
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		name := filepath.Base(p)

		files := []string{
			"./ui/html/base.tmpl",
			p,
		}

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
