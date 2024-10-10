package handlers

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
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
	ln         net.Listener
	srv        *http.Server
	router     *gin.Engine
	tokenMaker pkg.Maker
	config     pkg.Config

	repo MySQLRepository
}

func NewHttpServer(maker pkg.Maker, config pkg.Config) *HttpServer {
	router := gin.Default()

	s := &HttpServer{
		router: router,

		srv: &http.Server{
			Addr:    config.HTTP_PORT,
			Handler: router.Handler(),
		},
		tokenMaker: maker,
		config:     config,
	}

	s.setRoutes()

	return s
}

func (s *HttpServer) setRoutes() {
	users := s.router.Group("/users")

	s.router.GET("/health", s.healthCheckHandler)

	users.GET("/", s.listUsers)
	users.POST("/register", s.createUser)
	users.GET("/:id", s.getUser)
	users.POST("/login", s.loginUser)
	users.GET("/refresh-token", s.refreshToken)
	users.POST("/reset-password", s.resetPassword)
	users.POST("/update-subscription", s.updateUserSubscription)
}

func (s *HttpServer) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *HttpServer) Start() error {
	var err error
	if s.ln, err = net.Listen("tcp", s.config.HTTP_PORT); err != nil {
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
	log.Println("Shutting down http server...")

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

func errorResponse(err error) gin.H {
	return gin.H{
		"status_code": pkg.ErrorCode(err),
		"message":     pkg.ErrorMessage(err),
	}
}

func (s *HttpServer) Port() int {
	if s.ln == nil {
		return 0
	}

	port, _ := s.ln.Addr().(*net.TCPAddr)

	return port.Port
}
