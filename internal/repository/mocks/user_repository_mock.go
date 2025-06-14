package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/timuraipov/diploma/internal/domain"
)

type MockUserRepository struct {
	mock.Mock
	db map[string]domain.User
}

func NewMockUserRepository() MockUserRepository {
	db := make(map[string]domain.User)
	return MockUserRepository{db: db}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.db[user.Login] = *user
	return nil
}

func (m *MockUserRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	if user, ok := m.db[login]; ok {
		return user, nil
	}
	// If user not found, return an error

	return domain.User{}, domain.ErrUserNotFound
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	for _, user := range m.db {
		if user.ID == id {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrUserNotFound
}
