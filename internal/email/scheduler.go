package email

import (
	"github.com/go-co-op/gocron"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func StartEmailScheduler(ch *amqp091.Channel, queueName string) *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Minutes().Do(ProcessEmails, ch, queueName)
	if err != nil {
		log.Fatalf("Error when scheduling a dispatch task Email: %v", err)
	}

	s.StartAsync()
	return s
}
