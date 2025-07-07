package server

import (
	"sync"
	"venturo-core/configs"
	"venturo-core/internal/handler/http"
	"venturo-core/internal/middleware"
	"venturo-core/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func registerRoutes(app *fiber.App, db *gorm.DB, conf *configs.Config, wg *sync.WaitGroup) {
	app.Static("/public", "./public")
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Welcome to Venturo Core!",
		})
	})

	api := app.Group("/api/v1")

	// --- Setups ---
	authMiddleware := middleware.NewAuthMiddleware(conf.JWTSecretKey)

	// --- Setup services ---
	authService := service.NewAuthService(db, conf)
	userService := service.NewUserService(db, wg)
	postService := service.NewPostService(db)
	transactionService := service.NewTransactionService(db)

	// --- Setup handlers ---
	authHandler := http.NewAuthHandler(authService)
	userHandler := http.NewUserHandler(userService)
	postHandler := http.NewPostHandler(postService)
	transactionHandler := http.NewTransactionHandler(transactionService)

	// --- Auth routes ---
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// --- User routes ---
	api.Get("/profile", authMiddleware, userHandler.GetProfile)
	api.Put("/profile", authMiddleware, userHandler.UpdateProfile)

	// --- Post routes ---
	postRoutes := api.Group("/posts")
	postRoutes.Get("/", postHandler.GetAllPosts)                      // Public
	postRoutes.Get("/:id", postHandler.GetPostByID)                   // Public
	postRoutes.Post("/", authMiddleware, postHandler.CreatePost)      // Protected
	postRoutes.Put("/:id", authMiddleware, postHandler.UpdatePost)    // Protected
	postRoutes.Delete("/:id", authMiddleware, postHandler.DeletePost) // Protected

	// --- Transaction routes ---
	transactionRoutes := api.Group("/transactions")
	transactionRoutes.Post("/", authMiddleware, transactionHandler.CreateTransaction) // Protected
}
