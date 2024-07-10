package services

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt" // Import bcrypt
	"pastebin/graph/model"
	"time"
)

type documentService struct {
	db *sql.DB
}

func NewDocumentService(db *sql.DB) DocumentService {
	return &documentService{db: db}
}

func (s *documentService) StartExpiredDocumentsCleaner(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			err := s.DeleteExpiredDocuments()
			if err != nil {
				fmt.Println("Error deleting expired documents:", err)
			}
			err = s.DeleteDocumentsWithMaxViews()
			if err != nil {
				fmt.Println("Error deleting documents with max views:", err)
			}
		}
	}()
}

func (s *documentService) CreateDocument(value string, accessKey string, maxViewCount, ttlMs int) (int, error) {
	now := time.Now()

	hashedAccessKey, err := bcrypt.GenerateFromPassword([]byte(accessKey), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	res, err := s.db.Exec(`INSERT INTO documents (createdAt, updatedAt, value, accessKey, viewCount, maxViewCount, ttlMs) 
                           VALUES (?, ?, ?, ?, ?, ?, ?)`,
		now, now, value, string(hashedAccessKey), 0, maxViewCount, ttlMs)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *documentService) UpdateDocument(id int, value *string, accessKey string, maxViewCount, ttlMs int) (bool, error) {
	println("id", id)
	var storedHashedAccessKey string
	var createdAt time.Time
	var viewCount int

	err := s.db.QueryRow(`SELECT accessKey, createdAt, viewCount FROM documents WHERE id = ?`, id).Scan(&storedHashedAccessKey, &createdAt, &viewCount)
	if err != nil {
		return false, err
	}

	if ttlMs > 0 && time.Since(createdAt).Milliseconds() > int64(ttlMs) {
		_, err := s.db.Exec(`DELETE FROM documents WHERE id = ?`, id)
		if err != nil {
			return false, err
		}
		return false, errors.New("document expired and has been deleted")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedAccessKey), []byte(accessKey))
	if err != nil {
		return false, errors.New("access denied")
	}

	_, err = s.db.Exec(`UPDATE documents SET updatedAt = ?, value = ?, maxViewCount = ?, ttlMs = ? WHERE id = ?`,
		time.Now(), value, maxViewCount, ttlMs, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *documentService) DeleteDocument(id int, accessKey string) (bool, error) {
	var storedHashedAccessKey string
	var ttlMs int
	var createdAt time.Time

	err := s.db.QueryRow(`SELECT accessKey, createdAt, ttlMs FROM documents WHERE id = ?`, id).Scan(&storedHashedAccessKey, &createdAt, &ttlMs)
	if err != nil {
		return false, err
	}

	if ttlMs > 0 && time.Since(createdAt).Milliseconds() > int64(ttlMs) {
		_, err := s.db.Exec(`DELETE FROM documents WHERE id = ?`, id)
		if err != nil {
			return false, err
		}
		return false, errors.New("document expired and has been deleted")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedAccessKey), []byte(accessKey))
	if err != nil {
		return false, errors.New("access denied")
	}

	res, err := s.db.Exec(`DELETE FROM documents WHERE id = ?`, id)
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
	var storedHashedAccessKey string
	row := s.db.QueryRow(`SELECT id, createdAt, updatedAt, value, accessKey, viewCount, maxViewCount, ttlMs 
                          FROM documents WHERE id = ?`, id)

	var doc model.Document

	err := row.Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt, &doc.Value, &storedHashedAccessKey, &doc.ViewCount, &doc.MaxViewCount, &doc.TTLMs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("document not found")
		}
		return nil, err
	}

	if doc.TTLMs > 0 && time.Since(doc.CreatedAt).Milliseconds() > int64(doc.TTLMs) {
		_, err := s.db.Exec(`DELETE FROM documents WHERE id = ?`, doc.ID)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("document expired and has been deleted")
	}

	if doc.MaxViewCount == -1 || doc.ViewCount < doc.MaxViewCount {
		doc.ViewCount++
		_, err = s.db.Exec(`UPDATE documents SET viewCount = ? WHERE id = ?`, doc.ViewCount, doc.ID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("max view count reached")
	}

	return &doc, nil
}

func (s *documentService) DeleteExpiredDocuments() error {
	now := time.Now().Unix()
	_, err := s.db.Exec(`DELETE FROM documents WHERE ttlMs > 0 AND (createdAt + ttlMs / 1000) < ?`, now)
	return err
}

func (s *documentService) DeleteDocumentsWithMaxViews() error {
	_, err := s.db.Exec(`DELETE FROM documents WHERE maxViewCount > 0 AND viewCount >= maxViewCount`)
	return err
}
