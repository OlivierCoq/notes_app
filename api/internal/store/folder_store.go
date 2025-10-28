package store

import (
	"database/sql"
	"fmt"
)

type Folder struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	UserID         int           `json:"user_id"`
	IsFavorite     bool          `json:"is_favorite"`
	ParentFolderID sql.NullInt64 `json:"parent_folder_id"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
}

type PostgresFolderStore struct {
	db *sql.DB
}

func NewPostgresFolderStore(db *sql.DB) *PostgresFolderStore {
	return &PostgresFolderStore{db: db}
}

// Interface for FolderStore to allow decoupling and easier testing:
type FolderStore interface {
	CreateFolder(*Folder) (*Folder, error)
	GetFolderByID(id int) (*Folder, error)
	UpdateFolder(*Folder) error
	DeleteFolder(id int) error
	ListFoldersByUserID(userID int) ([]*Folder, error)
	GetFolderOwner(id int) (int, error)
}

// CRUD operations:

// Create folder:
func (pg *PostgresFolderStore) CreateFolder(folder *Folder) (*Folder, error) {

	// transaction:
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO folders (title, user_id, is_favorite, parent_folder_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(query, folder.Title, folder.UserID, folder.IsFavorite, folder.ParentFolderID).Scan(&folder.ID, &folder.CreatedAt, &folder.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return folder, nil
}

func (pg *PostgresFolderStore) GetFolderByID(id int) (*Folder, error) {
	// Create a folder instance first:
	folder := &Folder{}
	query := `
		SELECT id, title, user_id, is_favorite, parent_folder_id, created_at, updated_at
		FROM folders
		WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(
		&folder.ID,
		&folder.Title,
		&folder.UserID,
		&folder.IsFavorite,
		&folder.ParentFolderID,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Folder not found
	}

	return folder, nil
}

func (pg *PostgresFolderStore) UpdateFolder(folder *Folder) error {

	// Transaction
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE folders
		SET title = $1,
		    is_favorite = $2,
		    parent_folder_id = $3,
		    updated_at = NOW()
		WHERE id = $4
	`
	_, err = tx.Exec(query, folder.Title, folder.IsFavorite, folder.ParentFolderID, folder.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (pg *PostgresFolderStore) DeleteFolder(id int) error {

	query := `
		DELETE FROM folders
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
		return sql.ErrNoRows // Folder not found
	}
	return nil
}

func (pg *PostgresFolderStore) GetFolderOwner(id int) (int, error) {
	var userID int
	query := `
		SELECT user_id
		FROM folders
		WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no folder found with id %d", id)
		}
		return 0, err
	}
	return userID, nil
}

func (pg *PostgresFolderStore) ListFoldersByUserID(userID int) ([]*Folder, error) {
	query := `
		SELECT id, title, user_id, is_favorite, parent_folder_id, created_at, updated_at
		FROM folders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := pg.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*Folder
	for rows.Next() {
		folder := &Folder{}
		err := rows.Scan(
			&folder.ID,
			&folder.Title,
			&folder.UserID,
			&folder.IsFavorite,
			&folder.ParentFolderID,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return folders, nil
}
