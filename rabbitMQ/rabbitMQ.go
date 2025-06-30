package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing_RabbitMQ/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	rabbitMQHost     string
	rabbitMQPort     int
	RabbitMQExchange string

	rabbitMQDurable  bool
	rabbitMQUserName string
	rabbitMQPassword string

	RabbitMQConsumerQueue amqp.Queue

	RabbitMQConnection *amqp.Connection
	RabbitMQChannel    *amqp.Channel

	NotifyConfirm         chan amqp.Confirmation
	NotifyChannelClose    chan *amqp.Error
	NotifyConnectionClose chan *amqp.Error
	NotifyBlocked         chan amqp.Blocking
}

func (rmq *RabbitMQ) Initialize() bool {

	var err error

	//panic recovery
	defer func() {
		if err := recover(); err != nil {
			e := err.(error)
			fmt.Println("Error occurred in RabbitMQInitialization", e)
		}
	}()

	rmq.rabbitMQHost = config.CfgObj.RabbitMQConfigDetails.RabbitMQHost
	rmq.rabbitMQPort = config.CfgObj.RabbitMQConfigDetails.RabbitMQPort
	rmq.RabbitMQExchange = config.CfgObj.RabbitMQConfigDetails.RabbitMQExchange
	rmq.rabbitMQDurable = config.CfgObj.RabbitMQConfigDetails.RabbitMQDurable
	rmq.rabbitMQUserName = config.CfgObj.RabbitMQConfigDetails.RabbitMQUserName

	rmq.rabbitMQPassword, err = config.Decompress(config.CfgObj.RabbitMQConfigDetails.RabbitMQPassword)

	if err != nil {
		fmt.Println("Error decompressing password:", err)
		return false
	}

	return true
}

func (rmq *RabbitMQ) Connect() bool {
	var flag bool
	var err error

	if config.CfgObj.RabbitMQConfigDetails.RabbitMQSSLFlag == 0 {
		flag = rmq.RabbitMQConnectWithoutSSL()
	} else {
		flag = rmq.RabbitMQConnectWithSSl()
	}

	if flag {
		rmq.RabbitMQChannel, err = rmq.RabbitMQConnection.Channel()

		if err == nil {
			fmt.Println("RabbitMQChannel created successfully")

			if rmq.RabbitMQConnection.IsClosed() {
				//connection is not active
				fmt.Println("RabbitMQConnect failed with Host=[" + rmq.rabbitMQHost + "] and Port=[" + strconv.Itoa(rmq.rabbitMQPort) + "]")
				return false
			} else {
				//connection activated
				fmt.Println("RabbitMQ is connected successfully with Host=[" + rmq.rabbitMQHost + "] and Port=[" + strconv.Itoa(rmq.rabbitMQPort) + "]")
			}
		} else {
			fmt.Println("RabbitMQChannel creation failed")
			return false
		}

		err = rmq.RabbitMQChannel.Confirm(false)

		if err != nil {
			fmt.Println("RabbitMQChannel confirmation failed")
			return false
		}

		rmq.NotifyConfirm = make(chan amqp.Confirmation, 1)
		rmq.RabbitMQChannel.NotifyPublish(rmq.NotifyConfirm)
	}
	return true
}

func (rmq *RabbitMQ) RabbitMQConnectWithoutSSL() bool {

	//create connection
	conn, err := amqp.Dial("amqp://" + rmq.rabbitMQUserName + ":" + rmq.rabbitMQPassword + "@" + rmq.rabbitMQHost + ":" + strconv.Itoa(rmq.rabbitMQPort))
	if err != nil {
		//connection failed
		fmt.Println("RabbitMQConnect failed")
		return false
	}

	//store connection and channel
	rmq.RabbitMQConnection = conn
	return true
}

