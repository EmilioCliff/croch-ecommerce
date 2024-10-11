package handlers

import (
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type reviewResponse struct {
	user    userResponse
	review  repository.Review
	product repository.Product
}

type listResponse struct {
	review  repository.Review
	product repository.Product
}

func (s *HttpServer) createReview(ctx *gin.Context) {}

// func (s *HttpServer) listReviews(ctx *gin.Context) {
// 	reviews, err := s.repo.r.ListReviews(ctx)
// 	if err != nil {
// 		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

// 		return
// 	}

// 	var result []reviewResponse
// 	for _, review := range reviews {
// 		product, err := s.repo.p.GetProductById(ctx, review.ProductID)
// 		if err != nil {
// 			ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

// 			return
// 		}

// 		result = append(result, reviewResponse{
// 			pr
// 		})
// 	}

// }

func (s *HttpServer) listProductsReviews(ctx *gin.Context) {}

func (s *HttpServer) ListUsersReviews(ctx *gin.Context) {}

func (s *HttpServer) getReview(ctx *gin.Context) {}

func (s *HttpServer) deleteReview(ctx *gin.Context) {}
