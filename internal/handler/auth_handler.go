package handler

import (
	"JWTService/internal/email"
	"JWTService/internal/models"
	"JWTService/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

type AuthHandler struct {
	authService     *service.AuthService
	rabbitCh        *amqp091.Channel
	rabbitQueueName string
}

func NewAuthHandler(authService *service.AuthService, rabbitCh *amqp091.Channel, rabbitQueueName string) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		rabbitCh:        rabbitCh,
		rabbitQueueName: rabbitQueueName,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input models.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	tokens, err := h.authService.Register(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	emailMsg := models.EmailMessage{
		RecipientEmail: input.Email,
		Subject:        "Добро пожаловать!",
		Body:           "Спасибо за регистрацию!",
	}
	err = email.PublishEmailMessage(h.rabbitCh, h.rabbitQueueName, emailMsg)
	if err != nil {
		log.Printf("Ошибка публикации email: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(tokens)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	tokens, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	emailMsg := models.EmailMessage{
		RecipientEmail: req.Email,
		Subject:        "Добро пожаловать!",
		Body:           "Спасибо за то, что пользуетесь нашим сервисом!",
	}
	err = email.PublishEmailMessage(h.rabbitCh, h.rabbitQueueName, emailMsg)
	if err != nil {
		log.Printf("Ошибка публикации email: %v", err)
	}
	return c.JSON(tokens)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	tokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := h.authService.Logout(c.Context(), req.RefreshToken); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := h.authService.LogoutAll(c.Context(), req.RefreshToken); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
func (h *AuthHandler) TestPrint(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "done",
	})
}
