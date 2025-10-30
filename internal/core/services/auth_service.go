package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/JGCaceres97/parking/internal/api/dtos"
	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
)

type AuthService struct {
	repo          ports.UserRepository
	secretKey     []byte
	tokenDuration time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(repo ports.UserRepository, secretKey string, tokenDuration time.Duration) ports.AuthService {
	return &AuthService{
		repo:          repo,
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

func (s *AuthService) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, ports.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("error del repositorio al buscar usuario: %w", err)
	}

	if !user.IsActive {
		return nil, ports.ErrUserBlocked
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ports.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("error al comparar hash: %w", err)
	}

	tokenStr, expirationTime, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	response := &dtos.LoginResponse{
		Role:      user.Role,
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresIn: int64(time.Until(expirationTime).Seconds()),
	}

	return response, nil
}

func (s *AuthService) ParseToken(tokenStr string) (userID string, role domain.Role, err error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}

		return s.secretKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", ports.ErrTokenExpired
		}

		return "", "", ports.ErrInvalidToken
	}

	if !token.Valid {
		return "", "", ports.ErrInvalidToken
	}

	return claims.UserID, claims.Role, nil
}

func (s *AuthService) generateToken(userID, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(s.tokenDuration)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.secretKey)

	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, expirationTime, nil
}
