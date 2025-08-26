package repository

import (
	"chat-service/internal/constant"
	"chat-service/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UserClient interface {
	GetUserByID(ctx context.Context, userID uint) (*dto.UserResponse, error)
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

func (c *userClient) GetUserByID(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	fullURL, err := url.JoinPath(c.baseURL, fmt.Sprintf(constant.GetUserUrl, userID))
	if err != nil {
		return nil, errors.New(constant.ErrCreateHTTPRequest)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, errors.New(constant.ErrCreateHTTPRequest)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.New(constant.ErrSendHTTPRequest)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New(constant.ErrUserNotFound)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(constant.ErrInternalServer)
	}

	var wrapper struct {
		Message string           `json:"message"`
		Data    dto.UserResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, errors.New(constant.ErrUnmarshalResponse)
	}

	return &wrapper.Data, nil
}
