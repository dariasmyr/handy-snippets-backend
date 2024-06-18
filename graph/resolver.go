package graph

import (
	"pastebin/services"
)

type Resolver struct {
	DocumentService services.DocumentService
}

func NewResolver(ds services.DocumentService) *Resolver {
	return &Resolver{
		DocumentService: ds,
	}
}
