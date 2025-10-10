package store

import (
	"database/sql"
	"time"

	"github.com/codemaverick07/api/internals/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

func (p *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1,$2,$3,$4)
	`
	_, err := p.db.Exec(query, token.Hash, token.UserId, token.Expiry, token.Scope)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresTokenStore) CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = p.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (p *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2 
	`
	_, err := p.db.Exec(query, scope, userID)

	return err
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}
