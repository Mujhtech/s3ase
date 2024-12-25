package job

import "time"

type JobName string
type QueueName string

const (
	JobNameWebhook JobName = "webhook"
	JobNameAppSync JobName = "app_sync"

	QueueNameDefault QueueName = "default"
)

type ClientPayload struct {
	Data  []byte        `json:"data"`
	Delay time.Duration `json:"delay"`
}
