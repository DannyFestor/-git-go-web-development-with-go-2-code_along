package models

import "database/sql"

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
	// TODO
	// 1. Create Session Token
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO
	return nil, nil
}
