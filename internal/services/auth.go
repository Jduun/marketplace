package services

import (
	"context"
	"errors"
	"time"

	"marketplace/config"
	"marketplace/internal/dto"
	"marketplace/internal/entities"
	"marketplace/internal/repositories"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepository repositories.UserRepository
	cfg            *config.Config
}

func NewAuthServiceImpl(userRepository repositories.UserRepository, cfg *config.Config) AuthService {
	return &AuthServiceImpl{
		userRepository: userRepository,
		cfg:            cfg,
	}
}

func (s *AuthServiceImpl) CreateUser(ctx context.Context, userData *dto.UserCreateRequest) (*dto.UserResponse, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordHashing
	}
	hashedPassword := string(hashedPasswordBytes)
	createdUser := entities.User{
		Login:    userData.Login,
		Password: hashedPassword,
	}
	err = s.userRepository.CreateUser(ctx, &createdUser)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return nil, ErrUserAlreadyExists
		} else if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrUserNotFound
		} else {
			return nil, ErrCannotCreateUser
		}
	}
	return &dto.UserResponse{
		ID:        createdUser.ID,
		Login:     createdUser.Login,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *AuthServiceImpl) LoginUser(ctx context.Context, userData *dto.LoginUserRequest) (string, error) {
	user, err := s.userRepository.GetUserByLogin(ctx, userData.Login)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return "", ErrCannotFindUser
		} else {
			return "", ErrCannotLoginUser
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * time.Duration(s.cfg.TokenTTLMinutes)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", ErrCannotSignToken
	}
	return tokenString, nil
}

func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrCannotFindUser
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Login:     user.Login,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthServiceImpl) ParseToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWTSecret), nil
	})
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse((*claims)["sub"].(string))
		if err != nil {
			return uuid.Nil, ErrInvalidToken
		}
		return userID, nil
	}
	return uuid.Nil, ErrInvalidToken
}
