package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/danakin/web-dev-with-go-2-code_along/rand"
)

const (
	MinBytesPerToken = 32 // Token should be a minimum of 32 bytes
)

type Session struct {
	ID        int
	UserId    int
	Token     string // only set when creating a new session, not filled when looking up a session
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB // TODO: Redis?
	BytesPerToken int     // How many bytes are used to generate session token
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(ss.BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("session token create: %w", err)
	}
	session := Session{
		UserId:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	query := `
	INSERT INTO sessions (user_id, token_hash)
	VALUES ($1, $2)
	RETURNING id;
	`
	row := ss.DB.QueryRow(query, session.UserId, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create session error: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// convert tokenHash<Array> to tokenHash<Slice>
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
