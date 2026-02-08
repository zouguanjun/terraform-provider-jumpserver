package jumpserver

import (
	"encoding/json"
	"fmt"
)

// PlatformListResponse represents a paginated list of platforms
type PlatformListResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next,omitempty"`
	Previous *string    `json:"previous,omitempty"`
	Results  []Platform `json:"results"`
}

// ListPlatforms retrieves a list of platforms
func (c *Client) ListPlatforms() ([]Platform, error) {
	// Try to unmarshal as paginated response first
	var paginated PlatformListResponse
	resp, err := c.getRaw("/api/v1/assets/platforms/")
	if err != nil {
		return nil, err
	}

	// Try paginated format first
	if err := json.Unmarshal(resp, &paginated); err == nil {
		return paginated.Results, nil
	}

	// If that fails, try direct array format
	var platforms []Platform
	if err := json.Unmarshal(resp, &platforms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response as paginated or array: %w", err)
	}

	return platforms, nil
}

// GetPlatformByName retrieves a platform by name
func (c *Client) GetPlatformByName(name string) (*Platform, error) {
	platforms, err := c.ListPlatforms()
	if err != nil {
		return nil, err
	}

	for _, p := range platforms {
		if p.Name == name || p.DisplayName == name {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("platform not found: %s", name)
}
