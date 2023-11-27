package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

// io.Reader keeps track of where in a file it is, so read bytes get lost when handling it afterwards
// while an io.ReadSeeker can be reset so it starts back at the beginning
func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes) // read first 512 bytes of a file because MIME information should be in those
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	_, err = r.Seek(0, 0) // seek to the beginning of the file for handling it afterwards
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, t := range allowedTypes {
		if contentType == t {
			return nil
		}
	}

	return FileError{
		Issue: fmt.Sprintf("invalid content type: %v", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename)),
		}
	}

	return nil
}
