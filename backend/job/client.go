package job

import (
	"encoding/json"

	"github.com/danvixent/asynqmon"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/mujhtech/s3ase/internal/pkg/encrypt"
)

type Client struct {
	redisConnOpt asynq.RedisConnOpt
	client       *asynq.Client
	inspector    *asynq.Inspector
	aesCfb       encrypt.Encrypt
}

func NewClient(opts asynq.RedisConnOpt, aesCfb encrypt.Encrypt) *Client {

	return &Client{
		redisConnOpt: opts,
		client:       asynq.NewClient(opts),
		inspector:    asynq.NewInspector(opts),
		aesCfb:       aesCfb,
	}
}

func (c *Client) Enqueue(queue QueueName, job JobName, payload *ClientPayload) error {

	id := uuid.New().String()

	q := string(queue)

	data, err := c.aesCfb.Encrypt(payload.Data)

	if err != nil {
		return err
	}

	t := asynq.NewTask(string(job), []byte(data), asynq.Queue(q), asynq.TaskID(id), asynq.ProcessIn(payload.Delay))

	_, err = c.inspector.GetTaskInfo(q, id)
	if err != nil {

		// message := err.Error()
		// if ErrQueueNotFound.Error() == message || ErrTaskNotFound.Error() == message {
		// 	_, err := c.client.Enqueue(t, nil)
		// 	return err
		// }

		return err
	}

	// Delete the task if it already exists
	err = c.inspector.DeleteTask(q, id)

	if err != nil {
		return err
	}

	// Enqueue the task
	if _, err := c.client.Enqueue(t, nil); err != nil {
		return err
	}

	return nil
}

type Formatter struct {
	aesCfb encrypt.Encrypt
}

func (c *Client) Monitor() *asynqmon.HTTPHandler {
	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/queue/monitoring",
		RedisConnOpt: c.redisConnOpt,
		PayloadFormatter: Formatter{
			aesCfb: c.aesCfb,
		},
		ResultFormatter: Formatter{
			aesCfb: c.aesCfb,
		},
	})
	return h
}

func (q *Client) Inspector() *asynq.Inspector {
	return q.inspector
}

func (f Formatter) FormatPayload(_ string, payload []byte) string {
	data, err := f.aesCfb.Decrypt(string(payload))

	if err != nil {
		return ""
	}

	bytes, _ := json.Marshal(data)
	return string(bytes)
}

func (f Formatter) FormatResult(_ string, payload []byte) string {
	data, err := f.aesCfb.Decrypt(string(payload))

	if err != nil {
		return ""
	}

	bytes, _ := json.Marshal(data)
	return string(bytes)
}
