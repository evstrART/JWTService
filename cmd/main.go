package main

import (
	"JWTService/internal/handler"
	"JWTService/internal/middleware"
	"JWTService/internal/repository"
	"JWTService/internal/service"
	"JWTService/pkg/postgres"
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

	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
	)

	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "внутренняя ошибка сервера",
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Fatal(app.Listen(":" + port))
}
