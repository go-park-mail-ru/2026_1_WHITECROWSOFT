package mock

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/models"
	"github.com/google/uuid"
)

type MockData struct {
	Notes  []models.Note
	Blocks []models.Block
}

func NewMockData() *MockData {
	mock := &MockData{}
	mock.init()
	return mock
}

func (m *MockData) init() {
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	for i := 1; i <= 5; i++ {
		noteID := uuid.New()

		blockID := uuid.New()
		state := map[string]interface{}{
			"text":  "Это пример текста для заметки",
			"title": "Заголовок mock",
			"content": []string{
				"Первый параграф",
				"Второй параграф",
				"Третий параграф",
			},
			"format":  "text",
			"tags":    []string{"пример", "тест", "мок"},
			"created": time.Now().AddDate(0, -i, 0).Format(time.RFC3339),
			"updated": time.Now().Format(time.RFC3339),
		}

		stateJSON, _ := json.Marshal(state)

		block := models.Block{
			ID:     blockID,
			NoteID: noteID,
			Type:   "text",
			State:  stateJSON,
		}
		m.Blocks = append(m.Blocks, block)

		note := models.Note{
			ID:       noteID,
			UserID:   userID,
			Title:    "Моя заметка" + strconv.Itoa(i),
			ParentID: nil,
			Blocks:   []uuid.UUID{blockID},
		}
		m.Notes = append(m.Notes, note)
	}
}
