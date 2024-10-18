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
	// routes groups
	users := s.router.Group("/users")
	usersAuth := s.router.Group("/users").Use(authMiddleware(s.tokenMaker))

	products := s.router.Group("/products")
	productsAuth := s.router.Group("/products").Use(authMiddleware(s.tokenMaker))

	cart := s.router.Group("/categories")
	cartAuth := s.router.Group("/categories").Use(authMiddleware(s.tokenMaker))

	reviews := s.router.Group("/reviews")
	reviewsAuth := s.router.Group("/reviews").Use(authMiddleware(s.tokenMaker))

	blogs := s.router.Group("/blogs")
	// blogsAuth := blogs.Use(authMiddleware(s.tokenMaker))
	carts := s.router.Group("/carts")
	// cartsAuth := carts.Use(authMiddleware(s.tokenMaker))

	s.router.GET("/health", s.healthCheckHandler)

	// users routes
	usersAuth.GET("/", s.listUsers)
	users.POST("/register", s.createUser)
	usersAuth.GET("/:id", s.getUser)
	users.POST("/login", s.loginUser)
	usersAuth.GET("/:id/refresh-token", s.refreshToken)
	users.POST("/reset-password", s.resetPassword)
	usersAuth.PUT("/:id/update-subscription", s.updateUserSubscription)

	usersAuth.GET("/:id/reviews", s.listUsersReviews)

	usersAuth.POST("/:id/blogs", s.createBlog)
	users.GET("/:id/blogs", s.getBlogsByAuthor)
	users.GET("/:id/blogs/:blogId", s.getBlog)
	usersAuth.DELETE("/:id/blogs/:blogId", s.deleteBlog)
	usersAuth.PUT("/:id/blogs/:blogId", s.updateBlog)

	usersAuth.GET("/:id/cart", s.getCart)
	usersAuth.PUT("/:id/cart", s.updateCart)
	usersAuth.POST("/:id/cart", s.createCart)
	usersAuth.DELETE("/:id/cart", s.deleteCart)

	// product routes
	products.GET("/", s.listProducts) // use query params
	productsAuth.POST("/create-product", s.createProduct)
	products.GET("/:id", s.getProduct)
	productsAuth.PUT("/:id", s.updateProduct)
	productsAuth.PUT("/:id/stock", s.updateProductQuantity)
	productsAuth.DELETE("/:id", s.deleteProduct)

	productsAuth.POST("/:id/reviews", s.createReview)
	products.GET("/:id/reviews", s.listProductsReviews)

	// categories routes
	cart.GET("/", s.listCategories)
	cartAuth.POST("/create-category", s.createCategory)
	cart.GET("/:id", s.getCategory)
	cartAuth.PUT("/:id", s.updateCategory)
	cartAuth.DELETE("/:id", s.deleteCategory)

	// reviews routes
	reviews.GET("/", s.listReviews)
	reviews.GET("/:id", s.getReview)
	reviewsAuth.DELETE("/:id", s.deleteReview)

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

func getPayload(ctx *gin.Context) (*pkg.Payload, error) {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return &pkg.Payload{}, pkg.Errorf(pkg.INVALID_ERROR, "authorization payload not found")
	}

	p, ok := payload.(*pkg.Payload)
	if !ok {
		return &pkg.Payload{}, pkg.Errorf(pkg.INVALID_ERROR, "authorization payload not found")
	}

	return p, nil
}
