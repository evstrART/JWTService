package email

import (
	"JWTService/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func PublishEmailMessage(ch *amqp091.Channel, queueName string, msg models.EmailMessage) error {
	if ch == nil || ch.IsClosed() {
		return fmt.Errorf("Chan is nil or closed")
	}

	body, marshalErr := json.Marshal(msg)
	if marshalErr != nil {
		log.Printf("Error marshaling message to Email: %v", marshalErr)
		return marshalErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publishErr := ch.PublishWithContext(
		ctx,
		"",        // Exchange (по умолчанию, для прямой отправки в очередь)
		queueName, // Routing key (имя очереди)
		false,     // Mandatory
		false,     // Immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent, // Сохраняет сообщение на диск
		},
	)
	if publishErr != nil {
		log.Printf("Error publishing message in RabbitMQ to Email: %v", publishErr)
		return publishErr
	}
	log.Println("Сообщение о Email опубликовано в RabbitMQ.")
	return nil
}
