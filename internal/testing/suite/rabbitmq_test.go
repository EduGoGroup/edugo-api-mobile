//go:build integration
// +build integration

package suite_test

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	testifySuite "github.com/stretchr/testify/suite"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQSuite prueba la funcionalidad de RabbitMQ usando la suite compartida
type RabbitMQSuite struct {
	suite.IntegrationTestSuite
}

// TestRabbitMQSuite ejecuta los tests de RabbitMQ
func TestRabbitMQSuite(t *testing.T) {
	testifySuite.Run(t, new(RabbitMQSuite))
}

// TestRabbitMQConnection verifica que RabbitMQ está disponible
func (s *RabbitMQSuite) TestRabbitMQConnection() {
	ctx := context.Background()
	
	// Obtener URL de conexión
	rabbitURL, err := s.RabbitContainer.AmqpURL(ctx)
	s.NoError(err, "Debe obtener URL de RabbitMQ")
	s.NotEmpty(rabbitURL, "URL no debe estar vacía")
	
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	s.NoError(err, "Debe conectar a RabbitMQ")
	s.NotNil(conn, "Conexión no debe ser nil")
	defer conn.Close()
	
	s.Logger.Info("✅ Conexión a RabbitMQ exitosa")
}

// TestPublishMessage verifica que se pueden publicar mensajes
func (s *RabbitMQSuite) TestPublishMessage() {
	ctx := context.Background()
	
	// Obtener URL y conectar
	rabbitURL, err := s.RabbitContainer.AmqpURL(ctx)
	s.NoError(err)
	
	conn, err := amqp.Dial(rabbitURL)
	s.NoError(err)
	defer conn.Close()
	
	// Crear canal
	ch, err := conn.Channel()
	s.NoError(err)
	defer ch.Close()
	
	// Declarar exchange
	exchangeName := "test-exchange"
	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	s.NoError(err, "Debe declarar exchange")
	
	// Publicar mensaje
	message := []byte(`{"test": "message", "timestamp": "2025-11-16"}`)
	err = ch.PublishWithContext(
		ctx,
		exchangeName,
		"test.routing.key",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
			Timestamp:   time.Now(),
		},
	)
	s.NoError(err, "Debe publicar mensaje correctamente")
	
	s.Logger.Info("✅ Mensaje publicado exitosamente", "exchange", exchangeName)
}

// TestConsumeMessage verifica que se pueden consumir mensajes
func (s *RabbitMQSuite) TestConsumeMessage() {
	ctx := context.Background()
	
	// Obtener URL y conectar
	rabbitURL, err := s.RabbitContainer.AmqpURL(ctx)
	s.NoError(err)
	
	conn, err := amqp.Dial(rabbitURL)
	s.NoError(err)
	defer conn.Close()
	
	// Crear canal
	ch, err := conn.Channel()
	s.NoError(err)
	defer ch.Close()
	
	// Declarar exchange
	exchangeName := "test-consume-exchange"
	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	s.NoError(err)
	
	// Declarar queue
	queue, err := ch.QueueDeclare(
		"test-queue",
		false, // durable
		true,  // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	s.NoError(err, "Debe declarar queue")
	
	// Bind queue a exchange
	err = ch.QueueBind(
		queue.Name,
		"test.#",
		exchangeName,
		false,
		nil,
	)
	s.NoError(err, "Debe hacer bind de queue")
	
	// Publicar mensaje
	testMessage := []byte(`{"event": "test.created", "data": "sample"}`)
	err = ch.PublishWithContext(
		ctx,
		exchangeName,
		"test.event.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        testMessage,
		},
	)
	s.NoError(err, "Debe publicar mensaje")
	
	// Consumir mensaje
	msgs, err := ch.Consume(
		queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	s.NoError(err, "Debe iniciar consumo")
	
	// Esperar mensaje con timeout
	select {
	case msg := <-msgs:
		s.Equal(testMessage, msg.Body, "Mensaje recibido debe ser igual al enviado")
		s.Logger.Info("✅ Mensaje consumido exitosamente")
	case <-time.After(3 * time.Second):
		s.Fail("Timeout esperando mensaje")
	}
}

// TestMultiplePublishers verifica que múltiples publicadores pueden trabajar simultáneamente
func (s *RabbitMQSuite) TestMultiplePublishers() {
	ctx := context.Background()
	
	rabbitURL, err := s.RabbitContainer.AmqpURL(ctx)
	s.NoError(err)
	
	exchangeName := "test-multi-exchange"
	
	// Función para publicar mensajes
	publishMessages := func(publisherID int, count int) error {
		conn, err := amqp.Dial(rabbitURL)
		if err != nil {
			return err
		}
		defer conn.Close()
		
		ch, err := conn.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()
		
		// Declarar exchange
		err = ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
		if err != nil {
			return err
		}
		
		// Publicar mensajes
		for i := 0; i < count; i++ {
			err = ch.PublishWithContext(
				ctx,
				exchangeName,
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte("message"),
				},
			)
			if err != nil {
				return err
			}
		}
		
		return nil
	}
	
	// Publicar desde 3 goroutines simultáneas
	done := make(chan error, 3)
	
	for i := 0; i < 3; i++ {
		go func(id int) {
			done <- publishMessages(id, 10)
		}(i)
	}
	
	// Esperar que todos terminen
	for i := 0; i < 3; i++ {
		err := <-done
		s.NoError(err, "Todos los publishers deben completar sin errores")
	}
	
	s.Logger.Info("✅ Múltiples publishers funcionando correctamente")
}
