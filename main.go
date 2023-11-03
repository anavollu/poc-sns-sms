package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", index)
	r.POST("/send-sms-message", sendSMSMessage)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "POC AWS SNS",
	})
}

func pubTextSMS(snsClient *sns.Client, message string, phoneNumber string) (*sns.PublishOutput, error) {
	input := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}

	result, err := snsClient.Publish(context.Background(), input)
	if err != nil {
		log.Fatalf("failed to send message, %v", err)
		return nil, err
	}

	fmt.Printf("%s Message sent. Status was %v\n", *result.MessageId, result.ResultMetadata)
	return result, nil
}

func sendSMSMessage(c *gin.Context) {
	var json struct {
		Message     string `json:"message" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	result, err := pubTextSMS(snsClient, json.Message, json.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messageId":  result.MessageId,
		"status":     "Message sent.",
		"statusCode": result.ResultMetadata,
	})
}
