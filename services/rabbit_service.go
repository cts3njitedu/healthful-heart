package services

import (
	"github.com/cts3njitedu/healthful-heart/connections"
	"github.com/cts3njitedu/healthful-heart/models"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"encoding/json"
	"github.com/cts3njitedu/healthful-heart/utils"
)

type RabbitService struct {
	rabbitConnection connections.IRabbitConnection
	environmentUtil utils.IEnvironmentUtility

}

var (
	_conn *amqp.Connection;
)

func NewRabbitService(rabbitConnection connections.IRabbitConnection, environmentUtil utils.IEnvironmentUtility) *RabbitService {
	log.Println("Making rabbit connection...");
	conn, err := rabbitConnection.GetConnection();

	if err!= nil {
		fmt.Println(err)
	}
	_conn = conn;
	
	go pullFileMetaDataFromQueue(environmentUtil);
	
	return &RabbitService{rabbitConnection, environmentUtil}
}

func (rabbitService *RabbitService) PushFileMetaDataToQueue(file *models.WorkoutFile) error {
	conn, err := rabbitService.rabbitConnection.GetConnection();

	if err!= nil {
		fmt.Println(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()

	if err!=nil {
		fmt.Println(err)
		return err;
	}

	defer ch.Close()

	exchangeName := rabbitService.environmentUtil.GetEnvironmentString("EXCHANGE_NAME")
	routingKey := rabbitService.environmentUtil.GetEnvironmentString("ROUTING_KEY") 


	err = ch.ExchangeDeclare (
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err!=nil {
		fmt.Println(err)
		return err;
	}


	body, err := json.Marshal(&file)

	if err!=nil {
		fmt.Println(err)
		return err;
	}

	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
		},

	)

	if err!=nil {
		fmt.Println(err)
		return err;
	}

	log.Printf(" [x] Sent %s", body)
	// data := []string{"a","b","c"}
	// go func() {
	// 	time.Sleep(5*time.Second)
	// 	for d := range data {
	// 			log.Printf("This is the string: %s ", data[d])
	// 	}
	// }()

	return nil
}

func pullFileMetaDataFromQueue(environmentUtil utils.IEnvironmentUtility) {

	// log.Println("Begining process to retrieve messages from queue....")
	// conn, err := rabbitService.rabbitConnection.GetConnection();

	exchangeName := environmentUtil.GetEnvironmentString("EXCHANGE_NAME")
	routingKey := environmentUtil.GetEnvironmentString("ROUTING_KEY") 
	queueName := environmentUtil.GetEnvironmentString("QUEUE_NAME")
	

	ch, err := _conn.Channel()

	if err!=nil {
		fmt.Println("Two: ", err)
	}

	defer ch.Close()

	// err = ch.ExchangeDeclare (
	// 	exchangeName,
	// 	"direct",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )

	// if err!=nil {
	// 	fmt.Println("Three: ", err)
	// }

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err!=nil {
		fmt.Println("Four: ", err)
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if err!=nil {
		fmt.Println("Five: ", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err!=nil {
		fmt.Println("Six: ", err)
	}

	forever := make(chan bool)

	go func() {

			for d := range msgs {
					log.Printf(" [x] Received %s %T", d.Body, d.Body)
			}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}