func (rmq *RabbitMQ) RabbitMQConnectWithSSl() bool {
	cfg := new(tls.Config)
	cfg.RootCAs = x509.NewCertPool()

	file, err := os.Open(config.CfgObj.RabbitMQConfigDetails.RabbitMQCACert)
	if err != nil { // handling error
		fmt.Println("Certificate load failed.")
		return false
	}
	defer file.Close()

	if ca, err := io.ReadAll(file); err == nil {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	} else {
		fmt.Println("Certificate read failed.")
		return false
	}
	if !config.CfgObj.RabbitMQConfigDetails.RABBITMQ_TLS_ONLY_CA_CERT {
		if cert, err := tls.LoadX509KeyPair(config.CfgObj.RabbitMQConfigDetails.RabbitMQClientCert, config.CfgObj.RabbitMQConfigDetails.RabbitMQClientKey); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		}
	}
	cfg.InsecureSkipVerify = config.CfgObj.RabbitMQConfigDetails.RABBITMQ_TLS_INSECURE_SKIP_VERIFY

	conn, err := amqp.DialTLS(fmt.Sprintf("amqp://%s:%s@%s:%d", rmq.rabbitMQUserName, rmq.rabbitMQPassword, rmq.rabbitMQHost, rmq.rabbitMQPort), cfg)

	if err != nil {
		fmt.Println("Error connecting to RabbitMQ:", err)
		return false
	}

	rmq.RabbitMQConnection = conn

	return true
}

func (rmq *RabbitMQ) QueueDeclare(queueName string) (amqp.Queue, bool) {
	rmQueue, err := rmq.RabbitMQChannel.QueueDeclare(queueName, config.CfgObj.RabbitMQConfigDetails.RabbitMQDurable, false, false, false, nil)

	if err != nil {
		fmt.Println("failed to declare Queue: ", err)
		return rmQueue, false
	}

	fmt.Println("queue delcared successful")

	return rmQueue, true
}

func (rmq *RabbitMQ) PublishMessage(queueName string, message string) bool {
	rmQueue, flag := rmq.QueueDeclare(queueName)

	if !flag {
		return false
	}

	err := rmq.RabbitMQChannel.PublishWithContext(
		context.Background(),
		"",
		rmQueue.Name,
		true,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			ContentEncoding: "UTF-8",
			DeliveryMode:    amqp.Persistent,
			Body:            []byte(message),
		},
	)

	if err != nil {
		fmt.Println("publish failed", err)
	}

	for {
		select {
		case confirm := <-rmq.NotifyConfirm:
			if confirm.Ack {
				fmt.Print("Message {" + message + "} published successfully in queue {" + rmQueue.Name + "}")
				return true
			}
		}
	}
}

func (rmq *RabbitMQ) ConsumeMessage(queueName string) {

	rmQueue, flag := rmq.QueueDeclare(queueName)

	if !flag {
		return
	}

	msgs, err := rmq.RabbitMQChannel.Consume(
		rmQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println("error consuming from queue ,", err)
	}

	rmq.NotifyChannelClose = make(chan *amqp.Error, 1)
	rmq.RabbitMQChannel.NotifyClose(rmq.NotifyChannelClose)

	rmq.NotifyConnectionClose = make(chan *amqp.Error, 1)
	rmq.RabbitMQConnection.NotifyClose(rmq.NotifyConnectionClose)

	rmq.NotifyBlocked = make(chan amqp.Blocking, 1)
	rmq.RabbitMQConnection.NotifyBlocked(rmq.NotifyBlocked)

	for {
		select {
		case msg := <-msgs:
			fmt.Println(string(msg.Body))

			err := rmq.RabbitMQChannel.Ack(msg.DeliveryTag, true)

			if err != nil {
				fmt.Println("error", "Error while giving acknowledgement for consumed message is "+string(err.Error()))
			}

		case <-rmq.NotifyChannelClose:
			fmt.Println("RabbitMQ Channel Connection is Closed")
			return

		case <-rmq.NotifyBlocked:
			fmt.Println("RabbitMQ connection blocked")
			return

		case <-rmq.NotifyConnectionClose:
			fmt.Println("RabbitMQ channler is closed")
			return
		}
	}
}

func (rmq *RabbitMQ) Disconnect() bool {
	if rmq.RabbitMQChannel != nil {
		rmq.RabbitMQChannel.Close()
	}

	rmq.RabbitMQConnection.Close()

	return true
}
