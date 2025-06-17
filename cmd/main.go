package main

import (
	"JWTService/internal/handler"
	"JWTService/internal/middleware"
	"JWTService/internal/repository"
	"JWTService/internal/service"
	"JWTService/pkg/postgres"
	"github.com/redis/go-redis/v9"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: "",
		DB:       0,
	})
	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
		rdb,
	)

	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		},
	})

	app.Use(logger.New())

	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/reg", authHandler.Register)
	auth.Post("/refresh", authHandler.Refresh)

	protected := auth.Group("/", middleware.AuthMiddleware(authService))
	protected.Post("/logout", authHandler.Logout)
	protected.Post("/logout_all", authHandler.LogoutAll)
	protected.Get("/test", authHandler.TestPrint)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Fatal(app.Listen(":" + port))
}
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
