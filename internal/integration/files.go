package integration

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/HiogoPariz/api-notez/internal/dto"
)

type FileIntegration struct {
	dto *dto.NoteDTO
}

func CreateFileIntegration(dto *dto.NoteDTO) *FileIntegration {
	return &FileIntegration{dto}
}

func (integration *FileIntegration) GetFileContent() (string, error) {
	client := http.Client{}
	requestURL := fmt.Sprintf("http://localhost:3001/%s", integration.dto.FileName)
	requestBody := []byte{}
	bodyReader := bytes.NewReader(requestBody)

	req, err := http.NewRequest(http.MethodGet, requestURL, bodyReader)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func (integration *FileIntegration) CreateFileContent(content string, file_name string) error {
	client := http.Client{}
	requestURL := fmt.Sprintf("http://localhost:3001/%s", file_name)
	requestBody := []byte(fmt.Sprintf(`{"content": %q}`, content))
	bodyReader := bytes.NewReader(requestBody)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return err
	}

	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil
}
