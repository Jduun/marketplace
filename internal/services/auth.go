package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"marketplace/internal/dto"
	"marketplace/internal/entities"
	slogger "marketplace/internal/logger"
	"marketplace/internal/repositories"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepository  repositories.UserRepository
	tokenTTLMinutes time.Duration
	JWTSecret       string
}

func NewAuthServiceImpl(userRepository repositories.UserRepository, tokenTTLMinutes time.Duration, JWTSecret string) AuthService {
	return &AuthServiceImpl{
		userRepository:  userRepository,
		tokenTTLMinutes: tokenTTLMinutes,
		JWTSecret:       JWTSecret,
	}
}

func (s *AuthServiceImpl) CreateUser(ctx context.Context, userData *dto.UserCreateRequest) (*dto.UserResponse, error) {
	logger := slogger.GetLoggerFromContext(ctx).
		With(slog.String("op", "services.advertisement.CreateAdvertisement"))

	logger.Info("Creating user", slog.String("login", userData.Login))

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Password hashing failed", slog.Any("error", err))
		return nil, ErrPasswordHashing
	}
	hashedPassword := string(hashedPasswordBytes)
	createdUser := entities.User{
		Login:    userData.Login,
		Password: hashedPassword,
	}
	err = s.userRepository.CreateUser(ctx, &createdUser)
	if err != nil {
		logger.Error("User creation failed", slog.Any("error", err))
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return nil, ErrUserAlreadyExists
		} else if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrUserNotFound
		} else {
			return nil, ErrCannotCreateUser
		}
	}

	logger.Info("User created successfully", slog.String("userID", createdUser.ID.String()))

	return &dto.UserResponse{
		ID:        createdUser.ID,
		Login:     createdUser.Login,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *AuthServiceImpl) LoginUser(ctx context.Context, userData *dto.LoginUserRequest) (string, error) {
	logger := slogger.GetLoggerFromContext(ctx).
		With(slog.String("op", "services.auth.LoginUser"))

	logger.Info("Attempting login", slog.String("login", userData.Login))

	user, err := s.userRepository.GetUserByLogin(ctx, userData.Login)
	if err != nil {
		logger.Error("User lookup failed", slog.Any("error", err))
		if errors.Is(err, repositories.ErrNotFound) {
			return "", ErrCannotFindUser
		} else {
			return "", ErrCannotLoginUser
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		logger.Warn("Invalid credentials", slog.String("login", userData.Login))
		return "", ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"login": user.Login,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * s.tokenTTLMinutes).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		logger.Error("Token signing failed", slog.Any("error", err))
		return "", ErrCannotSignToken
	}

	logger.Info("User logged in successfully", slog.String("userID", user.ID.String()))

	return tokenString, nil
}

func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	logger := slogger.GetLoggerFromContext(ctx).
		With(slog.String("op", "services.auth.GetUserByID"))

	logger.Info("Fetching user by ID", slog.String("userID", id.String()))

	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		logger.Error("User fetch failed", slog.Any("error", err))
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrCannotFindUser
	}

	logger.Info("User fetched successfully", slog.String("userID", user.ID.String()))

	return &dto.UserResponse{
		ID:        user.ID,
		Login:     user.Login,
		CreatedAt: user.CreatedAt,
	}, nil
}
