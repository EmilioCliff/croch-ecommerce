package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

var _ repository.BlogRepository = (*BlogRepository)(nil)

type BlogRepository struct {
	db      *Store
	queries generated.Querier
}

func NewBlogRepository(db *Store) *BlogRepository {
	q := generated.New(db.db)

	return &BlogRepository{
		db:      db,
		queries: q,
	}
}

func (b *BlogRepository) CreateBlog(ctx context.Context, blog *repository.Blog) (*repository.Blog, error) {
	if err := blog.Validate(); err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	result, err := b.queries.CreateBlog(ctx, generated.CreateBlogParams{
		Author:  blog.Author,
		Title:   blog.Title,
		Content: blog.Content,
		ImgUrls: blog.ImgUrls,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create blog: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last inserted id: %v", err)
	}

	blog.ID = uint32(id)

	// send email to subscribed users

	return blog, err
}

func (b *BlogRepository) GetBlog(ctx context.Context, id uint32) (*repository.Blog, error) {
	blog, err := b.queries.GetBlog(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no blog found with id %d", id)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get blog: %v", err)
	}

	return &repository.Blog{
		ID:        blog.ID,
		Author:    blog.Author,
		Title:     blog.Title,
		Content:   blog.Content,
		ImgUrls:   blog.ImgUrls,
		CreatedAt: blog.CreatedAt,
	}, nil
}

func (b *BlogRepository) GetBlogsByAuthor(ctx context.Context, author uint32) ([]*repository.Blog, error) {
	blogs, err := b.queries.GetBlogsByAuthor(ctx, author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no blog found with author %s", author)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get blogs: %v", err)
	}

	var result []*repository.Blog

	for _, blog := range blogs {
		result = append(result, &repository.Blog{
			ID:        blog.ID,
			Author:    blog.Author,
			Title:     blog.Title,
			Content:   blog.Content,
			ImgUrls:   blog.ImgUrls,
			CreatedAt: blog.CreatedAt,
		})
	}

	return result, nil
}

func (b *BlogRepository) ListBlogs(ctx context.Context) ([]*repository.Blog, error) {
	blogs, err := b.queries.ListBlogs(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get blogs: %v", err)
	}

	var result []*repository.Blog

	for _, blog := range blogs {
		result = append(result, &repository.Blog{
			ID:        blog.ID,
			Author:    blog.Author,
			Title:     blog.Title,
			Content:   blog.Content,
			ImgUrls:   blog.ImgUrls,
			CreatedAt: blog.CreatedAt,
		})
	}

	return result, nil
}

func (b *BlogRepository) UpdateBlog(ctx context.Context, blog *repository.UpdateBlog) error {
	// check if the blog exists
	_, err := b.GetBlog(ctx, blog.ID)
	if err != nil {
		return err
	}

	var req generated.UpdateBlogParams

	req.ID = blog.ID

	if blog.Title != nil {
		req.Title = sql.NullString{
			Valid:  true,
			String: *blog.Title,
		}
	}

	if blog.Content != nil {
		req.Content = sql.NullString{
			Valid:  true,
			String: *blog.Content,
		}
	}

	if blog.ImgUrls != nil {
		req.ImgUrls = *blog.ImgUrls
	}

	err = b.queries.UpdateBlog(ctx, req)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update blog: %v", err)
	}

	return nil
}

func (b *BlogRepository) DeleteBlog(ctx context.Context, id uint32) error {
	// check if the blog exists
	_, err := b.GetBlog(ctx, id)
	if err != nil {
		return err
	}

	err = b.queries.DeleteBlog(ctx, id)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete blog: %v", err)
	}

	return nil
}
