package sheets

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"go.uber.org/zap"
)

func ConvertExcelToPDF(ctx context.Context, filename string) (string, error) {
	gotenbergURL := os.Getenv("GOTENBERG_URL")
	if gotenbergURL == "" {
		gotenbergURL = "http://localhost:3000"
	}

	url := gotenbergURL + "/forms/libreoffice/convert"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(filename)
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to open Excel file for PDF conversion", zap.String("filename", filename), zap.Error(err))
		return "", err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("files", filepath.Base(filename))
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to create form file for PDF conversion", zap.String("filename", filename), zap.Error(err))
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to copy file content for PDF conversion", zap.String("filename", filename), zap.Error(err))
		return "", err
	}

	err = writer.Close()
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to close multipart writer for PDF conversion", zap.Error(err))
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to create HTTP request for PDF conversion", zap.String("url", url), zap.Error(err))
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to execute HTTP request for PDF conversion", zap.String("url", url), zap.Error(err))
		return "", err
	}
	defer res.Body.Close()

	// Save the PDF to a file
	outFile, err := os.Create(strings.Replace(filename, ".xlsx", ".pdf", 1))
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to create PDF output file", zap.String("filename", filename), zap.Error(err))
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		logger.FromCtx(ctx).Debug("Failed to copy PDF response to output file", zap.String("filename", outFile.Name()), zap.Error(err))
		return "", err
	}

	return outFile.Name(), nil
}
