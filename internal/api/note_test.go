package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HiogoPariz/api-notez/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockNoteRepository struct{}

func (m *MockNoteRepository) CreateNote(note *dto.NoteDTO) error {
	return nil
}

func (m *MockNoteRepository) GetNotes() ([]*dto.NoteDTO, error) {
	return nil, nil
}

func (m *MockNoteRepository) GetNoteByID(id int) (*dto.NoteDTO, error) {
	return nil, nil
}

func (m *MockNoteRepository) GetNoteByUserId(userId int) (*dto.NoteListById, error) {
	return nil, nil
}

func (m *MockNoteRepository) DeleteNote(id int) error {
	return nil
}

func (m *MockNoteRepository) UpdateNote(note *dto.NoteDTO) error {
	return nil
}

func TestCreateNote(t *testing.T) {
	mockRepo := &MockNoteRepository{}

	noteService := &NoteService{
		repo: mockRepo,
	}

	router := gin.Default()
	router.POST("/note", noteService.CreateNote)

	requestBody := dto.NoteRequest{
		Title:   "Test Note",
		Content: "This is a test note.",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/note", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "fileName")

	fileName := response["fileName"].(string)
	fmt.Println("Input hash", fileName)
	assert.Regexp(t, `^[a-f0-9]{16}$`, fileName)

	time.Sleep(1 * time.Millisecond)
	req2, _ := http.NewRequest(http.MethodPost, "/note", bytes.NewBuffer(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	var response2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response2)
	fileName2 := response2["fileName"].(string)
	assert.NotEqual(t, fileName, fileName2)
}
