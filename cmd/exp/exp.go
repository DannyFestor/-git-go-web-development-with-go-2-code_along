package main

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "orange")

	val := ctx.Value(favoriteColorKey)
	strVal, ok := val.(string) // type assertion
	if !ok {
		fmt.Println("it isn't an int")
		return
	}

	fmt.Println(strVal)
	fmt.Println(strings.HasPrefix(strVal, "or"))
}
