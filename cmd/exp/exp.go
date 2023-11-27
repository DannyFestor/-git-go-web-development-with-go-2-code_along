package main

import (
	"fmt"

	"github.com/danakin/web-dev-with-go-2-code_along/models"
)

func main() {
	gs := models.GalleryService{}
	fmt.Println(gs.Images(2))
}
