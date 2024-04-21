package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Content string `json:"content" binding:"required"`
}

func main() {
	// Initialize Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka-cluster.cdc.svc.cluster.local:9092", // Replace with your Kafka broker address
	})
	if err != nil {
		fmt.Printf("Error creating producer: %v\n", err)
		os.Exit(1)
	}
	defer p.Close()

	// Initialize Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka-cluster.cdc.svc.cluster.local:9092", // Replace with your Kafka broker address
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Error creating consumer: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	// Subscribe to Kafka topic
	if err := c.SubscribeTopics([]string{"test-sample-client"}, nil); err != nil {
		fmt.Printf("Error subscribing to topic: %v\n", err)
		os.Exit(1)
	}

	// Run Kafka consumer in a separate goroutine
	go func() {
		for {
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
			case kafka.Error:
				fmt.Printf("Error: %v\n", e)
			}

			//fmt.Printf("Poll message: ")
		}
	}()

	// Initialize HTTP server
	router := gin.Default()

	// Define an API endpoint to produce a Kafka message
	router.POST("/message", func(c *gin.Context) {
		var msg Message

		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		topic := "test-sample-client" // Replace with your Kafka topic
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(msg.Content),
		}, nil)

		// Here you can handle the message, for example, print it out
		fmt.Printf("Send message: %s\n", msg.Content)
		c.JSON(http.StatusOK, gin.H{"status": "Send Message"})
	})

	// Start HTTP server
	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	fmt.Println("Shutting down gracefully...")
}
