package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
)

type userResponse struct {
	ID           uint32 `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Subscription bool   `json:"subscription"`
}

type createUserRequest struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

type createUserResponse struct {
	ID                      uint32 `json:"id"`
	AccessToken             string `json:"access_token"`
	AccessTokenExpiresAfter int64  `json:"access_token_expires_after"`
}

func (s *HttpServer) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	user, err := s.repo.u.CreateUser(ctx, &repository.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     "USER",
	})
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, createUserResponse{
		ID:                      user.ID,
		AccessToken:             user.RefreshToken,
		AccessTokenExpiresAfter: int64(s.config.TOKEN_DURATION.Seconds()),
	})
}

type updateUserRoleRequest struct {
	Role string `binding:"required,oneof=USER ADMIN" json:"role"`
}

func (s *HttpServer) updateUserRole(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	userId, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req updateUserRoleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	if err := s.repo.u.UpdateUserRole(ctx, payload.UserID, userId, req.Role); err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *HttpServer) getUser(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
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
		ID:           user.ID,
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
	ID                       uint32 `json:"id"`
	AccessToken              string `json:"access_token"`
	RefreshToken             string `json:"refresh_token"`
	AccessTokenExpiresAfter  int64  `json:"access_token_expires_after"`
	RefreshTokenExpiresAfter int64  `json:"refresh_token_expires_after"`
}

func (s *HttpServer) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

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

	refreshToken, err := s.repo.u.UpdateRefreshToken(ctx, user.ID)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	accesstoken, err := s.tokenMaker.CreateToken(user.ID, user.Email, user.Role, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		ID:                       user.ID,
		AccessToken:              accesstoken,
		RefreshToken:             refreshToken,
		AccessTokenExpiresAfter:  int64(s.config.TOKEN_DURATION.Seconds()),
		RefreshTokenExpiresAfter: int64(s.config.REFRESH_TOKEN_DURATION.Seconds()),
	})
}

func (s *HttpServer) refreshToken(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
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
	if _, err := s.tokenMaker.VerifyToken(user.RefreshToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "refresh token expired, kindly login")))

		return
	}

	accesstoken, err := s.tokenMaker.CreateToken(user.ID, user.Email, user.Role, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		ID:                       user.ID,
		AccessToken:              accesstoken,
		RefreshToken:             user.RefreshToken,
		AccessTokenExpiresAfter:  int64(s.config.TOKEN_DURATION.Seconds()),
		RefreshTokenExpiresAfter: int64(s.config.REFRESH_TOKEN_DURATION.Seconds()),
	})
}

type resetPasswordRequest struct {
	Email string `binding:"required" json:"email"`
}

func (s *HttpServer) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	// send password reset email
	// email will have a token that will enable one to send the data
	ctx.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *HttpServer) listUsers(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if isAdmin, err := isAdmin(payload); !isAdmin {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	users, err := s.repo.u.ListUsers(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	response := []userResponse{}
	for _, user := range users {
		response = append(response, userResponse{
			ID:           user.ID,
			Email:        user.Email,
			Role:         user.Role,
			Subscription: user.Subscription,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

type updateUserSubscriptionRequestBody struct {
	Subscription bool `binding:"required" json:"subscription"`
}

func (s *HttpServer) updateUserSubscription(ctx *gin.Context) {
	payload, err := getPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if id != payload.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "unauthorized to this action")))

		return
	}

	// method put
	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req updateUserSubscriptionRequestBody
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	if err := s.repo.u.UpdateUserSubscriptionStatus(ctx, id, req.Subscription); err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
