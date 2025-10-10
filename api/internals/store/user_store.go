package store

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainPasswordText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPasswordText), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plainPasswordText
	p.hash = hash

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type postgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *postgresUserStore {
	return &postgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUserName(username string) (*User, error)
	UpdateUser(*User) error
}

func (pg *postgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username,email,password_hash,bio) 
	VALUES ($1,$2,$3,$4) 
	RETURNING id,created_at,updated_at
	`
	err := pg.db.QueryRow(query, user.UserName, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (pg *postgresUserStore) GetUserByUserName(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	SELECT id,username,email,password_hash,bio, created_at, updated_at FROM users 
	where username = $1
	`
	err := pg.db.QueryRow(query, username).Scan(user.ID, user.UserName, user.Email, user.PasswordHash.hash, user.Bio, user.CreatedAt, user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (pg *postgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users 
	SET username = $1, email = $2,  bio = $3, updated_at = CURRENT_TIMESTAMP
	WHERE id = $4
	RETURNING updated_at
	`
	res, err := pg.db.Exec(query, user.UserName, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
