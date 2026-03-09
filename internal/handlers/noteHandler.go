package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/mock"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/models"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/helpers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"
	"github.com/google/uuid"
)

var (
	ErrNoteIDRequired = errors.New("note id is required")
	ErrInvalidNoteID  = errors.New("invalid note id")
	ErrNoteNotFound   = errors.New("note not found")
	ErrInvalidPath    = errors.New("invalid path")
)

const (
	notesKey   = "notes"
	totalKey   = "total"
	noteKey    = "note"
	blocksKey  = "blocks"
	noteIDKey  = "note_id"

	blocksPath = "blocks"
	
	minNotePathLength = 3
	minBlocksPathLength = 4
)

type NoteHandler struct {
	mockData *mock.MockData
}

func NewNoteHandler(mockData *mock.MockData) *NoteHandler {
	return &NoteHandler{
		mockData: mockData,
	}
}

func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserIDKey).(string)
	if !ok {
		helpers.JSONErrorResponse(w, http.StatusUnauthorized, jwt.ErrNoUserID)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidNoteID)
		return
	}

	var userNotes []models.Note
	for _, note := range h.mockData.Notes {
		if note.UserID == userUUID {
			userNotes = append(userNotes, note)
		}
	}

	if len(userNotes) == 0 {
		userNotes = h.mockData.Notes
	}

	response := map[string]interface{}{
		notesKey: userNotes,
		totalKey: len(userNotes),
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < minNotePathLength {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrNoteIDRequired)
		return
	}

	noteIDStr := pathParts[2]
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidNoteID)
		return
	}

	var foundNote *models.Note
	for _, note := range h.mockData.Notes {
		if note.ID == noteID {
			foundNote = &note
			break
		}
	}

	if foundNote == nil {
		helpers.JSONErrorResponse(w, http.StatusNotFound, ErrNoteNotFound)
		return
	}

	var noteBlocks []models.Block
	for _, blockID := range foundNote.Blocks {
		for _, block := range h.mockData.Blocks {
			if block.ID == blockID {
				noteBlocks = append(noteBlocks, block)
				break
			}
		}
	}

	response := map[string]interface{}{
		noteKey:   foundNote,
		blocksKey: noteBlocks,
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}

func (h *NoteHandler) GetNoteBlocks(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < minBlocksPathLength || pathParts[3] != blocksPath {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidPath)
		return
	}

	noteIDStr := pathParts[2]
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		helpers.JSONErrorResponse(w, http.StatusBadRequest, ErrInvalidNoteID)
		return
	}

	var foundNote *models.Note
	for _, note := range h.mockData.Notes {
		if note.ID == noteID {
			foundNote = &note
			break
		}
	}

	if foundNote == nil {
		helpers.JSONErrorResponse(w, http.StatusNotFound, ErrNoteNotFound)
		return
	}

	var noteBlocks []models.Block
	for _, blockID := range foundNote.Blocks {
		for _, block := range h.mockData.Blocks {
			if block.ID == blockID {
				noteBlocks = append(noteBlocks, block)
				break
			}
		}
	}

	response := map[string]interface{}{
		blocksKey:  noteBlocks,
		noteIDKey: foundNote.ID,
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}
