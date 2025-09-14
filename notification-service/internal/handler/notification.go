package handler

import (
	"encoding/json"
	"log"
	"time"
)

type NotificationEvent struct {
	Event     string                 `json:"event"`
	UserID    int64                  `json:"user_id"`
	Channel   string                 `json:"channel"`
	Template  string                 `json:"template"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
}

func HandleNotification(key, value string) error {
	var event NotificationEvent
	if err := json.Unmarshal([]byte(value), &event); err != nil {
		log.Printf("Failed to unmarshal notification event: %v", err)
		return err
	}

	log.Printf("[Notification-Service] Processing notification for user %d via %s", event.UserID, event.Channel)
	log.Printf("Template: %s, Data: %+v", event.Template, event.Data)
	log.Printf("Simulating sending %s notification to user %d ...", event.Channel, event.UserID)
	time.Sleep(5*time.Second)
	return nil
}