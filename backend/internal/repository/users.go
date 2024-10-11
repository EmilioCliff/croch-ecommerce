package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

type User struct {
	ID           uint32 `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Subscription bool   `json:"subscription"`
	Role         string `json:"role"`
	RefreshToken string `json:"refresh_token"`
	UpdatedBy    uint32 `json:"updated_by"`

	// Timestamps
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "email is required")
	}

	if u.Password == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "password is required")
	}

	if u.Role != "USER" && u.Role != "ADMIN" {
		return pkg.Errorf(pkg.INVALID_ERROR, "invalid user role")
	}

	if u.RefreshToken == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "refresh token cannot be empty")
	}

	return nil
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserById(ctx context.Context, id uint32) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetSubscribedUsers(ctx context.Context) ([]*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	UpdateUserCredentials(ctx context.Context, id uint32, password string) error
	UpdateUserSubscriptionStatus(ctx context.Context, id uint32, status bool) error
	UpdateRefreshToken(ctx context.Context, id uint32) (string, error)
	DeleteUser(ctx context.Context, id uint32) error
}
