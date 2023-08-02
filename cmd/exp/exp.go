package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "favorite-color", "orange")

	val := ctx.Value("favorite-color")
	fmt.Println(val)
}
