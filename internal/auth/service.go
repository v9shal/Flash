package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, FullName, email, password string, phone *string) (*User, error)
}
type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) RegisterUser(
	ctx context.Context,
	fullName string,
	email string,
	password string,
	phone *string,
) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password %w", err)
	}
	user := &User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
	}
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	// 4. Return created user
	return user, nil

}
