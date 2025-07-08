package storage

import (
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalUploaderAdapter saves files to the local disk.
type LocalUploaderAdapter struct {
	basePath string
}

// NewLocalUploaderAdapter creates a new local uploader.
func NewLocalUploaderAdapter(basePath string) *LocalUploaderAdapter {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		slog.Error("could not create local upload directory", "error", err)
		os.Exit(1)
	}
	return &LocalUploaderAdapter{basePath: basePath}
}

// Upload implements the StorageAdapter interface.
func (a *LocalUploaderAdapter) Upload(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := filepath.Join(a.basePath, objectName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	slog.Info("Successfully saved file to local disk", "path", dstPath)
	publicURL := "/" + filepath.ToSlash(dstPath)
	return publicURL, nil
}
