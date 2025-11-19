package model

import (
	"bytes"
	"completion-agent/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SangforCompletion struct {
	cfg *config.ModelConfig
}

func NewSangforCompletion(c *config.ModelConfig) LLM {
	return &SangforCompletion{
		cfg: c,
	}
}

func (m *SangforCompletion) Config() *config.ModelConfig {
	return m.cfg
}

func (m *SangforCompletion) Completions(ctx context.Context, p *CompletionParameter) (*CompletionResponse, CompletionStatus, error) {
	// 将data转换为JSON
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil, StatusServerError, err
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", m.cfg.CompletionsUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, StatusReqError, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", m.cfg.Authorization)

	// 发送请求
	client := &http.Client{
		Timeout: m.cfg.Timeout.Duration(),
	}
	resp, err := client.Do(req)
	if err != nil {
		status := StatusServerError
		switch err {
		case context.Canceled:
			status = StatusCanceled
		case context.DeadlineExceeded:
			status = StatusTimeout
		}
		return nil, status, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, StatusServerError, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, StatusModelError, fmt.Errorf("Invalid StatusCode(%d)", resp.StatusCode)
	}
	var rsp CompletionResponse
	if err := json.Unmarshal(body, &rsp); err != nil {
		return nil, StatusServerError, err
	}
	return &rsp, StatusSuccess, nil
}
