package uploader

import (
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"venturo-core/internal/adapter/storage"
)

// FileUploader handles the low-level logic of saving and uploading files.
type FileUploader struct {
	storageAdapter storage.StorageAdapter
	localPath      string
}

// NewFileUploader creates a new uploader instance.
func NewFileUploader(storageAdapter storage.StorageAdapter, localPath string) *FileUploader {
	// Ensure the temporary upload directory exists
	if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
		slog.Error("could not create temp upload directory", "error", err)
		os.Exit(1)
	}
	return &FileUploader{storageAdapter: storageAdapter, localPath: localPath}
}

// UploadAsync handles the entire file processing pipeline in the background.
func (u *FileUploader) UploadAsync(
	file *multipart.FileHeader,
	objectName string,
	onLocalUploadSuccess func(),
	onCloudUploadSuccess func(),
) {
	// --- Step 1: Save Locally ---
	localFilePath := filepath.Join(u.localPath, objectName)
	if err := saveToLocal(file, localFilePath); err != nil {
		slog.Error("Error saving temp file", "file", objectName, "error", err)
		return
	}
	slog.Info("Successfully saved temp file", "path", localFilePath)

	// Execute the first callback
	onLocalUploadSuccess()

	// --- Step 2: Upload to Cloud ---
	// In a real app, you would read the local file's contents to upload.
	// We pass the header for placeholder simplicity.
	if _, err := u.storageAdapter.Upload(context.Background(), file, objectName); err != nil {
		slog.Error("Error uploading to cloud", "file", objectName, "error", err)
		return // Don't continue if cloud upload fails
	}
	slog.Info("Successfully uploaded to cloud", "file", objectName)

	// Execute the second callback
	onCloudUploadSuccess()

	// --- Step 3: Cleanup ---
	if err := os.Remove(localFilePath); err != nil {
		slog.Error("Error cleaning up temp file", "file", objectName, "error", err)
	} else {
		slog.Info("Successfully cleaned up temp file", "file", objectName)
	}
}

// saveToLocal is a helper function containing the file-saving logic.
func saveToLocal(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
