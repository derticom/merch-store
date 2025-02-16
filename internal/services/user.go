package services

import (
	"context"
	"errors"

	"github.com/derticom/merch-store/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	RegisterUser(ctx context.Context, username, password string) (*models.User, error)
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUserCoins(ctx context.Context, id uuid.UUID, coins int) error
}

func (s *Service) RegisterUser(ctx context.Context, username, password string) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким именем
	existingUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем нового пользователя
	user := &models.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
		Coins:    1000, // Начальный баланс
	}

	// Сохраняем пользователя в базе данных
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	// Получаем пользователя по имени
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) UpdateUserCoins(ctx context.Context, id uuid.UUID, coins int) error {
	return s.repo.UpdateUserCoins(ctx, id, coins)
}
