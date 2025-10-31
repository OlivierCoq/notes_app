package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

// Password hashing and verification methods can be added here
/*
	In software development, a salt is a value that determins the computational complexity of the hashing algorithm.
	It is used to make password hashing more secure by increasing the time it takes to compute the hash.
	The higher the cost, the more computationally expensive it is to hash a password, which makes it harder for attackers to brute-force or crack passwords.

	Bcrypt uses a cost parameter that determines how many rounds of hashing are performed.
	The default cost is 10, but it can be increased for added security.
	However, increasing the cost also increases the time it takes to hash a password, so it's important to find a balance between security and performance.
*/
func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}
	p.plainText = &plainText
	p.hash = hash
	return nil
}

// For verifying if the provided plaintext password matches the stored hash. Authentication and login purposes.
func (p *password) Matches(plainText string) (bool, error) {
	if p == nil || p.hash == nil {
		return false, errors.New("password hash is nil")
	}
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err
			}
		}
	}
	return true, nil
}

type User struct {
	ID             int      `json:"id"`
	Username       string   `json:"username"`
	Email          string   `json:"email"`
	PasswordHash   password `json:"-"`
	Bio            string   `json:"bio"`
	AuthLevel      int      `json:"auth_level"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	PfpURL         string   `json:"pfp_url"`
	AddressLine1   string   `json:"address_line1"`
	AddressLine2   string   `json:"address_line2"`
	AddressCity    string   `json:"city"`
	AddressState   string   `json:"state"`
	AddressZip     string   `json:"zip"`
	AddressCountry string   `json:"country"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnonymousUser is a placeholder for unauthenticated users.
var AnonymousUser = &User{}

// IsAnonymous checks if the user is anonymous.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

// Interface for UserStore to allow decoupling and easier testing:
type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserById(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) (*User, error)
	GetUserToken(scope, tokenPlaintext string) (*User, error)
	UpdateUserPassword(userID int, newPassword string) error
}

// CRUUD operations:

// Create user:
func (s *PostgresUserStore) CreateUser(user *User) (*User, error) {

	query := `
		INSERT INTO users (
		username, 
		email, 
		password_hash, 
		bio, 
		auth_level, 
		first_name, 
		last_name, 
		pfp_url, 
		address_line1, 
		address_line2, 
		address_city, 
		address_state, 
		address_zip_code, 
		address_country, 
		created_at, 
		updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash.hash,
		user.Bio,
		user.AuthLevel,
		user.FirstName,
		user.LastName,
		user.PfpURL,
		user.AddressLine1,
		user.AddressLine2,
		user.AddressCity,
		user.AddressState,
		user.AddressZip,
		user.AddressCountry,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Read (Get) user by ID:
func (s *PostgresUserStore) GetUserById(id int) (*User, error) {
	query := `
		SELECT id, 
		username, 
		email, 
		password_hash, 
		bio, 
		auth_level, 
		first_name, 
		last_name, 
		pfp_url, 
		address_line1, 
		address_line2, 
		address_city, 
		address_state, 
		address_zip_code, 
		address_country, 
		created_at, 
		updated_at
		FROM users
		WHERE id = $1
	`
	user := &User{
		PasswordHash: password{},
	}
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.AuthLevel,
		&user.FirstName,
		&user.LastName,
		&user.PfpURL,
		&user.AddressLine1,
		&user.AddressLine2,
		&user.AddressCity,
		&user.AddressState,
		&user.AddressZip,
		&user.AddressCountry,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

// Read (Get) user by Username:
func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	query := `
		SELECT id,
		username,
		email,
		password_hash,
		bio,
		auth_level,
		first_name,
		last_name,
		pfp_url,
		address_line1,
		address_line2,
		address_city,
		address_state,
		address_zip_code,
		address_country,
		created_at,
		updated_at
		FROM users
		WHERE username = $1
	`
	user := &User{
		PasswordHash: password{},
	}
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.AuthLevel,
		&user.FirstName,
		&user.LastName,
		&user.PfpURL,
		&user.AddressLine1,
		&user.AddressLine2,
		&user.AddressCity,
		&user.AddressState,
		&user.AddressZip,
		&user.AddressCountry,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update user:
func (s *PostgresUserStore) UpdateUser(user *User) (*User, error) {
	query := `
		UPDATE users
		SET username = $1, 
		email = $2, 
		bio = $3, 
		auth_level = $4, 
		first_name = $5, 
		last_name = $6, 
		pfp_url = $7, 
		address_line1 = $8, 
		address_line2 = $9, 
		address_city = $10, 
		address_state = $11, 
		address_zip_code = $12, 
		address_country = $13, 
		updated_at = NOW()
		WHERE id = $14
	`
	result, err := s.db.Exec(
		query,
		user.Username,
		user.Email,
		user.Bio,
		user.AuthLevel,
		user.FirstName,
		user.LastName,
		user.PfpURL,
		user.AddressLine1,
		user.AddressLine2,
		user.AddressCity,
		user.AddressState,
		user.AddressZip,
		user.AddressCountry,
		user.ID)
	if err != nil {
		return nil, err
	}
	// rows:
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	// Check if any row was actually updated
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

// Get user by token (for authentication):
func (s *PostgresUserStore) GetUserToken(scope, plaintextPassword string) (*User, error) {
	// Implementation for retrieving a user by token from PostgreSQL

	tokenHash := sha256.Sum256([]byte(plaintextPassword))

	// INNER JOIN tokens t ON u.id = t.user_id (Not sure if order matters here)
	query := `
		SELECT 
		u.id, 
		u.username, 
		u.email, 
		u.password_hash, 
		u.bio, 
		u.auth_level, 
		u.first_name, 
		u.last_name, 
		u.pfp_url, 
		u.address_line1, 
		u.address_line2, 
		u.address_city, 
		u.address_state, 
		u.address_zip_code, 
		u.address_country, 
		u.created_at, 
		u.updated_at
		FROM users u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3
	`
	// the t.expiry is to ensure the token is still valid (not expired)
	user := &User{
		PasswordHash: password{},
	}
	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.AuthLevel,
		&user.FirstName,
		&user.LastName,
		&user.PfpURL,
		&user.AddressLine1,
		&user.AddressLine2,
		&user.AddressCity,
		&user.AddressState,
		&user.AddressZip,
		&user.AddressCountry,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update user password:
func (s *PostgresUserStore) UpdateUserPassword(userID int, newPassword string) error {
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`
	result, err := s.db.Exec(query, newPasswordHash, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
