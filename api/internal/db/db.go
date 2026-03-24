package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id         TEXT PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS prefs (
		session_id    TEXT PRIMARY KEY REFERENCES sessions(id),
		rev_index     INTEGER DEFAULT 0,
		filter        TEXT DEFAULT 'all',
		langs         TEXT DEFAULT '{"arabic":true,"transliteration":false,"bosnian":false,"english":false,"turkish":false}',
		seerah_open   INTEGER DEFAULT 0,
		seerah_lang   TEXT DEFAULT 'en',
		updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS bookmarks (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL REFERENCES sessions(id),
		surah_num  INTEGER NOT NULL,
		ayah_num   INTEGER NOT NULL,
		rev_index  INTEGER NOT NULL,
		surah_name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(session_id, surah_num, ayah_num)
	);
	`
	_, err := db.Exec(schema)
	return err
}
