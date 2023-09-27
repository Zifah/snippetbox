package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.hafiz.com.ng/internal/models"
)

type templateData struct {
	Snippet        *models.Snippet
	LatestSnippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		files := []string{
			"./ui/html/base.tmpl",
			p,
		}

		templateSet, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		name := filepath.Base(p)
		cache[name] = templateSet
	}

	return cache, nil
}
