package handlers

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Subscription bool   `json:"subscription"`
}

type createUserRequest struct {
	Email    string `binding:"required"                  json:"email"`
	Password string `binding:"required"                  json:"password"`
	Role     string `binding:"required,oneof=USER ADMIN" json:"role"`
}

type createUserResponse struct {
	ID                      string        `json:"id"`
	AccessToken             string        `json:"access_token"`
	AccessTokenExpiresAfter time.Duration `json:"access_token_expires_after"`
}

func (s *HttpServer) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	user, err := s.repo.u.CreateUser(ctx, &repository.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, createUserResponse{
		ID:                      user.ID.String(),
		AccessToken:             user.RefreshToken,
		AccessTokenExpiresAfter: s.config.TOKEN_DURATION,
	})
}

type getUserRequest struct {
	ID string `binding:"required" uri:"id"`
}

func (s *HttpServer) getUser(ctx *gin.Context) {
	var uri getUserRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	user, err := s.repo.u.GetUserById(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, userResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
	})
}

type loginUserRequest struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

type loginUserResponse struct {
	ID                       string        `json:"id"`
	AccessToken              string        `json:"access_token"`
	RefreshToken             string        `json:"refresh_token"`
	AccessTokenExpiresAfter  time.Duration `json:"access_token_expires_after"`
	RefreshTokenExpiresAfter time.Duration `json:"refresh_token_expires_after"`
}

func (s *HttpServer) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	user, err := s.repo.u.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	if err := pkg.ComparePasswordAndHash(user.Password, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	// add new refresh token in the db

	accesstoken, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		ID:                       user.ID.String(),
		AccessToken:              accesstoken,
		RefreshToken:             user.RefreshToken,
		AccessTokenExpiresAfter:  s.config.TOKEN_DURATION,
		RefreshTokenExpiresAfter: s.config.REFRESH_TOKEN_DURATION,
	})
}

type refreshTokenRequest struct {
	ID string `binding:"required" uri:"id"`
}

func (s *HttpServer) refreshToken(ctx *gin.Context) {
	var uri refreshTokenRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	user, err := s.repo.u.GetUserById(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	// check if the refresh token has expired

	// if expired ask user to login so as to create a new refresh token

	// if not expired, create new access token

	accesstoken, err := s.tokenMaker.CreateToken(user.ID, user.Email, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	refreshToken, err := s.repo.u.UpdateRefreshToken(ctx, user.ID)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		ID:                       user.ID.String(),
		AccessToken:              accesstoken,
		RefreshToken:             refreshToken,
		AccessTokenExpiresAfter:  s.config.TOKEN_DURATION,
		RefreshTokenExpiresAfter: s.config.REFRESH_TOKEN_DURATION,
	})
}

type resetPasswordRequest struct {
	Email string `binding:"required" json:"email"`
}

func (s *HttpServer) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	// send password reset email
	ctx.JSON(http.StatusNotImplemented, gin.H{})
}

func (s HttpServer) listUsers(ctx *gin.Context) {
	users, err := s.repo.u.ListUsers(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	var response []userResponse
	for _, user := range users {
		response = append(response, userResponse{
			ID:           user.ID.String(),
			Email:        user.Email,
			Role:         user.Role,
			Subscription: user.Subscription,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

type updateUserSubscriptionRequest struct {
	ID string `binding:"required" uri:"id"`

	Subscription bool `binding:"required" json:"subscription"`
}

func (s *HttpServer) updateUserSubscription(ctx *gin.Context) {
	var uri updateUserSubscriptionRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var req updateUserSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.u.UpdateUserSubscriptionStatus(ctx, id, req.Subscription)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "OK"})
}
