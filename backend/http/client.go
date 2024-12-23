package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func NewClient(disableSSLVerification bool, timeout time.Duration) *http.Client {

	tr := http.DefaultTransport.(*http.Transport).Clone()

	tr.TLSClientConfig.InsecureSkipVerify = disableSSLVerification

	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
}

func CreateRequestBody(body any) (*bytes.Buffer, error) {
	bBuff := &bytes.Buffer{}
	switch v := body.(type) {
	case io.Reader:

		bBytes, err := io.ReadAll(v)
		if err != nil {

			return nil, err
		}

		bBuff.Write(bBytes)

	default:
		err := json.NewEncoder(bBuff).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize body to json: %w", err)
		}
	}

	return bBuff, nil
}

func DecodeRespBody[T any](value *T, resp *http.Response) error {

	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close response body: %v", err)
			}
		}(resp.Body)
	}

	found := false
	for _, status := range []int{http.StatusOK, http.StatusCreated} {
		if resp.StatusCode == status {
			found = true
			break
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	bodyStr := string(bodyBytes)

	if !found {
		return fmt.Errorf("error: %s", bodyStr)
	}

	if err = json.Unmarshal(bodyBytes, &value); err != nil {
		return err
	}

	return nil
}
