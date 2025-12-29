package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/JGCaceres97/parking/internal/application/user"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/pkg/ulid"
)

type service struct {
	repo          user.Repository
	secretKey     []byte
	tokenDuration time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewService(repo user.Repository, secretKey string, tokenDuration time.Duration) Service {
	return &service{
		repo:          repo,
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

func (s *service) CreateAdmin(ctx context.Context, password string) error {
	exists := s.repo.ExistsUsername(ctx, domain.AdminUsername)
	if exists {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al hashear contraseña de administrador: %w", err)
	}

	admin := &domain.User{
		ID:        ulid.GenerateNewULID(),
		Username:  domain.AdminUsername,
		Password:  string(hashedPassword),
		Role:      domain.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	if err := s.repo.Create(ctx, admin); err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, req LoginInput) (*LoginOutput, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("error del repositorio al buscar usuario: %w", err)
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("error al comparar hash: %w", err)
	}

	tokenStr, expirationTime, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	response := &LoginOutput{
		Role:      user.Role,
		Token:     tokenStr,
		TokenType: "Bearer",
		ExpiresIn: int64(time.Until(expirationTime).Seconds()),
	}

	return response, nil
}

func (s *service) ParseToken(tokenStr string) (userID string, role domain.Role, err error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}

		return s.secretKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", ErrExpiredToken
		}

		return "", "", ErrInvalidToken
	}

	if !token.Valid {
		return "", "", ErrInvalidToken
	}

	return claims.UserID, claims.Role, nil
}

func (s *service) generateToken(userID, role string) (string, time.Time, error) {
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
