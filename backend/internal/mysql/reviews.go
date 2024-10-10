package mysql

import (
	"context"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

var _ repository.ReviewRepository = (*ReviewsRepository)(nil)

type ReviewsRepository struct {
	db      *Store
	queries generated.Querier
}

func NewReviewRepository(db *Store) *ReviewsRepository {
	q := generated.New(db.db)

	return &ReviewsRepository{
		db:      db,
		queries: q,
	}
}

func (r *ReviewsRepository) CreateReview(ctx context.Context, review *repository.Review) (*repository.Review, error) {
	if err := review.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	result, err := r.queries.CreateReview(ctx, generated.CreateReviewParams{
		ProductID: review.ProductID.String(),
		UserID:    review.UserID.String(),
		Rating:    review.Rating,
		Review:    review.Review,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last inserted id: %v", err)
	}

	review.ID = uint32(id)

	// update the products new rating
	err = r.queries.UpdateRating(ctx, review.ProductID.String())
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update rating: %v", err)
	}

	return review, nil
}

func (r *ReviewsRepository) GetReview(ctx context.Context, id uint32) (*repository.Review, error) {
	review, err := r.queries.GetReview(ctx, id)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	return &repository.Review{
		ID:        review.ID,
		ProductID: uuid.MustParse(review.ProductID),
		UserID:    uuid.MustParse(review.UserID),
		Rating:    review.Rating,
		Review:    review.Review,
	}, nil
}

func (r *ReviewsRepository) ListReviews(ctx context.Context) ([]*repository.Review, error) {
	reviews, err := r.queries.ListReviews(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	var result []*repository.Review
	for _, review := range reviews {
		result = append(result, &repository.Review{
			ID:        review.ID,
			ProductID: uuid.MustParse(review.ProductID),
			UserID:    uuid.MustParse(review.UserID),
			Rating:    review.Rating,
			Review:    review.Review,
		})
	}

	return result, nil
}

func (r *ReviewsRepository) ListUsersReviews(ctx context.Context, userID uuid.UUID) ([]*repository.Review, error) {
	reviews, err := r.queries.ListUsersReviews(ctx, userID.String())
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	var result []*repository.Review
	for _, review := range reviews {
		result = append(result, &repository.Review{
			ID:        review.ID,
			ProductID: uuid.MustParse(review.ProductID),
			UserID:    uuid.MustParse(review.UserID),
			Rating:    review.Rating,
			Review:    review.Review,
		})
	}

	return result, nil
}

func (r *ReviewsRepository) ListProductsReviews(ctx context.Context, productID uuid.UUID) ([]*repository.Review, error) {
	reviews, err := r.queries.ListProductsReviews(ctx, productID.String())
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	var result []*repository.Review
	for _, review := range reviews {
		result = append(result, &repository.Review{
			ID:        review.ID,
			ProductID: uuid.MustParse(review.ProductID),
			UserID:    uuid.MustParse(review.UserID),
			Rating:    review.Rating,
			Review:    review.Review,
		})
	}

	return result, nil
}

func (r *ReviewsRepository) DeleteReview(ctx context.Context, id uint32) error {
	err := r.queries.DeleteReview(ctx, id)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	return nil
}
