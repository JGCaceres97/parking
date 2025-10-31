package services

import (
	"context"
	"fmt"
	"time"

	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/ulid"
	"golang.org/x/crypto/bcrypt"
)

type InitService struct {
	repo ports.InitRepository
}

func NewInitService(repo ports.InitRepository) ports.InitService {
	return &InitService{repo: repo}
}

func (s *InitService) CreateAdmin(ctx context.Context, password string) error {
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

	if err := s.repo.CreateAdmin(ctx, admin); err != nil {
		return err
	}

	return nil
}
