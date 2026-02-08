package jumpserver

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config holds the JumpServer client configuration
type Config struct {
	Endpoint           string
	KeyID              string
	KeySecret          string
	OrgID              string
	Timeout            time.Duration
	InsecureSkipVerify bool
}

// Client represents a JumpServer API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new JumpServer API client
func NewClient(config *Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.OrgID == "" {
		config.OrgID = "00000000-0000-0000-0000-000000000000"
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
			},
		},
	}
}

// SignRequest signs an HTTP request with JumpServer HMAC authentication
func (c *Client) SignRequest(req *http.Request) error {
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set("Date", date)

	// Build string to sign
	requestTarget := "(request-target): " + strings.ToLower(req.Method) + " " + req.URL.Path
	stringToSign := requestTarget + "\ndate: " + date

	// Create HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(c.config.KeySecret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Set Authorization header
	authHeader := fmt.Sprintf(`Signature keyId="%s",algorithm="hmac-sha256",headers="(request-target) date",signature="%s"`,
		c.config.KeyID, signature)
	req.Header.Set("Authorization", authHeader)

	// Set organization header
	if c.config.OrgID != "" {
		req.Header.Set("X-JMS-ORG", c.config.OrgID)
	}

	return nil
}

// DoRequest executes an HTTP request with authentication
func (c *Client) DoRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	var jsonData []byte

	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			fmt.Printf("[DEBUG] [CLIENT] Failed to marshal request body: %v\n", err)
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.config.Endpoint + path
	fmt.Printf("[DEBUG] [CLIENT] =========================================\n")
	fmt.Printf("[DEBUG] [CLIENT] REQUEST START\n")
	fmt.Printf("[DEBUG] [CLIENT] Method: %s\n", method)
	fmt.Printf("[DEBUG] [CLIENT] URL: %s\n", url)
	fmt.Printf("[DEBUG] [CLIENT] Path: %s\n", path)
	if body != nil {
		fmt.Printf("[DEBUG] [CLIENT] Request Body:\n%s\n", string(jsonData))
	} else {
		fmt.Printf("[DEBUG] [CLIENT] Request Body: (none)\n")
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		fmt.Printf("[DEBUG] [CLIENT] Failed to create request: %v\n", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if err := c.SignRequest(req); err != nil {
		fmt.Printf("[DEBUG] [CLIENT] Failed to sign request: %v\n", err)
		return fmt.Errorf("failed to sign request: %w", err)
	}

	fmt.Printf("[DEBUG] [CLIENT] Request Headers:\n")
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("[DEBUG] [CLIENT]   %s: %s\n", key, value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("[DEBUG] [CLIENT] Failed to execute request: %v\n", err)
		fmt.Printf("[DEBUG] [CLIENT] =========================================\n")
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[DEBUG] [CLIENT] Failed to read response body: %v\n", err)
		fmt.Printf("[DEBUG] [CLIENT] =========================================\n")
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("[DEBUG] [CLIENT] Response Status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("[DEBUG] [CLIENT] Response Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("[DEBUG] [CLIENT]   %s: %s\n", key, value)
		}
	}
	fmt.Printf("[DEBUG] [CLIENT] Response Body:\n%s\n", string(respBody))
	fmt.Printf("[DEBUG] [CLIENT] RESPONSE END\n")
	fmt.Printf("[DEBUG] [CLIENT] =========================================\n")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			fmt.Printf("[DEBUG] [CLIENT] Failed to unmarshal response: %v\n", err)
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
		fmt.Printf("[DEBUG] [CLIENT] Unmarshaled Result (type %T): %+v\n", result, result)
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(path string, result interface{}) error {
	return c.DoRequest("GET", path, nil, result)
}

// getRaw performs a GET request and returns raw response bytes
func (c *Client) getRaw(path string) ([]byte, error) {
	var result interface{}
	err := c.DoRequest("GET", path, nil, &result)

	// Extract the raw response from DoRequest
	// We need to modify DoRequest to support returning raw bytes, or parse from debug
	// For now, let's make a direct request
	url := c.config.Endpoint + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := c.SignRequest(req); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

// Post performs a POST request
func (c *Client) Post(path string, body, result interface{}) error {
	return c.DoRequest("POST", path, body, result)
}

// Put performs a PUT request
func (c *Client) Put(path string, body, result interface{}) error {
	return c.DoRequest("PUT", path, body, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string, result interface{}) error {
	return c.DoRequest("DELETE", path, nil, result)
}
