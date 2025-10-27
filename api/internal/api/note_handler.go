package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/middleware"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/store"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/utils"
)


type NoteHandler struct {
	// Dependencies for the NoteHandler can be added here, such as a NoteStore or Logger
	notesStore store.NoteStore
	logger     *log.Logger
}


// Constructor for NoteHandler
func NewNoteHandler(notesStore store.NoteStore, logger *log.Logger) *NoteHandler {
	return &NoteHandler{
		notesStore: notesStore,
		logger:     logger,
	}
}


// NoteHandler methods for handling HTTP requests related to notes can be added here.
// CRUD

func (nh *NoteHandler) HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a note

	var note store.Note

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		nh.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"}) // 400
		return
	}

	// Ensure current user is the owner of the note
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() {
		nh.logger.Printf("Unauthorized: anonymous user cannot create notes")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	note.UserID = currentUser.ID // Set the note's UserID to the current user's ID

	createdNote, err := nh.notesStore.CreateNote(&note)
	if err != nil {
		nh.logger.Printf("Error creating note: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create note"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"note": createdNote}) // 201

}

func (nh *NoteHandler) HandleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting a note by ID
	noteId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		nh.logger.Printf("Invalid note ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid note ID parameter"}) // 400
		return
	}
	
	note, err := nh.notesStore.GetNoteByID(int(noteId))
	if err != nil {
		nh.logger.Printf("Error retrieving note: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve note"})
		return
	}
	if note == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Note not found"})
		return
	}

	// Ensure current user is the owner of the note
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || note.UserID != currentUser.ID {
		nh.logger.Printf("Unauthorized access to note ID %d by user ID %d", note.ID, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"note": note}) // 200
}

func (nh *NoteHandler) HandleUpdateNote(w http.ResponseWriter, r *http.Request) {

	// Implementation for updating a note
	paramsNoteId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		nh.logger.Printf("Invalid note ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid note ID parameter"}) // 400
		return
	}

	// Fetch existing note
	existingNote, err := nh.notesStore.GetNoteByID(int(paramsNoteId))
	if err != nil {
		nh.logger.Printf("Error retrieving note: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve note"})
		return
	}
	if existingNote == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Note not found"})
		return
	}

	// Ensure current user is the owner of the note
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || existingNote.UserID != currentUser.ID {
		nh.logger.Printf("Unauthorized access to note ID %d by user ID %d", existingNote.ID, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	// struct
	var updatedNoteRequest struct {
		Title      *string `json:"title"`
		Content    *string `json:"content"`
		IsFavorite *bool    `json:"is_favorite"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedNoteRequest)
	if err != nil {
		nh.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"}) // 400
		return
	}

	// Validation
	if updatedNoteRequest.Title != nil {
		existingNote.Title = *updatedNoteRequest.Title
	}
	if updatedNoteRequest.Content != nil {
		existingNote.Content = *updatedNoteRequest.Content
	}
	if updatedNoteRequest.IsFavorite != nil {
		existingNote.IsFavorite = *updatedNoteRequest.IsFavorite
	}

	// Save the updated note
	err = nh.notesStore.UpdateNote(existingNote)
	if err != nil {
		nh.logger.Printf("Error updating note: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update note"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"note": existingNote}) // 200
}

func (nh *NoteHandler) HandleDeleteNote(w http.ResponseWriter, r *http.Request) {

	noteId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		nh.logger.Printf("Invalid note ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid note ID parameter"}) // 400
		return
	}
	
	// Fetch existing note
	_, err = nh.notesStore.GetNoteByID(int(noteId))
	if err != nil {
		nh.logger.Printf("Error retrieving note: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve note"})
		return
	}

	// Ensure current user is the owner of the note
	currentUser := middleware.GetUser(r)
	noteOwnerID, err := nh.notesStore.GetNoteOwner(int(noteId))
	if err != nil {
		nh.logger.Printf("Error retrieving note owner: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve note owner"})
		return
	}
	if currentUser.IsAnonymous() || noteOwnerID != currentUser.ID {
		nh.logger.Printf("Unauthorized access to note ID %d by user ID %d", noteId, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	// Delete the note
	err = nh.notesStore.DeleteNote(int(noteId))
	if err != nil {
		nh.logger.Printf("Error deleting note: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete note"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Note deleted successfully"}) // 200
}