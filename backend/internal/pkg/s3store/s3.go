package s3store

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsType "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mujhtech/s3ase/config"

	tusdStore "github.com/tus/tusd/v2/pkg/s3store"
)

type Options func(*S3Store)

type S3Store struct {
	bucket string
	client *s3.Client
	cfg    *config.Config
}

func NewS3Store(cfg *config.Config, opts ...Options) (*S3Store, error) {
	config, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.Aws.AccessKey,
				cfg.Aws.SecretKey,
				"",
			),
		),
		awsConfig.WithRegion(cfg.Aws.DefaultRegion),
	)

	if err != nil {
		return nil, err
	}

	return &S3Store{
		client: s3.NewFromConfig(config, func(o *s3.Options) {
			o.UseAccelerate = false

			// Disable HTTPS and only use HTTP (helpful for debugging requests).
			o.EndpointOptions.DisableHTTPS = false

			if cfg.Aws.Endpoint != "" {
				o.BaseEndpoint = &cfg.Aws.Endpoint
			}
			o.UsePathStyle = cfg.Aws.UsePathStyle
			o.Region = cfg.Aws.DefaultRegion
		}),
		cfg: cfg,
	}, nil
}

func (s *S3Store) GetClient() *s3.Client {
	return s.client
}

func (s *S3Store) OpenFile(ctx context.Context, bucket, key, region string) (*s3.GetObjectOutput, error) {
	if bucket == "" || key == "" {
		return nil, fmt.Errorf("bucket and object key are required")
	}
	return s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, func(options *s3.Options) {
		if region != "" {
			options.Region = region
		}
	})
}

func (s *S3Store) DeleteFile(ctx context.Context, bucket, key, region string) error {
	if bucket == "" || key == "" {
		return fmt.Errorf("bucket and object key are required")
	}
	deleteObject := func(objectKey string) error {
		_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &bucket,
			Key:    &objectKey,
		}, func(options *s3.Options) {
			if region != "" {
				options.Region = region
			}
		})
		return err
	}
	if err := deleteObject(key); err != nil {
		return err
	}
	if err := deleteObject(key + ".info"); err != nil {
		return err
	}
	return nil
}

// DeleteBucket removes every object (including tus metadata sidecars) before
// deleting the bucket itself. S3 only permits deletion of empty buckets.
func (s *S3Store) DeleteBucket(ctx context.Context, bucket, region string) error {
	if bucket == "" {
		return fmt.Errorf("bucket is required")
	}
	opts := func(options *s3.Options) {
		if region != "" {
			options.Region = region
		}
	}
	var continuationToken *string
	for {
		result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: &bucket, ContinuationToken: continuationToken,
		}, opts)
		if err != nil {
			return err
		}
		if len(result.Contents) > 0 {
			objects := make([]awsType.ObjectIdentifier, 0, len(result.Contents))
			for _, object := range result.Contents {
				objects = append(objects, awsType.ObjectIdentifier{Key: object.Key})
			}
			quiet := true
			if _, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: &bucket, Delete: &awsType.Delete{Objects: objects, Quiet: &quiet},
			}, opts); err != nil {
				return err
			}
		}
		if result.IsTruncated == nil || !*result.IsTruncated {
			break
		}
		continuationToken = result.NextContinuationToken
	}
	_, err := s.client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: &bucket}, opts)
	return err
}

func WithBucket(bucket string) Options {
	return func(s3 *S3Store) {
		s3.bucket = bucket
	}
}

func WithRegion(region string) Options {
	return func(s3 *S3Store) {
		// s3.client.ForcePathStyle = true
		// s3.client.Region = region
	}
}

func (s *S3Store) ListBuckets(ctx context.Context) ([]string, error) {
	results, err := s.client.ListBuckets(ctx, nil)

	if err != nil {
		return nil, err
	}

	var buckets []string

	for _, bucket := range results.Buckets {
		buckets = append(buckets, *bucket.Name)
	}

	return buckets, nil
}

