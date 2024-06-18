package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table for documents
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS documents (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            createdAt DATETIME NOT NULL,
            updatedAt DATETIME NOT NULL,
            title TEXT,
            value TEXT NOT NULL,
            accessKey TEXT NOT NULL,
            viewCount INTEGER NOT NULL DEFAULT 0,
            maxViewCount INTEGER NOT NULL DEFAULT -1,
            ttlMs INTEGER NOT NULL DEFAULT -1
        );
    `)
	if err != nil {
		return nil, err
	}

	// Additional initialization if needed

	return db, nil
}
