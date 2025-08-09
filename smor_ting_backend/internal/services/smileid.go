package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type SmileIDClient struct {
	http *http.Client
	base string
	pid  string
	key  string
}

func NewSmileIDClient(baseURL, partnerID, apiKey string) *SmileIDClient {
	return &SmileIDClient{http: &http.Client{Timeout: 30 * time.Second}, base: strings.TrimRight(baseURL, "/"), pid: partnerID, key: apiKey}
}

type KYCRequest struct {
	Country   string `json:"country"`
	IDType    string `json:"id_type"`
	IDNumber  string `json:"id_number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type KYCResult struct {
	Status    string `json:"status"`
	Reference string `json:"reference"`
}

func (c *SmileIDClient) SubmitKYC(ctx context.Context, req KYCRequest) (*KYCResult, error) {
	url := c.base + "/v1/kyc"
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Partner-ID", c.pid)
	httpReq.Header.Set("X-API-Key", c.key)
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kyc submit failed: %s", string(b))
	}
	var out KYCResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
