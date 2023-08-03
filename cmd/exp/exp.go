package main

import (
	stdctx "context"
	"fmt"

	"github.com/danakin/web-dev-with-go-2-code_along/context"
	"github.com/danakin/web-dev-with-go-2-code_along/models"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := stdctx.Background()
	user := models.User{
		Email: "danny@festor.info",
	}
	ctx = context.WithUser(ctx, &user)

	retrievedUser := context.User(ctx)

	fmt.Println(retrievedUser.Email)
}
