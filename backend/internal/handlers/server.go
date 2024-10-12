package handlers

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
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
	products := s.router.Group("/products")
	cart := s.router.Group("/categories")
	reviews := s.router.Group("/reviews")
	blogs := s.router.Group("/blogs")
	carts := s.router.Group("/carts")

	s.router.GET("/health", s.healthCheckHandler)

	// users routes
	users.GET("/", s.listUsers)
	users.POST("/register", s.createUser)
	users.GET("/:id", s.getUser)
	users.POST("/login", s.loginUser)
	users.GET("/:id/refresh-token", s.refreshToken)
	users.POST("/reset-password", s.resetPassword)
	users.PUT("/:id/update-subscription", s.updateUserSubscription)

	users.GET("/:id/reviews", s.ListUsersReviews)

	users.POST("/:id/blogs", s.createBlog)
	users.GET("/:id/blogs", s.getBlogsByAuthor)
	users.GET("/:id/blogs/:blogId", s.getBlog)
	users.DELETE("/:id/blogs/:blogId", s.deleteBlog)
	users.PUT("/:id/blogs/:blogId", s.updateBlog)

	users.GET("/:id/cart", s.getCart)
	users.PUT("/:id/cart", s.updateCart)
	users.POST("/:id/cart", s.createCart)
	users.DELETE("/:id/cart", s.deleteCart)

	// product routes
	products.GET("/", s.listProducts) // use query params
	products.POST("/create-product", s.createProduct)
	products.GET("/:id", s.getProduct)
	products.PUT("/:id", s.updateProduct)
	products.PUT("/:id/stock", s.updateProductQuantity)
	products.DELETE("/:id", s.deleteProduct)

	products.POST("/:id/reviews", s.createReview)
	products.GET("/:id/reviews", s.listProductsReviews)

	// categories routes
	cart.GET("/", s.listCategories)
	cart.POST("/create-category", s.createCategory)
	cart.GET("/:id", s.getCategory)
	cart.PUT("/:id", s.updateCategory)
	cart.DELETE("/:id", s.deleteCategory)

	// reviews routes
	reviews.GET("/", s.listReviews)
	reviews.GET("/:id", s.getReview)
	reviews.DELETE("/:id", s.deleteReview)

	// blogs route
	blogs.GET("/", s.listBlogs)

	// carts route
	carts.GET("/", s.listCarts)
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
	s.repo = MySQLRepository{
		u:    mysql.NewUserRepository(store),
		p:    mysql.NewProductRepository(store),
		cart: mysql.NewCartRepository(store),
		// o:    mysql.NewOrderRepository(store),
		cate: mysql.NewCategoryRepository(store),
		r:    mysql.NewReviewRepository(store),
		b:    mysql.NewBlogRepository(store),
	}
}

func (s *HttpServer) Port() int {
	if s.ln == nil {
		return 0
	}

	port, _ := s.ln.Addr().(*net.TCPAddr)

	return port.Port
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status_code": pkg.ErrorCode(err),
		"message":     pkg.ErrorMessage(err),
	}
}

func getParam(key string) (uint32, error) {
	intId, err := strconv.Atoi(key)
	if err != nil {
		return 0, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	if intId <= 0 {
		return 0, pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	return uint32(intId), nil
}
