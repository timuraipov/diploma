package repository

import (
	"context"
	"os"
	"testing"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/storage/db"
)

var testRepository *userRepository

func setupTestDBUser(t *testing.T) {

	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		t.Fatal("DATABASE_URI is not set")
	}

	db, err := db.NewDB(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_, err = db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS "user" (
    id       INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login    VARCHAR(200) NOT NULL UNIQUE,
    password VARCHAR(200) NOT NULL
);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	testRepository = NewUserRepository(*db)
}

func cleanupTestDBUser(t *testing.T) {
	_, err := testRepository.database.Pool.Exec(context.Background(), `DROP TABLE "user"`)
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}
}

func TestUserRepository_CreateAndGetUserByLogin(t *testing.T) {
	setupTestDBUser(t)
	defer cleanupTestDBUser(t)

	testCases := []struct {
		name     string
		user     domain.User
		expected domain.User
	}{
		{
			name: "Create and get user",
			user: domain.User{
				Login:    "testuser",
				Password: "testpassword",
			},
			expected: domain.User{
				Login:    "testuser",
				Password: "testpassword",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testRepository.Create(context.Background(), &tc.user)
			if err != nil {
				t.Fatalf("failed to create user: %v", err)
			}

			gotUser, err := testRepository.GetByLogin(context.Background(), tc.user.Login)
			if err != nil {
				t.Fatalf("failed to get user by login: %v", err)
			}

			if gotUser.Login != tc.expected.Login || gotUser.Password != tc.expected.Password {
				t.Errorf("expected user %v, got %v", tc.expected, gotUser)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	setupTestDBUser(t)
	defer cleanupTestDBUser(t)

	// Create a test user
	user := domain.User{
		Login:    "testuser2",
		Password: "testpassword",
	}
	err := testRepository.Create(context.Background(), &user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Fetch the user by ID
	gotUser, err := testRepository.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("failed to get user by ID: %v", err)
	}

	// Validate the result
	if gotUser.ID != user.ID || gotUser.Login != user.Login || gotUser.Password != user.Password {
		t.Errorf("expected user %v, got %v", user, gotUser)
	}
}
