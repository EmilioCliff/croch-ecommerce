package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db      *Store
	queries generated.Querier
}

func NewUserRepository(db *Store) *UserRepository {
	q := generated.New(db.db)

	return &UserRepository{
		db:      db,
		queries: q,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate uuid: %v", err)
	}

	accessToken, err := u.db.tokenMaker.CreateToken(user.ID, user.Email, u.db.config.TOKEN_DURATION)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create token: %v", err)
	}

	refreshToken, err := u.db.tokenMaker.CreateToken(user.ID, user.Email, u.db.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create token: %v", err)
	}

	user.ID = id
	user.RefreshToken = refreshToken

	if err := user.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	hashPass, err := pkg.GenerateHashPassword(user.Password, u.db.config.PASSWORD_COST)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to hash password: %v", err)
	}

	_, err = u.queries.CreateUser(ctx, generated.CreateUserParams{
		ID:           user.ID.String(),
		Email:        user.Email,
		Password:     hashPass,
		Subscription: user.Subscription,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create user: %v", err)
	}

	// change user refresh token to access_token
	user.RefreshToken = accessToken

	return user, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*repository.User, error) {
	user, err := u.queries.GetUserById(ctx, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found with id %s", id.String())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	userID, _ := uuid.Parse(user.ID)
	updatedBy, _ := uuid.Parse(user.UpdatedBy)

	return &repository.User{
		ID:           userID,
		Email:        user.Email,
		Password:     user.Password,
		Subscription: user.Subscription,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
		UpdatedBy:    updatedBy,
		UpdatedAt:    user.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	user, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found with id %s", email)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	userID, _ := uuid.Parse(user.ID)
	updatedBy, _ := uuid.Parse(user.UpdatedBy)

	return &repository.User{
		ID:           userID,
		Email:        user.Email,
		Password:     user.Password,
		Subscription: user.Subscription,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
		UpdatedBy:    updatedBy,
		UpdatedAt:    user.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (u *UserRepository) GetSubscribedUsers(ctx context.Context) ([]*repository.User, error) {
	users, err := u.queries.GetSubscribedUsers(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get users: %v", err)
	}

	var result []*repository.User

	for _, user := range users {
		userID, _ := uuid.Parse(user.ID)
		updatedBy, _ := uuid.Parse(user.UpdatedBy)

		result = append(result, &repository.User{
			ID:           userID,
			Email:        user.Email,
			Password:     user.Password,
			Subscription: user.Subscription,
			Role:         user.Role,
			RefreshToken: user.RefreshToken,
			UpdatedBy:    updatedBy,
			UpdatedAt:    user.UpdatedAt,
			CreatedAt:    user.CreatedAt,
		})
	}

	return result, nil
}

func (u *UserRepository) ListUsers(ctx context.Context) ([]*repository.User, error) {
	users, err := u.queries.ListUsers(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get users: %v", err)
	}

	var result []*repository.User

	for _, user := range users {
		userID, _ := uuid.Parse(user.ID)
		updatedBy, _ := uuid.Parse(user.UpdatedBy)

		result = append(result, &repository.User{
			ID:           userID,
			Email:        user.Email,
			Password:     user.Password,
			Subscription: user.Subscription,
			Role:         user.Role,
			RefreshToken: user.RefreshToken,
			UpdatedBy:    updatedBy,
			UpdatedAt:    user.UpdatedAt,
			CreatedAt:    user.CreatedAt,
		})
	}

	return result, nil
}

func (u *UserRepository) UpdateUserCredentials(ctx context.Context, id uuid.UUID, password string) error {
	hashPass, err := pkg.GenerateHashPassword(password, u.db.config.PASSWORD_COST)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to hash password: %v", err)
	}

	err = u.queries.UpdateUserCredentials(ctx, generated.UpdateUserCredentialsParams{
		ID:        id.String(),
		Password:  hashPass,
		UpdatedBy: id.String(),
		UpdatedAt: time.Now(),
	})

	return err
}

func (u *UserRepository) UpdateUserSubscriptionStatus(ctx context.Context, id uuid.UUID, status bool) error {
	err := u.queries.UpdateSubscriptionStatus(ctx, generated.UpdateSubscriptionStatusParams{
		ID:           id.String(),
		Subscription: status,
		UpdatedBy:    id.String(),
		UpdatedAt:    time.Now(),
	})

	return err
}

func (u *UserRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID) (string, error) {
	user, err := u.queries.GetUserById(ctx, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return "", pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found with id %v", id)
		}

		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	userID, _ := uuid.Parse(user.ID)

	refreshToken, err := u.db.tokenMaker.CreateToken(userID, user.Email, u.db.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create token: %v", err)
	}

	err = u.queries.UpdateRefreshToken(ctx, generated.UpdateRefreshTokenParams{
		ID:           id.String(),
		RefreshToken: refreshToken,
	})

	return refreshToken, err
}

func (u *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := u.queries.DeleteUser(ctx, id.String())

	return err
}
