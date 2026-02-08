package jumpserver

import "fmt"

// Asset represents a JumpServer asset
type Asset struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"` // Request field name
	Addrs    string   `json:"addrs"`   // Response field name
	Platform Platform `json:"platform"`
	Nodes    []Node   `json:"nodes,omitempty"`
	IsActive bool     `json:"is_active"`
	Comment  string   `json:"comment,omitempty"`
	Created  string   `json:"date_created,omitempty"`
	Updated  string   `json:"date_updated,omitempty"`
}

// PlatformRequest represents platform in request
type PlatformRequest struct {
	PK int `json:"pk"`
}

// Platform represents a platform (Linux, Windows, etc.)
type Platform struct {
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Type        interface{} `json:"type"`     // Can be string or object {"value":"", "label":""}
	Category    interface{} `json:"category"` // Can be string or object {"value":"", "label":""}
}

// GetID returns the platform ID as string
func (p *Platform) GetID() string {
	switch v := p.ID.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

// GetTypeValue returns the type value as string
func (p *Platform) GetTypeValue() string {
	switch v := p.Type.(type) {
	case string:
		return v
	case map[string]interface{}:
		if val, ok := v["value"].(string); ok {
			return val
		}
	}
	return ""
}

// GetCategoryValue returns the category value as string
func (p *Platform) GetCategoryValue() string {
	switch v := p.Category.(type) {
	case string:
		return v
	case map[string]interface{}:
		if val, ok := v["value"].(string); ok {
			return val
		}
	}
	return ""
}

// Protocol represents a protocol configuration
type Protocol struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// AssetAccount represents an asset account in create/update asset request
type AssetAccount struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Secret      string `json:"secret"`
	SecretType  string `json:"secret_type"`
	Privileged  bool   `json:"privileged"`
	PushNow     bool   `json:"push_now"`
	SecretReset bool   `json:"secret_reset"`
	OnInvalid   string `json:"on_invalid"`
	IsActive    bool   `json:"is_active"`
}

// Node represents an organization node
type Node struct {
	ID       string `json:"id"`
	FullName string `json:"full_value"` // Changed from full_name to full_value to match API
	Value    string `json:"value"`
	Name     string `json:"name"`
	Weight   int    `json:"weight"`
}

// CreateAssetRequest defines the request to create an asset
type CreateAssetRequest struct {
	Name      string          `json:"name"`
	Address   string          `json:"address"`
	Platform  PlatformRequest `json:"platform"`
	Nodes     []NodeRequest   `json:"nodes,omitempty"`
	Protocols []Protocol      `json:"protocols,omitempty"`
	Accounts  []AssetAccount  `json:"accounts,omitempty"`
	Labels    []string        `json:"labels,omitempty"`
	IsActive  bool            `json:"is_active"`
	Comment   string          `json:"comment,omitempty"`
}

// NodeRequest represents node in request
type NodeRequest struct {
	PK string `json:"pk"`
}

// UpdateAssetRequest defines the request to update an asset
type UpdateAssetRequest struct {
	Name      string          `json:"name,omitempty"`
	Address   string          `json:"address,omitempty"`
	Platform  PlatformRequest `json:"platform,omitempty"`
	Nodes     []NodeRequest   `json:"nodes,omitempty"`
	Protocols []Protocol      `json:"protocols,omitempty"`
	Accounts  []AssetAccount  `json:"accounts,omitempty"`
	Labels    []string        `json:"labels,omitempty"`
	IsActive  *bool           `json:"is_active,omitempty"`
	Comment   string          `json:"comment,omitempty"`
}

// AssetListResponse represents a paginated list of assets
type AssetListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
	Results  []Asset `json:"results"`
}

// CreateAsset creates a new asset
func (c *Client) CreateAsset(req *CreateAssetRequest) (*Asset, error) {
	var result Asset
	err := c.Post("/api/v1/assets/hosts/", req, &result)
	return &result, err
}

// GetAsset retrieves an asset by ID
func (c *Client) GetAsset(id string) (*Asset, error) {
	var result Asset
	err := c.Get(fmt.Sprintf("/api/v1/assets/hosts/%s/", id), &result)
	return &result, err
}

// ListAssets retrieves a list of assets
func (c *Client) ListAssets() ([]Asset, error) {
	var result AssetListResponse
	err := c.Get("/api/v1/assets/hosts/", &result)
	return result.Results, err
}

// UpdateAsset updates an existing asset
func (c *Client) UpdateAsset(id string, req *UpdateAssetRequest) (*Asset, error) {
	var result Asset
	err := c.Put(fmt.Sprintf("/api/v1/assets/hosts/%s/", id), req, &result)
	return &result, err
}

// DeleteAsset deletes an asset
func (c *Client) DeleteAsset(id string) error {
	return c.Delete(fmt.Sprintf("/api/v1/assets/hosts/%s/", id), nil)
}
