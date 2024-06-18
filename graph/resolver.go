package graph

import (
	"database/sql"
	"pastebin/services"
)

type Resolver struct {
	DocumentService services.DocumentService
}

func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{
		DocumentService: services.NewDocumentService(db),
	}
}
