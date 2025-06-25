package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func NewRabbitMQClient(queueName string) (*amqp091.Connection, *amqp091.Channel, amqp091.Queue, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	var conn *amqp091.Connection
	var err error

	conn, err = amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Println("Error connecting to RabbitMQ")

	}
	log.Println("Successful connection to RabbitMQ.")

	if err != nil {
		return nil, nil, amqp091.Queue{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, amqp091.Queue{}, err
	}
	log.Println("Chan RabbitMQ is open.")

	q, err := ch.QueueDeclare(
		queueName, // Название очереди
		true,      // Durable (сохраняет сообщения при перезапуске брокера)
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, amqp091.Queue{}, err
	}
	log.Printf("Queue RabbitMQ '%s' is created.", q.Name)

	return conn, ch, q, nil
}

func CloseRabbitMQConnections(conn *amqp091.Connection, ch *amqp091.Channel) {
	if ch != nil && !ch.IsClosed() {
		log.Println("Close RabbitMQ chan.")
		if err := ch.Close(); err != nil {
			log.Printf("Fail close RabbitMQ chan: %v", err)
		}
	}
	if conn != nil && !conn.IsClosed() {
		log.Println("Close RabbitMQ conn.")
		if err := conn.Close(); err != nil {
			log.Printf("Fail close RabbitMQ conn: %v", err)
		}
	}
}
