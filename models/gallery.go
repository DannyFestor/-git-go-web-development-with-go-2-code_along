package models

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
)

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

func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}
