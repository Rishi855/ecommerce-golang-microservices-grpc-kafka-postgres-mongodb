package model

import "time"

type Log struct {
	ID        string    `bson:"_id,omitempty"`
	Level     string    `bson:"level"`
	Service   string    `bson:"service"`
	Message   string    `bson:"message"`
	Timestamp time.Time `bson:"timestamp"`
}
