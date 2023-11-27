package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	ImagesDir string
}

func (service *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}

	query := `
	SELECT user_id, title
	FROM galleries
	WHERE id = $1;
	`

	row := service.DB.QueryRow(query, gallery.ID)
	err := row.Scan(&gallery.UserID, &gallery.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	query := `
	SELECT id, title
	FROM galleries
	WHERE user_id = $1;
	`

	rows, err := service.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("query gallery by user id: %w", err)
	}

	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}

		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user id: %w", err)
		}

		galleries = append(galleries, gallery)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user id: %w", err)
	}

	return galleries, nil
}

func (service *GalleryService) Create(userID int, title string) (*Gallery, error) {
	gallery := Gallery{
		UserID: userID,
		Title:  title,
	}

	query := `
	INSERT INTO galleries (user_id, title)
	VALUES ($1, $2) RETURNING id;
	`

	row := service.DB.QueryRow(query, gallery.UserID, gallery.Title)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	query := `
	UPDATE galleries
	SET title = $2
	WHERE id = $1;
	`

	_, err := service.DB.Exec(query, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}

	return nil
}

func (service *GalleryService) Delete(gallery *Gallery) error {
	query := `
	DELETE FROM galleries
	WHERE id = $1;
	`

	_, err := service.DB.Exec(query, gallery.ID)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	return nil
}

func (service *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := service.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	return nil
}

func (service *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}

	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, service.extensions()) {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}

	return images, nil
}

func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(service.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath) // determine if file exists
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}

		return Image{}, fmt.Errorf("querying for image: %w", err)
	}

	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func (service *GalleryService) extensions() []string {
	return []string{
		".png",
		".jpg",
		".jpeg",
		".gif",
	}
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)

		if filepath.Ext(file) == ext {
			return true
		}
	}

	return false
}
