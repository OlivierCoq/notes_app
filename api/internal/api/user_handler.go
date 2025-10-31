package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/middleware"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/store"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/utils"
)

// This is for user registration requests
type RegisterUserRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Bio            string `json:"bio"`
	AddressLine1   string `json:"address_line_1"`
	AddressLine2   string `json:"address_line_2"`
	AddressCity    string `json:"address_city"`
	AddressState   string `json:"address_state"`
	AddressZip     string `json:"address_zip"`
	AddressCountry string `json:"address_country"`
	PfpURL         string `json:"pfp_url"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	AuthLevel      int    `json:"auth_level"` // Optional, default to regular user
}

type UserHandler struct {
	// Add fields as necessary, e.g., a reference to the application or database
	userStore store.UserStore // Interface to interact with user data. This promotes db decoupling and easier testing.
	logger    *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

// Validation:
func (h *UserHandler) validateRegisterUserRequest(req *RegisterUserRequest) error {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return errors.New("username, email, and password are required")
	}
	if len(req.Username) > 50 {
		return errors.New("username must be less than 50 characters")
	}
	// Email
	if len(req.Email) > 100 {
		return errors.New("email must be less than 100 characters")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	// Password
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(req.Password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(req.Password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(req.Password)
	if !(hasLower && hasUpper && hasDigit) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

// Define methods for UserHandler to handle user-related requests. CRUD operations, etc.

// Create/Register User
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	// Decode the POST request body into the RegisterUserRequest struct:
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"}) // 400
		return
	}

	// Validate the request
	err = h.validateRegisterUserRequest(&req)
	if err != nil {
		h.logger.Printf("Validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()}) // 400
		return
	}

	// Send to DB via userStore
	user := &store.User{
		Username:       req.Username,
		Email:          req.Email,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		AddressCity:    req.AddressCity,
		AddressState:   req.AddressState,
		AddressZip:     req.AddressZip,
		AddressCountry: req.AddressCountry,
		PfpURL:         req.PfpURL,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Bio:            req.Bio,
		AuthLevel:      req.AuthLevel, // Default to 1 (regular user) if not provided
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"}) // 500
		return
	}

	createdUser, err := h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create user"}) // 500
		return
	}

	// Respond with the created user (excluding password hash) as JSON to the frontend:
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": createdUser}) // 201

}

// Additional user-related handlers can be added here (e.g., GetUser, UpdateUser, DeleteUser, etc.)
func (h *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting a user by ID
	userId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		h.logger.Printf("Invalid user ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID parameter"}) // 400
		return
	}

	user, err := h.userStore.GetUserById(int(userId))
	if err != nil {
		h.logger.Printf("Error fetching user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}) // 200
}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a user by ID
	userId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		h.logger.Printf("Invalid user ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID parameter"}) // 400
		return
	}

	var req RegisterUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	// Ensure current user is the one being updated
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || currentUser.ID != int(userId) {
		h.logger.Printf("Unauthorized update attempt for user ID %d by user ID %d", userId, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	// // Validate the request
	// err = h.validateRegisterUserRequest(&req)
	// if err != nil {
	// 	h.logger.Printf("Validation error: %v", err)
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()}) // 400
	// 	return
	// }

	user := &store.User{
		ID:             int(userId),
		Username:       req.Username,
		Email:          req.Email,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		AddressCity:    req.AddressCity,
		AddressState:   req.AddressState,
		AddressZip:     req.AddressZip,
		AddressCountry: req.AddressCountry,
		PfpURL:         req.PfpURL,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Bio:            req.Bio,
		AuthLevel:      req.AuthLevel,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"}) // 500
		return
	}

	updatedUser, err := h.userStore.UpdateUser(user)
	if err != nil {
		h.logger.Printf("Error updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": updatedUser}) // 200
}

func (h *UserHandler) HandleGetSelf(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting the currently authenticated user
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() {
		h.logger.Printf("Unauthorized access to self user data")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}

	user, err := h.userStore.GetUserById(currentUser.ID)
	if err != nil {
		h.logger.Printf("Error fetching user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}) // 200
}

// Update user password:
func (h *UserHandler) HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a user's password
	userId, err := utils.ReadIDParam(r, "id")
	if err != nil {
		h.logger.Printf("Invalid user ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID parameter"}) // 400
		return
	}

	var req struct {
		NewPassword string `json:"new_password"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Invalid request payload: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	// Ensure current user is the one being updated
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() || currentUser.ID != int(userId) {
		h.logger.Printf("Unauthorized password update attempt for user ID %d by user ID %d", userId, currentUser.ID)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"}) // 401
		return
	}
	err = h.userStore.UpdateUserPassword(int(userId), req.NewPassword)
	if err != nil {
		h.logger.Printf("Error updating user password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update user password"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Password updated successfully"}) // 200
}
