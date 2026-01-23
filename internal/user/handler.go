package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo UserRepository
}

func NewHandler(repo UserRepository) *Handler {
	return &Handler{repo: repo}
}

// Request/Response DTOs
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// POST /api/users/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Username == "" {
		h.writeError(w, http.StatusBadRequest, "email, password and username are required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	// Create user
	now := time.Now()
	user := &User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Username:  req.Username,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.repo.Create(r.Context(), user); err != nil {
		if errors.Is(err, ErrEmailExists) {
			h.writeError(w, http.StatusConflict, "email already exists")
			return
		}
		if errors.Is(err, ErrUsernameExists) {
			h.writeError(w, http.StatusConflict, "username already exists")
			return
		}
		log.Println("Create user error:", err)
		h.writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	h.writeJSON(w, http.StatusCreated, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

// GET /api/users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse ID from URL path
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	h.writeJSON(w, http.StatusOK, UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{Error: message})
}
