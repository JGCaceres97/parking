package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/pkg/ulid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	exists := s.repo.ExistsUsername(ctx, user.Username)
	if exists {
		return nil, domain.ErrUsernameAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
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

func (s *service) Update(ctx context.Context, id string, userUpdated *domain.User) (*domain.User, error) {
	if userUpdated.Username == domain.AdminUsername {
		return nil, domain.ErrAdminProtected
	}

	existingUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingName, _ := s.repo.FindByUsername(ctx, userUpdated.Username); existingName != nil && id != existingName.ID {
		return nil, domain.ErrUsernameAlreadyExists
	}

	existingUser.Username = userUpdated.Username
	existingUser.Role = userUpdated.Role
	existingUser.IsActive = userUpdated.IsActive

	if err := s.repo.Update(ctx, existingUser); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("error al actualizar usuario en repo: %w", err)
	}

	existingUser.Password = ""
	return existingUser, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Username == domain.AdminUsername {
		return domain.ErrAdminProtected
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return err
		}

		return fmt.Errorf("error al eliminar usuario en repo: %w", err)
	}

	return nil
}

func (s *service) ToggleActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.Username == domain.AdminUsername {
		return nil, domain.ErrAdminProtected
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

func (s *service) ListAll(ctx context.Context, id string) ([]domain.User, error) {
	return s.repo.ListAll(ctx, id)
}

func (s *service) UpdateUsername(ctx context.Context, id string, newUsername string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.Username == domain.AdminUsername {
		return nil, domain.ErrAdminProtected
	}

	if existingName, _ := s.repo.FindByUsername(ctx, newUsername); existingName != nil && id != existingName.ID {
		return nil, domain.ErrUsernameAlreadyExists
	}

	user.Username = newUsername
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error al actualizar username: %w", err)
	}

	user.Password = ""
	return user, nil
}
