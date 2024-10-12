package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/gin-gonic/gin"
)

type createBlogRequest struct {
	Title   string   `binding:"required" json:"title"`
	Content string   `binding:"required" json:"content"`
	ImgUrls []string `binding:""         json:"img_urls"`
}

func (s *HttpServer) createBlog(ctx *gin.Context) {
	// get user id from token
	var id uint32 = 3

	var req createBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	data := &repository.Blog{
		Author:  id,
		Title:   req.Title,
		Content: req.Content,
	}

	err := data.MarshalOptions(req.ImgUrls)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	blog, err := s.repo.b.CreateBlog(ctx, data)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, blog)
}

func (S *HttpServer) listBlogs(ctx *gin.Context) {
	blogs, err := S.repo.b.ListBlogs(ctx)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func (s *HttpServer) getBlog(ctx *gin.Context) {
	id, err := getParam(ctx.Param("blogId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	blog, err := s.repo.b.GetBlog(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, blog)
}

func (s *HttpServer) updateBlog(ctx *gin.Context) {
	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	var req createBlogRequest
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	id, err := getParam(ctx.Param("blogId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	data := &repository.UpdateBlog{
		ID:      id,
		Title:   pkg.StringPtr(req.Title),
		Content: pkg.StringPtr(req.Content),
	}

	err = data.MarshalOptions(req.ImgUrls)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err)))

		return
	}

	err = s.repo.b.UpdateBlog(ctx, &repository.UpdateBlog{
		ID:      id,
		Title:   pkg.StringPtr(req.Title),
		Content: pkg.StringPtr(req.Content),
		ImgUrls: data.ImgUrls,
	})
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *HttpServer) getBlogsByAuthor(ctx *gin.Context) {
	id, err := getParam(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	blogs, err := s.repo.b.GetBlogsByAuthor(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func (s *HttpServer) deleteBlog(ctx *gin.Context) {
	id, err := getParam(ctx.Param("blogId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	err = s.repo.b.DeleteBlog(ctx, id)
	if err != nil {
		ctx.JSON(pkg.PkgErrorToHttpError(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
