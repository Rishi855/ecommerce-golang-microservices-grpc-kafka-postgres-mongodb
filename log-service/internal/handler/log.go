package handler

import (
	"context"
	"encoding/json"
	"log"
	"log-service/internal/initializers"
	"log-service/internal/model"
	"time"
)

type LogEvent struct {
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func HandleLog(key, value string) error {
	var event LogEvent
	var logCollection = initializers.LogCollection
	if err := json.Unmarshal([]byte(value), &event); err != nil {
		log.Printf("[Log-Service] Failed to unmarshal log event: %v", err)
		return err
	}

	// Convert timestamp string to time.Time
	parsedTime, _ := time.Parse(time.RFC3339, event.Timestamp)

	// Build log model
	logDoc := model.Log{
		Level:     event.Level,
		Service:   event.Service,
		Message:   event.Message,
		Timestamp: parsedTime,
	}

	// Insert into MongoDB
	_, err := logCollection.InsertOne(context.Background(), logDoc)
	if err != nil {
		log.Printf("[Log-Service] Failed to insert log into MongoDB: %v", err)
		return err
	}

	log.Printf("[Log-Service] Stored log in MongoDB: [%s] %s - %s", logDoc.Level, logDoc.Service, logDoc.Message)
	return nil
}
