package s3store

import (
	"context"
	"io"
)

type HookEvent struct{}

type MetaData map[string]string

type FileInfo struct {
	ID             string
	Size           int64
	SizeIsDeferred bool

	Offset   int64
	MetaData MetaData

	IsPartial bool
	IsFinal   bool

	PartialUploads []string

	Storage map[string]string

	//stopUpload func(HTTPResponse)
}

type HTTPHeader map[string]string

type HTTPResponse struct {
	StatusCode int

	Body string

	Header HTTPHeader
}

type Upload interface {
	WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error)

	GetInfo(ctx context.Context) (FileInfo, error)

	GetReader(ctx context.Context) (io.ReadCloser, error)

	FinishUpload(ctx context.Context) error
}
