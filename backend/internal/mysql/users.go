package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/go-sql-driver/mysql"
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
	accessToken, err := u.db.tokenMaker.CreateToken(user.ID, user.Email, user.Role, u.db.config.TOKEN_DURATION)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create token: %v", err)
	}

	refreshToken, err := u.db.tokenMaker.CreateToken(user.ID, user.Email, user.Role, u.db.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create token: %v", err)
	}

	user.RefreshToken = refreshToken

	if err := user.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	hashPass, err := pkg.GenerateHashPassword(user.Password, u.db.config.PASSWORD_COST)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to hash password: %v", err)
	}

	result, err := u.queries.CreateUser(ctx, generated.CreateUserParams{
		Email:        user.Email,
		Password:     hashPass,
		Subscription: user.Subscription,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
	})
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "duplicate entry for email: %s", user.Email)
			}
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create user: %v", err)
	}

	createdId, err := result.LastInsertId()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last inserted id: %v", err)
	}

	user.ID = uint32(createdId)

	// change user refresh token to access_token
	user.RefreshToken = accessToken

	return user, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, id uint32) (*repository.User, error) {
	user, err := u.queries.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found with id %d", id)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	return &repository.User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		Subscription: user.Subscription,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
		UpdatedBy:    uint32(user.UpdatedBy.Int32),
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

	return &repository.User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		Subscription: user.Subscription,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
		UpdatedBy:    uint32(user.UpdatedBy.Int32),
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
		result = append(result, &repository.User{
			ID:           user.ID,
			Email:        user.Email,
			Password:     user.Password,
			Subscription: user.Subscription,
			Role:         user.Role,
			RefreshToken: user.RefreshToken,
			UpdatedBy:    uint32(user.UpdatedBy.Int32),
			UpdatedAt:    user.UpdatedAt,
			CreatedAt:    user.CreatedAt,
		})
	}

	return result, nil
}

func (u *UserRepository) GetUserEmail(ctx context.Context, id uint32) (string, error) {
	email, err := u.queries.GetUserEmail(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found with id %d", id)
		}

		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	return email, nil
}

func (u *UserRepository) ListUsers(ctx context.Context) ([]*repository.User, error) {
	users, err := u.queries.ListUsers(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get users: %v", err)
	}

	var result []*repository.User

	for _, user := range users {
		result = append(result, &repository.User{
			ID:           user.ID,
			Email:        user.Email,
			Password:     user.Password,
			Subscription: user.Subscription,
			Role:         user.Role,
			RefreshToken: user.RefreshToken,
			UpdatedBy:    uint32(user.UpdatedBy.Int32),
			UpdatedAt:    user.UpdatedAt,
			CreatedAt:    user.CreatedAt,
		})
	}

	return result, nil
}

func (u *UserRepository) UpdateUserCredentials(ctx context.Context, id uint32, password string) error {
	hashPass, err := pkg.GenerateHashPassword(password, u.db.config.PASSWORD_COST)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to hash password: %v", err)
	}

	if id <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "invalid user id")
	}

	err = u.queries.UpdateUserCredentials(ctx, generated.UpdateUserCredentialsParams{
		ID:       id,
		Password: hashPass,
		UpdatedBy: sql.NullInt32{
			Valid: true,
			Int32: int32(id),
		},
	})

	return err
}

func (u *UserRepository) UpdateUserSubscriptionStatus(ctx context.Context, id uint32, status bool) error {
	if id <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "invalid user id")
	}

	err := u.queries.UpdateSubscriptionStatus(ctx, generated.UpdateSubscriptionStatusParams{
		ID:           id,
		Subscription: status,
		UpdatedBy: sql.NullInt32{
			Valid: true,
			Int32: int32(id),
		},
	})

	return err
}

func (u *UserRepository) UpdateUserRole(ctx context.Context, adminId uint32, userId uint32, role string) error {
	if role != "ADMIN" && role != "USER" {
		return pkg.Errorf(pkg.INVALID_ERROR, "invalid user role")
	}

	err := u.queries.UpdateUserRole(ctx, generated.UpdateUserRoleParams{
		Role: role,
		UpdatedBy: sql.NullInt32{
			Valid: true,
			Int32: int32(adminId),
		},
		ID: userId,
	})

	return err
}

func (u *UserRepository) UpdateRefreshToken(ctx context.Context, id uint32) (string, error) {
	user, err := u.queries.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found with id %v", id)
		}

		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	refreshToken, err := u.db.tokenMaker.CreateToken(user.ID, user.Email, user.Role, u.db.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create token: %v", err)
	}

	err = u.queries.UpdateRefreshToken(ctx, generated.UpdateRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: refreshToken,
	})

	return refreshToken, err
}

func (u *UserRepository) DeleteUser(ctx context.Context, id uint32) error {
	err := u.queries.DeleteUser(ctx, id)

	return err
}
