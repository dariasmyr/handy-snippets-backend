package services

import (
	"pastebin/graph/model"
	"time"
)

type DocumentService interface {
	StartExpiredDocumentsCleaner(interval time.Duration)
	CreateDocument(value string, title *string, accessKey string, maxViewCount, ttlMs int) (int, error)
	UpdateDocument(id int, value *string, title *string, accessKey string, maxViewCount, ttlMs int) (bool, error)
	DeleteDocument(id int, accessKey string) (bool, error)
	GetDocument(id int) (*model.Document, error)
	DeleteExpiredDocuments() error
}
