package sheets

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ConvertExcelToPDF(filename string) (string, error) {

	url := "http://localhost:3000/forms/libreoffice/convert"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("files", filepath.Base(filename))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Save the PDF to a file
	outFile, err := os.Create(strings.Replace(filename, ".xlsx", ".pdf", 1))
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		return "", err
	}

	return outFile.Name(), nil
}
