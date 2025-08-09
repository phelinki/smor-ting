package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Minimal MTN MoMo client scaffold to integrate Collection and Disbursement
// Aligns with MTN documentation: https://momodeveloper.mtn.com/api-documentation

type MomoClient struct {
	httpClient         *http.Client
	baseURL            string
	targetEnv          string
	apiUser            string
	apiKey             string
	subKeyCollection   string
	subKeyDisbursement string
}

func NewMomoClient(baseURL, targetEnv, apiUser, apiKey, subCollection, subDisbursement string) *MomoClient {
	return &MomoClient{
		httpClient:         &http.Client{Timeout: 30 * time.Second},
		baseURL:            strings.TrimRight(baseURL, "/"),
		targetEnv:          targetEnv,
		apiUser:            apiUser,
		apiKey:             apiKey,
		subKeyCollection:   subCollection,
		subKeyDisbursement: subDisbursement,
	}
}

// Online-only guard; caller should prevent offline wallet actions
func (c *MomoClient) EnsureOnline(ctx context.Context) error {
	// Simple ping via HEAD on baseURL; in production, better health endpoint
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, c.baseURL, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("momo not reachable: %w", err)
	}
	_ = resp.Body.Close()
	return nil
}

// Token endpoints
type momoToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *MomoClient) basicAuth() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.apiUser+":"+c.apiKey))
}

func (c *MomoClient) GetCollectionToken(ctx context.Context) (*momoToken, error) {
	url := c.baseURL + "/collection/token/"
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	req.Header.Set("Authorization", c.basicAuth())
	req.Header.Set("Ocp-Apim-Subscription-Key", c.subKeyCollection)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("collection token failed: %s", string(b))
	}
	var t momoToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *MomoClient) GetDisbursementToken(ctx context.Context) (*momoToken, error) {
	url := c.baseURL + "/disbursement/token/"
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	req.Header.Set("Authorization", c.basicAuth())
	req.Header.Set("Ocp-Apim-Subscription-Key", c.subKeyDisbursement)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("disbursement token failed: %s", string(b))
	}
	var t momoToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Collection request-to-pay
type Party struct {
	PartyIdType string `json:"partyIdType"`
	PartyId     string `json:"partyId"`
}
type RequestToPay struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalId   string `json:"externalId"`
	Payer        Party  `json:"payer"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

func (c *MomoClient) RequestToPay(ctx context.Context, body RequestToPay) (string, error) {
	tok, err := c.GetCollectionToken(ctx)
	if err != nil {
		return "", err
	}
	ref := randomUUID()
	url := c.baseURL + "/collection/v1_0/requesttopay"
	data, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(data)))
	req.Header.Set("Authorization", tok.TokenType+" "+tok.AccessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.subKeyCollection)
	req.Header.Set("X-Target-Environment", c.targetEnv)
	req.Header.Set("X-Reference-Id", ref)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("r2p failed: %s", string(b))
	}
	return ref, nil
}

func (c *MomoClient) GetRequestToPayStatus(ctx context.Context, referenceId string) (string, error) {
	tok, err := c.GetCollectionToken(ctx)
	if err != nil {
		return "", err
	}
	url := c.baseURL + "/collection/v1_0/requesttopay/status/" + referenceId
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", tok.TokenType+" "+tok.AccessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.subKeyCollection)
	req.Header.Set("X-Target-Environment", c.targetEnv)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("r2p status failed: %s", string(b))
	}
	var out struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Status, nil
}

// Disbursement transfer
type TransferRequest struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalId   string `json:"externalId"`
	Payee        Party  `json:"payee"`
	PayerMessage string `json:"payerMessage,omitempty"`
	PayeeNote    string `json:"payeeNote,omitempty"`
}

func (c *MomoClient) Transfer(ctx context.Context, body TransferRequest) (string, error) {
	tok, err := c.GetDisbursementToken(ctx)
	if err != nil {
		return "", err
	}
	ref := randomUUID()
	url := c.baseURL + "/disbursement/v1_0/transfer"
	data, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(data)))
	req.Header.Set("Authorization", tok.TokenType+" "+tok.AccessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.subKeyDisbursement)
	req.Header.Set("X-Reference-Id", ref)
	req.Header.Set("X-Target-Environment", c.targetEnv)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transfer failed: %s", string(b))
	}
	return ref, nil
}

func (c *MomoClient) GetTransferStatus(ctx context.Context, referenceId string) (string, error) {
	tok, err := c.GetDisbursementToken(ctx)
	if err != nil {
		return "", err
	}
	url := c.baseURL + "/disbursement/v1_0/transfer/status/" + referenceId
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", tok.TokenType+" "+tok.AccessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.subKeyDisbursement)
	req.Header.Set("X-Target-Environment", c.targetEnv)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transfer status failed: %s", string(b))
	}
	var out struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Status, nil
}

func randomUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
