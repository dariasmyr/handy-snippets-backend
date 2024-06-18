package services

import (
	"database/sql"
	"errors"
	"fmt"
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

	fmt.Println(now, now, title, value, accessKey, 0, maxViewCount, ttlMs)

	res, err := s.db.Exec(`INSERT INTO documents (created_at, updated_at, title, value, access_key, view_count, max_view_count, ttl_ms) 
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
	res, err := s.db.Exec(`DELETE FROM documents WHERE id = ? AND access_key = ?`, id, accessKey)
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
	row := s.db.QueryRow(`SELECT id, created_at, updated_at, title, value, access_key, view_count, max_view_count, ttl_ms 
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
	_, err := s.db.Exec(`DELETE FROM documents WHERE ttl_ms > 0 AND created_at + ttl_ms/1000 < ?`, now)
	return err
}
