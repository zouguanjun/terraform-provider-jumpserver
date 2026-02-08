package jumpserver

import "fmt"

// Account represents a JumpServer account
type Account struct {
	ID           string      `json:"id"`
	Name         string      `json:"username"`
	Asset        interface{} `json:"asset"` // Can be string ID or object with id
	AssetDisplay string      `json:"asset_display,omitempty"`
	SecretType   interface{} `json:"secret_type"` // Can be string or object {"value":"", "label":""}
	Comment      string      `json:"comment,omitempty"`
	Created      string      `json:"date_created,omitempty"`
	Updated      string      `json:"date_updated,omitempty"`
}

// GetSecretTypeValue returns the secret_type value as string
func (a *Account) GetSecretTypeValue() string {
	switch v := a.SecretType.(type) {
	case string:
		return v
	case map[string]interface{}:
		if val, ok := v["value"].(string); ok {
			return val
		}
	}
	return ""
}

// GetAssetID returns the asset ID as string
func (a *Account) GetAssetID() string {
	switch v := a.Asset.(type) {
	case string:
		return v
	case map[string]interface{}:
		if id, ok := v["id"].(string); ok {
			return id
		}
	}
	return ""
}

// CreateAccountRequest defines the request to create an account
type CreateAccountRequest struct {
	Name       string `json:"username"`
	Asset      string `json:"asset"`
	Secret     string `json:"secret"`
	SecretType string `json:"secret_type"`
	Comment    string `json:"comment,omitempty"`
}

// UpdateAccountRequest defines the request to update an account
type UpdateAccountRequest struct {
	Name       string `json:"username,omitempty"`
	Secret     string `json:"secret,omitempty"`
	SecretType string `json:"secret_type,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// AccountListResponse represents a paginated list of accounts
type AccountListResponse struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next,omitempty"`
	Previous *string   `json:"previous,omitempty"`
	Results  []Account `json:"results"`
}

// CreateAccount creates a new account
func (c *Client) CreateAccount(req *CreateAccountRequest) (*Account, error) {
	var result Account
	err := c.Post("/api/v1/accounts/accounts/", req, &result)
	return &result, err
}

// GetAccount retrieves an account by ID
func (c *Client) GetAccount(id string) (*Account, error) {
	var result Account
	err := c.Get(fmt.Sprintf("/api/v1/accounts/accounts/%s/", id), &result)
	return &result, err
}

// ListAccounts retrieves a list of accounts
func (c *Client) ListAccounts(assetID string) ([]Account, error) {
	path := "/api/v1/accounts/accounts/"
	if assetID != "" {
		path = fmt.Sprintf("/api/v1/accounts/accounts/?asset=%s", assetID)
	}

	var result AccountListResponse
	err := c.Get(path, &result)
	return result.Results, err
}

// UpdateAccount updates an existing account
func (c *Client) UpdateAccount(id string, req *UpdateAccountRequest) (*Account, error) {
	var result Account
	err := c.Put(fmt.Sprintf("/api/v1/accounts/accounts/%s/", id), req, &result)
	return &result, err
}

// DeleteAccount deletes an account
func (c *Client) DeleteAccount(id string) error {
	return c.Delete(fmt.Sprintf("/api/v1/accounts/accounts/%s/", id), nil)
}
