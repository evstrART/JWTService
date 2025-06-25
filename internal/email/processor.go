package email

import (
	"JWTService/internal/models"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

func ProcessEmails(ch *amqp091.Channel, queueName string) {
	if ch == nil || ch.IsClosed() {
		log.Println("Chan is nil or closed")
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM_EMAIL")

	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" || smtpFrom == "" {
		log.Println("Invalid SMTP configuration")
		return
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("Invalod SMTP port '%s': %v.", smtpPortStr, err)
		return
	}

	for {
		msg, ok, err := ch.Get(queueName, false)

		if err != nil {
			log.Printf("Fail get inf from chan: %v", err)
			break
		}
		if !ok {
			log.Println("Queue is empty")
			break
		}

		var emailMsg models.EmailMessage
		if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
			log.Printf("Error Unmarshal message: %v", err)
			msg.Nack(false, false)
			continue
		}

		m := gomail.NewMessage()
		m.SetHeader("From", smtpFrom)
		m.SetHeader("To", emailMsg.RecipientEmail)
		m.SetHeader("Subject", emailMsg.Subject)
		m.SetBody("text/plain", emailMsg.Body)

		d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

		if err := d.DialAndSend(m); err != nil {
			log.Printf("Ошибка отправки письма для %s: %v", emailMsg.RecipientEmail, err)
			msg.Nack(false, true)
			continue
		}

		log.Printf("Письмо успешно отправлено для: %s ", emailMsg.RecipientEmail)
		msg.Ack(false)
	}
}
