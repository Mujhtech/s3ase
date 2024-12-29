package response

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/render"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/rs/zerolog/log"
)

type Response struct {
	StatusCode int `json:"-"`
}

type ServerResponse struct {
	Response
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
}

func (res Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	render.Status(r, res.StatusCode)
	return nil
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) error {
	_ = render.Render(w, r, ServerResponse{
		Response: Response{
			StatusCode: http.StatusBadRequest,
		},
		Message: "Bad Request",
		Error:   err.Error(),
	})

	return nil
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) error {
	_ = render.Render(w, r, ServerResponse{
		Response: Response{
			StatusCode: http.StatusInternalServerError,
		},
		Message: "Internal Server Error",
		Error:   err.Error(),
	})

	return nil
}

func Ok(w http.ResponseWriter, r *http.Request, message string, data interface{}) error {
	_ = render.Render(w, r, ServerResponse{
		Response: Response{
			StatusCode: http.StatusOK,
		},
		Message: message,
		Data:    data,
	})

	return nil
}

func Created(w http.ResponseWriter, r *http.Request, message string, data interface{}) error {
	_ = render.Render(w, r, ServerResponse{
		Response: Response{
			StatusCode: http.StatusCreated,
		},
		Message: message,
		Data:    data,
	})

	return nil
}

func Unauthorized(w http.ResponseWriter, r *http.Request, err error) error {
	_ = render.Render(w, r, ServerResponse{
		Response: Response{
			StatusCode: http.StatusUnauthorized,
		},
		Message: "Unauthorized",
		Error:   err.Error(),
	})

	return nil
}

func Redirect(w http.ResponseWriter, r *http.Request, u string, status int, ignoreUrl bool) error {
	if !ignoreUrl {
		parsedUrl, err := url.ParseRequestURI(u)
		if err != nil || parsedUrl.Host != "" {
			return BadRequest(w, r, fmt.Errorf("invalid redirect URL"))
		}
	}

	http.Redirect(w, r, u, status)
	return nil
}

func Stream(
	ctx context.Context,
	w http.ResponseWriter,
	chStop <-chan struct{},
	chEvents <-chan *sse.Event,
	chErr <-chan error,
) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	// Set the headers related to event streaming.
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")
	h.Set("Access-Control-Allow-Origin", "*")

	const (
		pingInterval = 30 * time.Second
		tailMaxTime  = 2 * time.Hour
	)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	stream := sseStream{
		enc:     enc,
		writer:  w,
		flusher: flusher,
	}

	ctx, ctxCancel := context.WithTimeout(ctx, tailMaxTime)
	defer ctxCancel()

	if err := stream.ping(); err != nil {
		return
	}

	defer func() {
		if err := stream.close(); err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to close SSE stream")
		}
	}()

	pingTimer := time.NewTimer(pingInterval)
	defer pingTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Ctx(ctx).Debug().Err(ctx.Err()).Msg("stream SSE request context done")
			return

		case <-chStop:
			log.Ctx(ctx).Debug().Msg("app shutdown")
			return

		case err := <-chErr:
			log.Ctx(ctx).Debug().Err(err).Msg("received error from SSE stream")
			return

		case <-pingTimer.C:
			if err := stream.ping(); err != nil {
				log.Ctx(ctx).Err(err).Msg("failed to send SSE ping")
				return
			}

		case event := <-chEvents:
			if err := stream.event(event); err != nil {
				log.Ctx(ctx).Err(err).Msgf("failed to send SSE event: %s", event.Type)
				return
			}
		}

		pingTimer.Stop() // stop timer

		select {
		case <-pingTimer.C: // drain channel
		default:
		}

		pingTimer.Reset(pingInterval) // reset timer
	}

}

type sseStream struct {
	enc     *json.Encoder
	writer  io.Writer
	flusher http.Flusher
}

func (r sseStream) event(event *sse.Event) error {
	_, err := io.WriteString(r.writer, fmt.Sprintf("event: %s\n", event.Type))
	if err != nil {
		return fmt.Errorf("failed to send event header: %w", err)
	}

	_, err = io.WriteString(r.writer, "data: ")
	if err != nil {
		return fmt.Errorf("failed to send data header: %w", err)
	}

	err = r.enc.Encode(event.Data)
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	// NOTE: enc.Encode is ending the data with a new line, only add one more
	// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/encoding/json/stream.go;l=220
	_, err = r.writer.Write([]byte{'\n'})
	if err != nil {
		return fmt.Errorf("failed to send end of message: %w", err)
	}

	r.flusher.Flush()
	return nil
}

func (r sseStream) close() error {
	_, err := io.WriteString(r.writer, "event: error\ndata: eof\n\n")
	if err != nil {
		return fmt.Errorf("failed to send EOF: %w", err)
	}
	r.flusher.Flush()
	return nil
}

func (r sseStream) ping() error {
	_, err := io.WriteString(r.writer, ": ping\n\n")
	if err != nil {
		return fmt.Errorf("failed to send ping: %w", err)
	}
	r.flusher.Flush()
	return nil
}
