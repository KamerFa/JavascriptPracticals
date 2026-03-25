package prefs

import (
	"database/sql"
	"encoding/json"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// EnsureSession creates a session row if it doesn't exist.
func (s *Store) EnsureSession(sessionID string) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO sessions (id) VALUES (?)`, sessionID,
	)
	return err
}

// GetPrefs returns preferences for a session, creating defaults if needed.
func (s *Store) GetPrefs(sessionID string) (*Prefs, error) {
	if err := s.EnsureSession(sessionID); err != nil {
		return nil, err
	}

	row := s.db.QueryRow(
		`SELECT rev_index, filter, langs, seerah_open, seerah_lang FROM prefs WHERE session_id = ?`,
		sessionID,
	)

	var p Prefs
	var langsJSON string
	var seerahOpen int

	err := row.Scan(&p.RevIndex, &p.Filter, &langsJSON, &seerahOpen, &p.SeerahLang)
	if err == sql.ErrNoRows {
		// Return defaults
		return &Prefs{
			RevIndex:   0,
			Filter:     "all",
			Langs:      Langs{Arabic: true},
			SeerahOpen: false,
			SeerahLang: "en",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(langsJSON), &p.Langs)
	p.SeerahOpen = seerahOpen == 1
	return &p, nil
}

// SavePrefs upserts preferences for a session.
func (s *Store) SavePrefs(sessionID string, p *Prefs) error {
	if err := s.EnsureSession(sessionID); err != nil {
		return err
	}

	langsJSON, _ := json.Marshal(p.Langs)
	seerahOpen := 0
	if p.SeerahOpen {
		seerahOpen = 1
	}

	_, err := s.db.Exec(`
		INSERT INTO prefs (session_id, rev_index, filter, langs, seerah_open, seerah_lang, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(session_id) DO UPDATE SET
			rev_index = excluded.rev_index,
			filter = excluded.filter,
			langs = excluded.langs,
			seerah_open = excluded.seerah_open,
			seerah_lang = excluded.seerah_lang,
			updated_at = CURRENT_TIMESTAMP
	`, sessionID, p.RevIndex, p.Filter, string(langsJSON), seerahOpen, p.SeerahLang)

	return err
}

// GetBookmark returns the pinned bookmark for a session (only one active pin).
func (s *Store) GetBookmark(sessionID string) (*Bookmark, error) {
	row := s.db.QueryRow(
		`SELECT id, surah_num, ayah_num, rev_index, surah_name FROM bookmarks WHERE session_id = ? ORDER BY created_at DESC LIMIT 1`,
		sessionID,
	)

	var b Bookmark
	err := row.Scan(&b.ID, &b.SurahNum, &b.AyahNum, &b.RevIndex, &b.SurahName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// SetBookmark replaces the current pin with a new one (or removes it).
func (s *Store) SetBookmark(sessionID string, b *Bookmark) error {
	if err := s.EnsureSession(sessionID); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing bookmarks for this session
	if _, err := tx.Exec(`DELETE FROM bookmarks WHERE session_id = ?`, sessionID); err != nil {
		return err
	}

	if b != nil {
		_, err = tx.Exec(
			`INSERT INTO bookmarks (session_id, surah_num, ayah_num, rev_index, surah_name) VALUES (?, ?, ?, ?, ?)`,
			sessionID, b.SurahNum, b.AyahNum, b.RevIndex, b.SurahName,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteBookmark removes the pin for a session.
func (s *Store) DeleteBookmark(sessionID string) error {
	_, err := s.db.Exec(`DELETE FROM bookmarks WHERE session_id = ?`, sessionID)
	return err
}
