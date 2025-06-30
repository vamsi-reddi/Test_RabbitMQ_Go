package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing_RabbitMQ/config"
	"testing_RabbitMQ/process"
	rabbitmq "testing_RabbitMQ/rabbitMQ"
)

func main(){
	if(len(os.Args) != 2) {
		fmt.Println("invalid command")
		return
	}

	configFile := os.Args[1]

	var secretData string
	var err error

	if filepath.Ext(configFile) == "json"{
		secretData, err = config.ReadConfigFile(configFile)

		if err != nil {
			return
		}
	} else{
		fmt.Println("invalid file format")
		return
	}

	if !config.LoadConfig(secretData) {
		return
	}

	process.RabbitMQ = &rabbitmq.RabbitMQ{}

	if !process.RabbitMQ.Initialize() {
		return
	}
	
	if !process.RabbitMQ.Connect() {
		return
	}

	defer process.RabbitMQ.Disconnect()

	if !process.RabbitMQ.PublishMessage("test-vamsi-queue", "Hello World!") {
		return
	}

	process.RabbitMQ.ConsumeMessage("test-vamsi-queue")
	
}