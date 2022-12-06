package models

import (
	"database/sql"
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
	// TODO: Hash the Session Token
	session := Session{
		UserId: userID,
		Token:  token,
		// TODO: Set the TokenHash
	}
	// TODO: Store Session in DP
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO
	return nil, nil
}