func (s *S3Store) CheckOrCreateNewBucket(ctx context.Context, bucket string, region string) (string, error) {

	opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.Region = region
		},
	}

	existBucket, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucket,
	}, opts...)

	if err != nil || existBucket == nil {
		input := &s3.CreateBucketInput{Bucket: &bucket}
		if region != "" && region != "us-east-1" {
			input.CreateBucketConfiguration = &awsType.CreateBucketConfiguration{
				LocationConstraint: awsType.BucketLocationConstraint(region),
			}
		}

		_, err := s.client.CreateBucket(ctx, input, func(o *s3.Options) {
			o.Region = region
		})

		if err != nil {
			return "", err
		}

		// set default bucket policy
		if err = s.SetBucketDefaultPolicy(ctx, bucket); err != nil {
			return "", err
		}

		return "", nil

	}

	if existBucket.BucketRegion != nil {
		return *existBucket.BucketRegion, nil
	}
	return region, nil
}

func (s *S3Store) SetBucketDefaultPolicy(ctx context.Context, bucket string) error {
	// policy := map[string]interface{}{
	// 	"Version": "2012-10-17",
	// 	"Statement": []map[string]interface{}{
	// 		{
	// 			"Effect":    "Allow",
	// 			"Principal": "*",
	// 			"Action": []string{
	// 				"s3:GetObject",
	// 				"s3:PutObject",
	// 				"s3:DeleteObject",
	// 			},
	// 			"Resource": "arn:aws:s3:::" + bucket + "/*",
	// 		},
	// 	},
	// }

	// policyJSON, err := json.Marshal(policy)
	// if err != nil {
	// 	return err
	// }

	// _, err = s.client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
	// 	Bucket: &bucket,
	// 	Policy: aws.String(string(policyJSON)),
	// 	CreateBucketConfiguration: &awsType.CreateBucketConfiguration{
	// 		LocationConstraint: awsType.BucketLocationConstraint(region),
	// 	},
	// })

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *S3Store) NewApp(ctx context.Context, bucket string) *S3UploadApi {

	return &S3UploadApi{
		store:             s,
		cfg:               s.cfg,
		Tusd:              tusdStore.New(bucket, s.client),
		CompleteUploads:   make(chan HookEvent),
		TerminatedUploads: make(chan HookEvent),
		UploadProgress:    make(chan HookEvent),
		CreatedUploads:    make(chan HookEvent),
	}
}

type S3UploadApi struct {
	cfg               *config.Config
	store             *S3Store
	Tusd              tusdStore.S3Store
	UploadProgress    chan HookEvent
	CreatedUploads    chan HookEvent
	TerminatedUploads chan HookEvent
	CompleteUploads   chan HookEvent
}

// func (s *S3UploadApi) writeChunk(ctx context.Context, r *http.Request, resp HTTPResponse, upload Upload, info FileInfo) (HTTPResponse, error) {
// 	// Get Content-Length if possible

// 	length := r.ContentLength
// 	offset := info.Offset

// 	// Test if this upload fits into the file's size
// 	// if !info.SizeIsDeferred && offset+length > info.Size {
// 	// 	return resp, ErrSizeExceeded
// 	// }

// 	// maxSize := info.Size - offset
// 	// // If the upload's length is deferred and the PATCH request does not contain the Content-Length
// 	// // header (which is allowed if 'Transfer-Encoding: chunked' is used), we still need to set limits for
// 	// // the body size.
// 	// if info.SizeIsDeferred {
// 	// 	if handler.config.MaxSize > 0 {
// 	// 		// Ensure that the upload does not exceed the maximum upload size
// 	// 		maxSize = handler.config.MaxSize - offset
// 	// 	} else {
// 	// 		// If no upload limit is given, we allow arbitrary sizes
// 	// 		maxSize = math.MaxInt64
// 	// 	}
// 	// }
// 	if length > 0 {
// 		maxSize = length
// 	}

// 	// c.log.Info("ChunkWriteStart", "maxSize", maxSize, "offset", offset)

// 	// var bytesWritten int64
// 	var err error
// 	// Prevent a nil pointer dereference when accessing the body which may not be
// 	// available in the case of a malicious request.
// 	if r.Body != nil {
// 		// Limit the data read from the request's body to the allowed maximum. We use
// 		// http.MaxBytesReader instead of io.LimitedReader because it returns an error
// 		// if too much data is provided (handled in bodyReader) and also stops the server
// 		// from reading the remaining request body.
// 		body = newBodyReader(ctx, maxSize)
// 		// c.body.onReadDone = func() {
// 		// 	// Update the read deadline for every successful read operation. This ensures that the request handler
// 		// 	// keeps going while data is transmitted but that dead connections can also time out and be cleaned up.
// 		// 	if err := c.resC.SetReadDeadline(time.Now().Add(handler.config.NetworkTimeout)); err != nil {
// 		// 		c.log.Warn("NetworkTimeoutError", "error", err)
// 		// 	}

