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

	// Approaches to unique user id in session table
	// 1
	//   1. Query for a user's session
	//   2. If found, update the user's session
	//   3. If not found, create a new session for the user
	// 2 implemented below
	//   1. Try to update session
	//   2. If err, create new session
	query := `
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;
	`
	row := ss.DB.QueryRow(query, session.UserId, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		query := `
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;
		`
		row := ss.DB.QueryRow(query, session.UserId, session.TokenHash)
		err = row.Scan(&session.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("create session error: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)

	// get session
	var user User
	query := `
		SELECT user_id
		FROM sessions 
		WHERE token_hash = $1;
	`
	row := ss.DB.QueryRow(query, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("sessionservice user: %w", err)
	}

	// get user
	query = `
		SELECT email
		FROM users
		WHERE id = $1;
	`
	row = ss.DB.QueryRow(query, user.ID)
	err = row.Scan(&user.Email)
	if err != nil {
		return nil, fmt.Errorf("sessionservice user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	query := `
		DELETE FROM sessions
		WHERE token_hash = $1;
	`
	_, err := ss.DB.Exec(query, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// convert tokenHash<Array> to tokenHash<Slice>
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
