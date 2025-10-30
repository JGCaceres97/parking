package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/pkg/ulid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := s.repo.FindByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, ports.ErrUserNotFound) {
		return nil, fmt.Errorf("error al buscar usuario por nombre: %w", err)
	}

	if existingUser != nil {
		return nil, ports.ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al hashear password: %w", err)
	}

	user.ID = ulid.GenerateNewULID()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC().Truncate(time.Second)

	if err = s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error al guardar el usuario: %w", err)
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id string, userUpdated *domain.User) (*domain.User, error) {
	existingUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingName, _ := s.repo.FindByUsername(ctx, userUpdated.Username); existingName != nil && id != existingName.ID {
		return nil, ports.ErrUsernameExists
	}

	existingUser.Username = userUpdated.Username
	existingUser.Role = userUpdated.Role
	existingUser.IsActive = userUpdated.IsActive

	if err := s.repo.Update(ctx, existingUser); err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("error al actualizar usuario en repo: %w", err)
	}

	existingUser.Password = ""
	return existingUser, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return err
		}

		return fmt.Errorf("error al eliminar usuario en repo: %w", err)
	}

	return nil
}

func (s *UserService) ToggleActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.IsActive == isActive {
		return user, nil
	}

	user.IsActive = isActive
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error al cambiar estado activo del usuario: %w", err)
	}

	return user, nil
}

func (s *UserService) ListAll(ctx context.Context, id string) ([]domain.User, error) {
	return s.repo.ListAll(ctx, id)
}

func (s *UserService) UpdateUsername(ctx context.Context, id string, newUsername string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Username = newUsername
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error al actualizar username: %w", err)
	}

	user.Password = ""
	return user, nil
}
