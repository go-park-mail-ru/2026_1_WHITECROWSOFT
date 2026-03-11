package notes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/mock"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNoteHandler_GetAllNotes(t *testing.T) {
	mockData := mock.NewMockData()
	handler := NewNoteHandler(mockData)
	userUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	tests := []struct {
		name             string
		setupContext     func() context.Context
		expectedStatus   int
		expectedError    string
		validateResponse func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful get all notes for user",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), types.UserIDKey, userUUID.String())
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				notes, ok := resp["notes"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, notes, 5) // В мок-данных 5 заметок для этого пользователя

				total, ok := resp["total"].(float64)
				assert.True(t, ok)
				assert.Equal(t, float64(5), total)
			},
		},
		{
			name: "unauthorized - no user id in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedStatus: http.StatusUnauthorized,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				err, ok := resp["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, jwt.ErrNoUserID.Error(), err)
			},
		},
		{
			name: "invalid user id format",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), types.UserIDKey, "invalid-uuid")
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				err, ok := resp["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, ErrInvalidNoteID.Error(), err)
			},
		},
		{
			name: "user with no notes - returns all notes (fallback)",
			setupContext: func() context.Context {
				newUserUUID := uuid.New()
				return context.WithValue(context.Background(), types.UserIDKey, newUserUUID.String())
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				notes, ok := resp["notes"].([]interface{})
				assert.True(t, ok)
				// Проверяем, что вернулись все заметки (fallback)
				assert.Len(t, notes, 5)

				total, ok := resp["total"].(float64)
				assert.True(t, ok)
				assert.Equal(t, float64(5), total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/notes", nil)
			req = req.WithContext(tt.setupContext())
			w := httptest.NewRecorder()

			handler.GetAllNotes(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			tt.validateResponse(t, response)
		})
	}
}

func TestNoteHandler_GetNote(t *testing.T) {
	mockData := mock.NewMockData()
	handler := NewNoteHandler(mockData)

	// Получаем первую заметку из мок-данных для тестов
	existingNote := mockData.Notes[0]
	nonExistentNoteID := uuid.New()

	tests := []struct {
		name             string
		path             string
		expectedStatus   int
		expectedError    string
		validateResponse func(*testing.T, map[string]interface{})
	}{
		{
			name:           "successful get note by id",
			path:           "/notes/" + existingNote.ID.String(),
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				// Проверяем заметку
				note, ok := resp["note"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, existingNote.ID.String(), note["ID"])
				assert.Equal(t, existingNote.Title, note["Title"])

				// Проверяем блоки
				blocks, ok := resp["blocks"].([]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, blocks)

				// Проверяем структуру первого блока
				firstBlock := blocks[0].(map[string]interface{})
				assert.Contains(t, firstBlock, "id")
				assert.Contains(t, firstBlock, "note_id")
				assert.Contains(t, firstBlock, "type_id")
				assert.Contains(t, firstBlock, "position")
				assert.Contains(t, firstBlock, "content")
				assert.Contains(t, firstBlock, "states")

				// Проверяем что states существуют
				states, ok := firstBlock["states"].([]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, states)
			},
		},
		{
			name:           "note id required - path too short",
			path:           "/notes",
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				err, ok := resp["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, ErrNoteIDRequired.Error(), err)
			},
		},
		{
			name:           "invalid note id format",
			path:           "/notes/invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				err, ok := resp["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, ErrInvalidNoteID.Error(), err)
			},
		},
		{
			name:           "note not found",
			path:           "/notes/" + nonExistentNoteID.String(),
			expectedStatus: http.StatusNotFound,
			validateResponse: func(t *testing.T, resp map[string]interface{}) {
				err, ok := resp["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, ErrNoteNotFound.Error(), err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.GetNote(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			tt.validateResponse(t, response)
		})
	}
}

// Тест для проверки корректности данных в ответе GetNote
func TestNoteHandler_GetNote_DataIntegrity(t *testing.T) {
	mockData := mock.NewMockData()
	handler := NewNoteHandler(mockData)

	// Берем конкретную заметку для теста
	testNote := mockData.Notes[0]

	req := httptest.NewRequest(http.MethodGet, "/notes/"+testNote.ID.String(), nil)
	w := httptest.NewRecorder()

	handler.GetNote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Проверяем соответствие данных заметки
	note, ok := response["note"].(map[string]interface{})
	assert.True(t, ok)

	// Проверяем все поля заметки
	assert.Equal(t, testNote.ID.String(), note["ID"])
	assert.Equal(t, testNote.UserID.String(), note["UserID"])
	assert.Equal(t, testNote.Title, note["Title"])
	assert.NotEmpty(t, note["CreatedAt"])
	assert.NotEmpty(t, note["UpdatedAt"])

	// Проверяем блоки
	blocks, ok := response["blocks"].([]interface{})
	assert.True(t, ok)

	// Проверяем, что блоки соответствуют заметке
	mockBlocks := mockData.GetBlocksByNoteID(testNote.ID)
	assert.Len(t, blocks, len(mockBlocks))

	// Проверяем соответствие каждого блока
	for i, blockData := range blocks {
		block := blockData.(map[string]interface{})
		mockBlock := mockBlocks[i]

		assert.Equal(t, mockBlock.ID.String(), block["id"])
		assert.Equal(t, mockBlock.NoteID.String(), block["note_id"])
		assert.Equal(t, float64(mockBlock.BlockTypeID), block["type_id"])
		assert.Equal(t, float64(mockBlock.Position), block["position"])
		assert.Equal(t, mockBlock.Content, block["content"])

		// Проверяем состояния блока
		states, ok := block["states"].([]interface{})
		assert.True(t, ok)

		mockStates := mockData.GetBlockStatesByBlockID(mockBlock.ID)
		assert.Len(t, states, len(mockStates))

		if len(mockStates) > 0 {
			firstState := states[0].(map[string]interface{})
			assert.Equal(t, mockStates[0].ID.String(), firstState["ID"])
			assert.Equal(t, mockStates[0].BlockID.String(), firstState["BlockID"])
			assert.Equal(t, mockStates[0].Formatting, firstState["Formatting"])
		}
	}
}

// Тест для проверки GetNote с разными вариантами путей
func TestNoteHandler_GetNote_PathVariations(t *testing.T) {
	mockData := mock.NewMockData()
	handler := NewNoteHandler(mockData)

	testNote := mockData.Notes[0]

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "path with trailing slash",
			path:           "/notes/" + testNote.ID.String() + "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "path with query parameters",
			path:           "/notes/" + testNote.ID.String() + "?format=full",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "path with extra segments",
			path:           "/notes/" + testNote.ID.String() + "/extra/segment",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.GetNote(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
