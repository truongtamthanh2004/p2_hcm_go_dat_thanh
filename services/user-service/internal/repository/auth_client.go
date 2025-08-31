package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
	"user-service/internal/constant"
	"user-service/internal/dto"
)

type AuthClient interface {
	UpdateUser(ctx context.Context, req dto.UpdateAuthUserRequest) error
}

type authClient struct {
	baseURL string
	http    *http.Client
}

func NewAuthClient(baseURL string) AuthClient {
	return &authClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *authClient) UpdateUser(ctx context.Context, req dto.UpdateAuthUserRequest) error {
	// Marshal body
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return errors.New(constant.ErrMarshalRequest)
	}

	// Build URL
	fullURL, err := url.JoinPath(c.baseURL, constant.UpdateAuthUserURL)
	if err != nil {
		return errors.New(constant.ErrCreateHTTPRequest)
	}

	// HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.New(constant.ErrCreateHTTPRequest)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return errors.New(constant.ErrSendHTTPRequest)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(constant.ErrInternalServer)
	}

	return nil
}
