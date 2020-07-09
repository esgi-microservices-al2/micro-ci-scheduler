package rabbit

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/System-Glitch/goyave/v2"
	"github.com/System-Glitch/goyave/v2/config"
	"github.com/streadway/amqp"
)

const (
	// QueueName the name of the AMQP queue used for message delivery.
	QueueName string = "al2.scheduler"
)

var (
	conn    *amqp.Connection = nil
	channel *amqp.Channel    = nil
	mu      sync.Mutex
)

// Connect to AMQP and open channel.
func Connect() {
	goyave.Logger.Println("Starting AMQP")
	mu.Lock()
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.GetString("rabbitUser"), config.GetString("rabbitPassword"), config.GetString("rabbitHost"), int(config.Get("rabbitPort").(float64)))

	var err error
	conn, err = amqp.Dial(connectionString)

	if err != nil {
		goyave.ErrLogger.Println(err, "Failed to open AMQP connection")
		mu.Unlock()
		time.Sleep(5 * time.Second)
		Connect()
		return
	}

	channel, err = conn.Channel()
	mu.Unlock()
	if err != nil {
		goyave.ErrLogger.Println(err, "Failed to open AMQP channel")
		Stop()
		return
	}
}

// Publish id to queue. Cannot be called if not connected.
func Publish(id int) {
	mu.Lock()
	defer mu.Unlock()
	queue, err := channel.QueueDeclare(
		QueueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		goyave.ErrLogger.Println(err, "Failed to declare AMQP queue")
		return
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(strconv.Itoa(id)),
	}

	err = channel.Publish("", queue.Name, false, false, msg)
	if err != nil {
		goyave.ErrLogger.Println(err, "Failed to publish to AMQP channel")
	}
}

// Stop close AMQP connection and channel.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	goyave.Logger.Println("Stopping AMQP")
	if channel != nil {
		channel.Close()
		channel = nil
	}
	if conn != nil {
		conn.Close()
		conn = nil
	}
}
