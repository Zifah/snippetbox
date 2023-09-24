package main

import "snippetbox.hafiz.com.ng/internal/models"

type templateData struct {
	Snippet        *models.Snippet
	LatestSnippets []*models.Snippet
}
