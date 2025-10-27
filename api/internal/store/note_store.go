package store

import (
	"database/sql"
	"fmt"
)


type Note struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	UserID     int    `json:"user_id"`
	IsFavorite bool   `json:"is_favorite"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PostgresNoteStore struct {
	db *sql.DB
}

func NewPostgresNoteStore(db *sql.DB) *PostgresNoteStore {
	return &PostgresNoteStore{db: db}
}

// Interface for NoteStore to allow decoupling and easier testing:
type NoteStore interface {
	CreateNote(*Note) (*Note, error)
	GetNoteByID(id int) (*Note, error)
	UpdateNote(*Note) error
	DeleteNote(id int) error
	GetNoteOwner(id int) (int, error)
}

// CRUD operations:

	// Create note:
func (pg *PostgresNoteStore) CreateNote(note *Note) (*Note, error) {

	// transaction:
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO notes (title, content, user_id, is_favorite, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	err = tx.QueryRow(query, note.Title, note.Content, note.UserID, note.IsFavorite).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return note, nil
}

func (pg *PostgresNoteStore) GetNoteByID(id int) (*Note, error) {
	// Create a note instance first:
	note := &Note{}
	query := `
		SELECT id, title, content, user_id, is_favorite, created_at, updated_at
		FROM notes	
		WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.UserID,
		&note.IsFavorite,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Note not found
	}

	return note, nil
}

func (pg *PostgresNoteStore) UpdateNote(note *Note) error {

	// Transaction
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE notes
		SET title = $1,
		    content = $2,
		    is_favorite = $3,
		    updated_at = NOW()
		WHERE id = $4
	`
	_, err = tx.Exec(query, note.Title, note.Content, note.IsFavorite, note.ID)
	if err != nil {
		return err
	}
	
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (pg *PostgresNoteStore) DeleteNote(id int) error {

	query := `
		DELETE FROM notes
		WHERE id = $1
	`
	res, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Note not found
	}
	return nil


}

func (pg *PostgresNoteStore) GetNoteOwner(id int) (int, error) {
	var userID int
	query := `
		SELECT user_id
		FROM notes
		WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no note found with id %d", id)
		}
		return 0, err
	}
	return userID, nil	
}