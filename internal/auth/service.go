package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"spotify-clone/internal/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthService defines the authentication business logic interface
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error)
}

type authService struct {
	userRepo   user.UserRepository
	jwtService JWTService
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo user.UserRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register creates a new user account and returns tokens
func (s *authService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user entity
	now := time.Now()
	newUser := &user.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(newUser.ID.String(), newUser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(newUser.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:        newUser.ID.String(),
			Email:     newUser.Email,
			Username:  newUser.Username,
			CreatedAt: newUser.CreatedAt,
		},
	}, nil
}

// Login authenticates user and returns tokens
func (s *authService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Find user by username
	foundUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(foundUser.ID.String(), foundUser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(foundUser.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:        foundUser.ID.String(),
			Email:     foundUser.Email,
			Username:  foundUser.Username,
			CreatedAt: foundUser.CreatedAt,
		},
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *authService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Find user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	foundUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(foundUser.ID.String(), foundUser.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(foundUser.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:        foundUser.ID.String(),
			Email:     foundUser.Email,
			Username:  foundUser.Username,
			CreatedAt: foundUser.CreatedAt,
		},
	}, nil
}
