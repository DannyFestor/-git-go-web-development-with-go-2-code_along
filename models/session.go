package models

import (
	"database/sql"
	"fmt"

	"github.com/danakin/web-dev-with-go-2-code_along/rand"
)

type Session struct {
	ID        int
	UserId    int
	Token     string // only set when creating a new session, not filled when looking up a session
	TokenHash string
}

type SessionService struct {
	DB *sql.DB // TODO: Redis?
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, err := rand.SessionToken()
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
