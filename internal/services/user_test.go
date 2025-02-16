package services

import (
	"context"
	"errors"
	"testing"

	"github.com/derticom/merch-store/internal/models"
	"github.com/derticom/merch-store/internal/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	username := "testuser"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		setup       func()
		username    string
		password    string
		expected    *models.User
		expectedErr error
	}{
		{
			name: "successful registration",
			setup: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, user *models.User) error {
						user.ID = uuid.New()
						user.Password = string(hashedPassword)
						return nil
					})
			},
			username: username,
			password: password,
			expected: &models.User{
				Username: username,
				Password: string(hashedPassword),
				Coins:    initialBalance,
			},
			expectedErr: nil,
		},
		{
			name: "username already exists",
			setup: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(&models.User{Username: username}, nil)
			},
			username:    username,
			password:    password,
			expected:    nil,
			expectedErr: errors.New("username already exists"),
		},
		{
			name: "error creating user",
			setup: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			username:    username,
			password:    password,
			expected:    nil,
			expectedErr: errors.New("create error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			user, err := service.RegisterUser(context.Background(), tt.username, tt.password)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Username, user.Username)
				assert.Equal(t, tt.expected.Coins, user.Coins)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)))
			}
		})
	}
}

func TestService_AuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	username := "testuser"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	tests := []struct {
		name        string
		setup       func()
		username    string
		password    string
		expected    *models.User
		expectedErr error
	}{
		{
			name: "successful authentication",
			setup: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(user, nil)
			},
			username:    username,
			password:    password,
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, nil)
			},
			username:    username,
			password:    password,
			expected:    nil,
			expectedErr: errors.New("user not found"),
		},
		{
			name: "invalid password",
			setup: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(user, nil)
			},
			username:    username,
			password:    "wrongpassword",
			expected:    nil,
			expectedErr: errors.New("invalid password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			user, err := service.AuthenticateUser(context.Background(), tt.username, tt.password)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}

func TestService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	userID := uuid.New()
	user := &models.User{
		ID: userID,
	}

	tests := []struct {
		name        string
		setup       func()
		userID      uuid.UUID
		expected    *models.User
		expectedErr error
	}{
		{
			name: "successful get user by ID",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
			},
			userID:      userID,
			expected:    user,
			expectedErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, nil)
			},
			userID:      userID,
			expected:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			user, err := service.GetUserByID(context.Background(), tt.userID)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}

func TestService_UpdateUserCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_services.NewMockRepository(ctrl)
	service := New(mockRepo)

	userID := uuid.New()
	coins := 500

	tests := []struct {
		name        string
		setup       func()
		userID      uuid.UUID
		coins       int
		expectedErr error
	}{
		{
			name: "successful update user coins",
			setup: func() {
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), userID, coins).Return(nil)
			},
			userID:      userID,
			coins:       coins,
			expectedErr: nil,
		},
		{
			name: "error updating user coins",
			setup: func() {
				mockRepo.EXPECT().UpdateUserCoins(gomock.Any(), userID, coins).Return(errors.New("update error"))
			},
			userID:      userID,
			coins:       coins,
			expectedErr: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.UpdateUserCoins(context.Background(), tt.userID, tt.coins)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
