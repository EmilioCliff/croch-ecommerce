package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

type Blog struct {
	ID      uint32          `json:"id"`
	Author  uint32          `json:"author"`
	Title   string          `json:"title"`
	Content string          `json:"content"`
	ImgUrls json.RawMessage `json:"img_urls"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}

func (p *Blog) UnmarshalOptions() ([]string, error) {
	var imgUrls []string
	if err := json.Unmarshal(p.ImgUrls, &imgUrls); err != nil {
		return nil, fmt.Errorf("failed to unmarshal img_urls: %w", err)
	}

	return imgUrls, nil
}

func (p *Blog) Validate() error {
	if p.Author <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "author is required")
	}

	if p.Title == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "title is required")
	}

	if p.Content == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "content is required")
	}

	return nil
}

type UpdateBlog struct {
	ID      uint32           `json:"id"`
	Title   *string          `json:"title"`
	Content *string          `json:"content"`
	ImgUrls *json.RawMessage `json:"img_urls"`
}

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *Blog) (*Blog, error)
	GetBlog(ctx context.Context, id uint32) (*Blog, error)
	GetBlogsByAuthor(ctx context.Context, author uint32) ([]*Blog, error)
	ListBlogs(ctx context.Context) ([]*Blog, error)
	UpdateBlog(ctx context.Context, blog *UpdateBlog) error
	DeleteBlog(ctx context.Context, id uint32) error
}
