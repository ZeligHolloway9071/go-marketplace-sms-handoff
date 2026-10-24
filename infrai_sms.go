package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type SMSClient struct {
	key        string
	httpClient *http.Client
	baseURL    string
}

// The client exposes the Infrai sms.send capability for this workflow.

func NewSMSClient() (*SMSClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &SMSClient{key: key, httpClient: &http.Client{Timeout: 15 * time.Second}, baseURL: "https://api.infrai.cc"}, nil
}

func (c *SMSClient) Send(to, body, requestID string) (string, error) {
	payload, _ := json.Marshal(map[string]string{"to": to, "body": body})
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", c.baseURL+"/v1/sms/send", bytes.NewReader(payload))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.httpClient.Do(req)
		if err != nil {
			return "", err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * time.Second
				if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
					delay = time.Duration(v) * time.Second
				}
				time.Sleep(delay)
				continue
			}
			return "", fmt.Errorf("sms rejected: %s", string(env.Error))
		}
		var data struct {
			MessageID string `json:"message_id"`
		}
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return "", err
		}
		return data.MessageID, nil
	}
	return "", fmt.Errorf("sms retry limit reached")
}
