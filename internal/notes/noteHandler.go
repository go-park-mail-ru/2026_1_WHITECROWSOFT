package notes

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
		"notes": userNotes,
		"total": len(userNotes),
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
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

	blocks := h.mockData.GetBlocksByNoteID(foundNote.ID)

	blocksWithStates := make([]map[string]interface{}, 0, len(blocks))
	for _, block := range blocks {
		states := h.mockData.GetBlockStatesByBlockID(block.ID)

		blockData := map[string]interface{}{
			"id":       block.ID,
			"note_id":  block.NoteID,
			"type_id":  block.BlockTypeID,
			"position": block.Position,
			"content":  block.Content,
			"states":   states,
		}
		blocksWithStates = append(blocksWithStates, blockData)
	}

	response := map[string]interface{}{
		"note":   foundNote,
		"blocks": blocksWithStates,
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}
