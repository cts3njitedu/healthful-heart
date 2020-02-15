package connections

import (
  "github.com/streadway/amqp"
  "github.com/cts3njitedu/healthful-heart/utils"
  "log"
)

type RabbitConnection struct {
	environmentUtil utils.IEnvironmentUtility
}

func NewRabbitConnection(environmentUtil utils.IEnvironmentUtility) *RabbitConnection {
	return &RabbitConnection{environmentUtil}
}


func (rconn *RabbitConnection) GetConnection() (*amqp.Connection, error)  {
	mqUrl:=rconn.environmentUtil.GetEnvironmentString("CLOUDAMQP_URL");
	conn, err := amqp.Dial(mqUrl);
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to Rabbit MQ", err);
		return nil, err
	}

	return conn, nil
}