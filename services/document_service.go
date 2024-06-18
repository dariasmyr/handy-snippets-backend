package services

import (
	"pastebin/graph/model"
)

type DocumentService interface {
	CreateDocument(value, title, accessKey string, maxViewCount, ttlMs int) (int, error)
	DeleteDocument(id int, accessKey string) (bool, error)
	GetDocument(id int) (*model.Document, error)
	DeleteExpiredDocuments() error
}
