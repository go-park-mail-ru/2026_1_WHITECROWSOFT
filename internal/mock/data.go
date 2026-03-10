package mock

import (
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/models"
	"github.com/google/uuid"
)

type MockData struct {
	Notes       []models.Note
	Accounts    []models.Account
	Blocks      []models.Block
	BlockTypes  []models.BlockType
	BlockStates []models.BlockState
}

func NewMockData() *MockData {
	mock := &MockData{}
	mock.init()
	return mock
}

func (m *MockData) init() {
	now := time.Now()

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	m.Accounts = append(m.Accounts, models.Account{
		ID:           userID,
		Username:     "testuser",
		Password:     []byte("hashed_password"),
		TokenVersion: 1,
		CreatedAt:    now,
		UpdatedAt:    now,
	})

	m.BlockTypes = []models.BlockType{
		{ID: 1, Name: "text"},
		{ID: 2, Name: "image"},
		{ID: 3, Name: "code"},
		{ID: 4, Name: "quote"},
	}

	for i := 1; i <= 5; i++ {
		noteID := uuid.New()

		note := models.Note{
			ID:        noteID,
			UserID:    userID,
			Title:     "Моя заметка " + strconv.Itoa(i),
			ParentID:  nil,
			CreatedAt: now.AddDate(0, -i, 0),
			UpdatedAt: now,
		}
		m.Notes = append(m.Notes, note)

		blockID := uuid.New()
		block := models.Block{
			ID:          blockID,
			NoteID:      noteID,
			BlockTypeID: 1,
			Position:    0,
			Content:     "Пример" + strconv.Itoa(i) + "\n\nПервый параграф текста\nВторой параграф\nТретий параграф",
			CreatedAt:   now.AddDate(0, -i, 0),
			UpdatedAt:   now,
		}
		m.Blocks = append(m.Blocks, block)

		stateID := uuid.New()
		state := models.BlockState{
			ID:         stateID,
			BlockID:    blockID,
			Formatting: `{"format":"text","tags":["пример","тест","мок"]}`,
			CreatedAt:  now.AddDate(0, -i, 0),
			UpdatedAt:  now,
		}
		m.BlockStates = append(m.BlockStates, state)
	}
}

func (m *MockData) GetBlocksByNoteID(noteID uuid.UUID) []models.Block {
	var blocks []models.Block
	for _, block := range m.Blocks {
		if block.NoteID == noteID {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func (m *MockData) GetBlockStatesByBlockID(blockID uuid.UUID) []models.BlockState {
	var states []models.BlockState
	for _, state := range m.BlockStates {
		if state.BlockID == blockID {
			states = append(states, state)
		}
	}
	return states
}
