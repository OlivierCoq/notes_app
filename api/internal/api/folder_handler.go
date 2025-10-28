package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/middleware"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/store"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/utils"
)

type FolderHandler struct {
	// Dependencies for the FolderHandler can be added here, such as a FolderStore or Logger
	folderStore store.FolderStore
	logger      *log.Logger
}

// Constructor for FolderHandler
func NewFolderHandler(folderStore store.FolderStore, logger *log.Logger) *FolderHandler {
	return &FolderHandler{
		folderStore: folderStore,
		logger:      logger,
	}
}

// FolderHandler methods for handling HTTP requests related to folders can be added here.
// CRUD

func (fh *FolderHandler) HandleCreateFolder(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a folder

	var folder store.Folder

	err := json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		fh.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"}) // 400
		return
	}

	// Ensure current user is the owner of the folder
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() {
		fh.logger.Printf("Unauthorized: anonymous user cannot create folders")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	folder.UserID = currentUser.ID // Set the folder's UserID to the current user's ID

	createdFolder, err := fh.folderStore.CreateFolder(&folder)
	if err != nil {
		fh.logger.Printf("Error creating folder: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create folder"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"folder": createdFolder}) // 201

}

func (fh *FolderHandler) HandleGetFolderByID(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting a folder by ID
	folderId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		fh.logger.Printf("Invalid folder ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid folder ID parameter"}) // 400
		return
	}

	folder, err := fh.folderStore.GetFolderByID(int(folderId))
	if err != nil {
		fh.logger.Printf("Error retrieving folder: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve folder"})
		return
	}
	if folder == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Folder not found"})
		return
	}

	// Ensure current user is the owner of the folder
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || folder.UserID != currentUser.ID {
		fh.logger.Printf("Unauthorized access to folder ID %d by user ID %d", folder.ID, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"folder": folder}) // 200
}

func (fh *FolderHandler) HandleUpdateFolder(w http.ResponseWriter, r *http.Request) {

	// Implementation for updating a folder
	paramsFolderId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		fh.logger.Printf("Invalid folder ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid folder ID parameter"}) // 400
		return
	}

	// Fetch existing folder
	existingFolder, err := fh.folderStore.GetFolderByID(int(paramsFolderId))
	if err != nil {
		fh.logger.Printf("Error retrieving folder: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve folder"})
		return
	}
	if existingFolder == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Folder not found"})
		return
	}

	// Ensure current user is the owner of the folder
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || existingFolder.UserID != currentUser.ID {
		fh.logger.Printf("Unauthorized access to folder ID %d by user ID %d", existingFolder.ID, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	// struct
	var updatedFolderRequest struct {
		Name *string `json:"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedFolderRequest)
	if err != nil {
		fh.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"}) // 400
		return
	}

	// Validation
	if updatedFolderRequest.Name != nil {
		existingFolder.Title = *updatedFolderRequest.Name
	}

	// Save the updated folder
	err = fh.folderStore.UpdateFolder(existingFolder)
	if err != nil {
		fh.logger.Printf("Error updating folder: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update folder"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"folder": existingFolder}) // 200
}

func (fh *FolderHandler) HandleDeleteFolder(w http.ResponseWriter, r *http.Request) {

	folderId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		fh.logger.Printf("Invalid folder ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid folder ID parameter"}) // 400
		return
	}

	// Fetch existing folder
	_, err = fh.folderStore.GetFolderByID(int(folderId))
	if err != nil {
		fh.logger.Printf("Error retrieving folder: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve folder"})
		return
	}

	// Ensure current user is the owner of the folder
	currentUser := middleware.GetUser(r)
	folderOwnerID, err := fh.folderStore.GetFolderOwner(int(folderId))
	if err != nil {
		fh.logger.Printf("Error retrieving folder owner: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve folder owner"})
		return
	}
	if currentUser.IsAnonymous() || folderOwnerID != currentUser.ID {
		fh.logger.Printf("Unauthorized access to folder ID %d by user ID %d", folderId, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	// Delete the folder
	err = fh.folderStore.DeleteFolder(int(folderId))
	if err != nil {
		fh.logger.Printf("Error deleting folder: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete folder"}) // 500
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Folder deleted successfully"}) // 200
}

func (fh *FolderHandler) HandleListFoldersByUserID(w http.ResponseWriter, r *http.Request) {

	// Implementation for listing notes by user ID
	userId, err := utils.ReadIDParam(r, "user_id")
	if err != nil {
		fh.logger.Printf("Invalid user ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID parameter"}) // 400
		return
	}

	// Ensure current user is the same as the requested user ID
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || currentUser.ID != int(userId) {
		fh.logger.Printf("Unauthorized access to folders of user ID %d by user ID %d", userId, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	folders, err := fh.folderStore.ListFoldersByUserID(int(userId))
	if err != nil {
		fh.logger.Printf("Error retrieving folders: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve folders"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"folders": folders}) // 200
}
