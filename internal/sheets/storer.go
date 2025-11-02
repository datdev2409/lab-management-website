package sheets

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ReportStorer interface {
	Store(ctx context.Context, dataReader io.Reader, reportName string) (string, error)
}

type LocalFileStoreStrategy struct {
	BaseDir string
}

func (l *LocalFileStoreStrategy) Store(ctx context.Context, dataReader io.Reader, reportName string) (string, error) {
	if err := os.MkdirAll(l.BaseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to ensure base directory exists: %w", err)
	}

	fullPath := filepath.Join(l.BaseDir, reportName)

	destFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file %s: %w", fullPath, err)
	}
	defer destFile.Close()

	bytesWritten, err := io.Copy(destFile, dataReader)
	if err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to copy report stream to file: %w", err)
	}

	fmt.Printf("Successfully wrote %d bytes to local file: %s\n", bytesWritten, fullPath)

	return fullPath, nil
}
