package handlers

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

const (
	ShutdownTimeout = 5 * time.Second
)

type MySQLRepository struct {
	u    repository.UserRepository
	p    repository.ProductRepository
	cart repository.CartRepository
	o    repository.OrderRepository
	cate repository.CategoryRepository
	r    repository.ReviewRepository
	b    repository.BlogRepository
}

type HttpServer struct {
	ln     net.Listener
	srv    *http.Server
	router *gin.Engine

	repo MySQLRepository
}

func NewHttpServer(addr string) *HttpServer {
	router := gin.Default()

	s := &HttpServer{
		router: router,

		srv: &http.Server{
			Addr:    addr,
			Handler: router.Handler(),
		},
	}

	return s
}

func (s *HttpServer) setRoutes() {
	s.router.GET("/health", s.healthCheckHandler)
}

func (s *HttpServer) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *HttpServer) Start(addr string) error {
	var err error
	if s.ln, err = net.Listen("tcp", addr); err != nil {
		return err
	}

	go func() {
		err := s.srv.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

func (s *HttpServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

func (s *HttpServer) SetDependencies(store *mysql.Store) {
	s.repo.u = mysql.NewUserRepository(store)
	s.repo.p = mysql.NewProductRepository(store)
	// order repo
	s.repo.cart = mysql.NewCartRepository(store)
	s.repo.cate = mysql.NewCategoryRepository(store)
	s.repo.r = mysql.NewReviewRepository(store)
	s.repo.b = mysql.NewBlogRepository(store)
}

func (s *HttpServer) Port() int {
	if s.ln == nil {
		return 0
	}

	port, _ := s.ln.Addr().(*net.TCPAddr)

	return port.Port
}
