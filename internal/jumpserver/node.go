package jumpserver

import (
	"encoding/json"
	"fmt"
)

// NodeListResponse represents a paginated list of nodes
type NodeListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
	Results  []Node  `json:"results"`
}

// ListNodes retrieves a list of nodes
func (c *Client) ListNodes() ([]Node, error) {
	// Try to unmarshal as paginated response first
	var paginated NodeListResponse
	resp, err := c.getRaw("/api/v1/assets/nodes/")
	if err != nil {
		return nil, err
	}

	// Try paginated format first
	if err := json.Unmarshal(resp, &paginated); err == nil {
		return paginated.Results, nil
	}

	// If that fails, try direct array format
	var nodes []Node
	if err := json.Unmarshal(resp, &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response as paginated or array: %w", err)
	}

	return nodes, nil
}

// GetNodeByFullName retrieves a node by full name
func (c *Client) GetNodeByFullName(fullName string) (*Node, error) {
	nodes, err := c.ListNodes()
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if n.FullName == fullName {
			return &n, nil
		}
	}

	return nil, fmt.Errorf("node not found: %s", fullName)
}
