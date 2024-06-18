package services

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"pastebin/graph/model"
	"time"
)

type documentService struct {
	db *sql.DB
}

func NewDocumentService(db *sql.DB) DocumentService {
	return &documentService{db: db}
}

func (s *documentService) CreateDocument(value, title, accessKey string, maxViewCount, ttlMs int) (int, error) {
	now := time.Now()

	res, err := s.db.Exec(`INSERT INTO documents (createdAt, updatedAt, title, value, accessKey, viewCount, maxViewCount, ttlMs) 
                           VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		now, now, title, value, accessKey, 0, maxViewCount, ttlMs)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *documentService) DeleteDocument(id int, accessKey string) (bool, error) {
	res, err := s.db.Exec(`DELETE FROM documents WHERE id = ? AND accessKey = ?`, id, accessKey)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func (s *documentService) GetDocument(id int) (*model.Document, error) {
	row := s.db.QueryRow(`SELECT id, createdAt, updatedAt, title, value, accessKey, viewCount, maxViewCount, ttlMs 
                          FROM documents WHERE id = ?`, id)

	var doc model.Document
	var title sql.NullString

	err := row.Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt, &title, &doc.Value, &doc.AccessKey, &doc.ViewCount, &doc.MaxViewCount, &doc.TTLMs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	doc.Title = &title.String
	return &doc, nil
}

func (s *documentService) DeleteExpiredDocuments() error {
	now := time.Now()
	_, err := s.db.Exec(`DELETE FROM documents WHERE ttlMs > 0 AND createdAt + ttlMs/1000 < ?`, now)
	return err
}