// 		// 	// The write deadline is updated accordingly to ensure that we can also write responses.
// 		// 	if err := c.resC.SetWriteDeadline(time.Now().Add(2 * handler.config.NetworkTimeout)); err != nil {
// 		// 		c.log.Warn("NetworkTimeoutError", "error", err)
// 		// 	}
// 		// }

// 		// We use a callback to allow the hook system to cancel an upload. The callback
// 		// cancels the request context causing the request body to be closed with the
// 		// provided error.
// 		// info.stopUpload = func(res HTTPResponse) {
// 		// 	cause := ErrUploadStoppedByServer
// 		// 	cause.HTTPResponse = cause.HTTPResponse.MergeWith(res)
// 		// 	c.cancel(cause)
// 		// }

// 		// if handler.config.NotifyUploadProgress {
// 		// 	handler.sendProgressMessages(c, info)
// 		// }

// 		bytesWritten, err = upload.WriteChunk(ctx, offset, r.Body)

// 		// If we encountered an error while reading the body from the HTTP request, log it, but only include
// 		// it in the response, if the store did not also return an error.
// 		bodyErr := c.body.hasError()
// 		if bodyErr != nil {
// 			c.log.Error("BodyReadError", "error", bodyErr.Error())
// 			if err == nil {
// 				err = bodyErr
// 			}
// 		}

// 		// Terminate the upload if it was stopped, as indicated by the ErrUploadStoppedByServer error.
// 		// terminateUpload := errors.Is(bodyErr, ErrUploadStoppedByServer)
// 		// if terminateUpload && handler.composer.UsesTerminater {
// 		// 	if terminateErr := handler.terminateUpload(c, upload, info); terminateErr != nil {
// 		// 		// We only log this error and not show it to the user since this
// 		// 		// termination error is not relevant to the uploading client
// 		// 		c.log.Error("UploadStopTerminateError", "error", terminateErr.Error())
// 		// 	}
// 		// }
// 	}

// 	// c.log.Info("ChunkWriteComplete", "bytesWritten", bytesWritten)

// 	// Send new offset to client
// 	// newOffset := offset + bytesWritten
// 	// resp.Header["Upload-Offset"] = strconv.FormatInt(newOffset, 10)
// 	// handler.Metrics.incBytesReceived(uint64(bytesWritten))
// 	// info.Offset = newOffset

// 	// We try to finish the upload, even if an error occurred. If we have a previous error,
// 	// we return it and its HTTP response.
// 	finishResp, finishErr := s.finishUploadIfComplete(c, resp, upload, info)
// 	if err != nil {
// 		return resp, err
// 	}

// 	return finishResp, finishErr
// }

// finishUploadIfComplete checks whether an upload is completed (i.e. upload offset
// matches upload size) and if so, it will call the data store's FinishUpload
// function and emit the necessary events for the hooks.
// func (s *S3UploadApi) finishUploadIfComplete(c *context.Context, resp HTTPResponse, upload Upload, info FileInfo) (HTTPResponse, error) {
// 	// If the upload is completed, ...
// 	if !info.SizeIsDeferred && info.Offset == info.Size {
// 		var err error
// 		// ... allow the data storage to finish and cleanup the upload
// 		if err = upload.FinishUpload(c); err != nil {
// 			return resp, err
// 		}

// 		// ... and call pre-finish callback and send post-finish notification.
// 		resp, err = s.emitFinishEvents(resp, info)
// 		if err != nil {
// 			return resp, err
// 		}
// 	}

// 	return resp, nil
// }

// func (s *S3UploadApi) emitFinishEvents(resp HTTPResponse, info FileInfo) (HTTPResponse, error) {
// 	// if handler.config.PreFinishResponseCallback != nil {
// 	// 	resp2, err := handler.config.PreFinishResponseCallback(newHookEvent(c, info))
// 	// 	if err != nil {
// 	// 		return resp, err
// 	// 	}
// 	// 	resp = resp.MergeWith(resp2)
// 	// }

// 	// c.log.Info("UploadFinished", "size", info.Size)
// 	// handler.Metrics.incUploadsFinished()

// 	//if handler.config.NotifyCompleteUploads {
// 	s.CompleteUploads <- HookEvent{}
// 	//}

// 	return resp, nil
// }
