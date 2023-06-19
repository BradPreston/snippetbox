package main

import "github.com/BradPreston/snippetbox/internal/models"

type templateData struct {
    Snippet *models.Snippet
    Snippets []*models.Snippet
}
