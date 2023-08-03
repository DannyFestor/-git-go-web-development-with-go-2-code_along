package context

import (
	"context"

	"github.com/danakin/web-dev-with-go-2-code_along/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User) // assert value is a User model
	if !ok {                       // nothing in the context, or invalid value
		return nil
	}

	return user
}
