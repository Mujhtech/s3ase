package job

import "time"

type JobName string
type QueueName string

const (
	JobNameWebhook   JobName   = "webhook"
	QueueNameDefault QueueName = "default"
)

type ClientPayload struct {
	Data  []byte        `json:"data"`
	Delay time.Duration `json:"delay"`
}
