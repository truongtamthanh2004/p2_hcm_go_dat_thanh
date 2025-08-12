package repository

import (
	"auth-service/internal/constant"
	"auth-service/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type UserClient interface {
	CreateUser(ctx context.Context, email, name, role string) (*dto.CreateUserResponse, error)
}

type userClient struct {
	baseURL string
	http    *http.Client
}

func NewUserClient(baseURL string) UserClient {
	return &userClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *userClient) CreateUser(ctx context.Context, email, name, role string) (*dto.CreateUserResponse, error) {
	// Prepare request body
	reqBody := dto.CreateUserRequest{
		Email: email,
		Name:  name,
		Role:  role,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New(constant.ErrMarshalRequest)
	}

	// Create HTTP request
	fullURL, err := url.JoinPath(c.baseURL, constant.CreateUserUrl)
	if err != nil {
		return nil, errors.New(constant.ErrCreateHTTPRequest)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.New(constant.ErrCreateHTTPRequest)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.New(constant.ErrSendHTTPRequest)
	}
	defer resp.Body.Close()

	// Handle non-201 status codes
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(constant.ErrInternalServer)
	}

	// Parse response body
	var wrapper struct {
		Message string                 `json:"message"`
		Data    dto.CreateUserResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, errors.New(constant.ErrUnmarshalResponse)
	}

	return &wrapper.Data, nil
}
