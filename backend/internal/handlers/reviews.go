package handlers

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
)

type reviewResponse struct {
	ID          uint32    `json:"id"`
	AuthorEmail string    `json:"author_email"`
	Rating      uint32    `json:"rating"`
	Review      string    `json:"review"`
	CreatedAt   time.Time `json:"created_at"`
	ProductName string    `json:"product_name"`
}

type createReviewRequest struct {
	AuthorID  uint32 `binding:"required" json:"author_id"`
	ProductID uint32 `binding:"required" json:"product_id"`
	Review    string `binding:"required" json:"review"`
	Rating    uint32 `binding:"required" json:"rating"`
}

func (s *HttpServer) createReview(ctx *gin.Context) {
	var req createReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	review, err := s.repo.r.CreateReview(ctx, &repository.Review{
		ProductID: req.ProductID,
		UserID:    req.AuthorID,
		Rating:    req.Rating,
		Review:    req.Review,
	})
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	result, err := s.structureReviewResponse([]*repository.Review{review}, ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *HttpServer) listReviews(ctx *gin.Context) {
	reviews, err := s.repo.r.ListReviews(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	result, err := s.structureReviewResponse(reviews, ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *HttpServer) listProductsReviews(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	reviews, err := s.repo.r.ListProductsReviews(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	result, err := s.structureReviewResponse(reviews, ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *HttpServer) listUsersReviews(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	reviews, err := s.repo.r.ListUsersReviews(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	result, err := s.structureReviewResponse(reviews, ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *HttpServer) getReview(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	review, err := s.repo.r.GetReview(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	result, err := s.structureReviewResponse([]*repository.Review{review}, ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *HttpServer) deleteReview(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.r.DeleteReview(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *HttpServer) structureReviewResponse(reviews []*repository.Review, ctx *gin.Context) ([]reviewResponse, error) {
	var result []reviewResponse

	for _, review := range reviews {
		email, err := s.repo.u.GetUserEmail(ctx, review.UserID)
		if err != nil {
			return nil, err
		}

		productName, err := s.repo.p.GetProductName(ctx, review.ProductID)
		if err != nil {
			return nil, err
		}

		result = append(result, reviewResponse{
			ID:          review.ID,
			AuthorEmail: email,
			Rating:      review.Rating,
			Review:      review.Review,
			CreatedAt:   review.CreatedAt,
			ProductName: productName,
		})
	}

	return result, nil
}
