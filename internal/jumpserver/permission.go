package jumpserver

import "fmt"

// Permission represents a JumpServer permission
type Permission struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Users       []interface{} `json:"users"` // Can be string array or object array
	UserGroups  []string      `json:"user_groups"`
	Assets      []interface{} `json:"assets"` // Can be string array or object array
	AssetGroups []string      `json:"asset_groups"`
	Actions     []interface{} `json:"actions"` // Can be string array or object array
	Comment     string        `json:"comment,omitempty"`
	Created     string        `json:"date_created,omitempty"`
	Updated     string        `json:"date_updated,omitempty"`
}

// PermissionItem represents a user or asset in permission
type PermissionItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ActionItem represents an action in permission
type ActionItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GetUserIDs extracts user IDs from users array
func (p *Permission) GetUserIDs() []string {
	var ids []string
	for _, u := range p.Users {
		switch v := u.(type) {
		case string:
			ids = append(ids, v)
		case map[string]interface{}:
			if id, ok := v["id"].(string); ok {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

// GetAssetIDs extracts asset IDs from assets array
func (p *Permission) GetAssetIDs() []string {
	var ids []string
	for _, a := range p.Assets {
		switch v := a.(type) {
		case string:
			ids = append(ids, v)
		case map[string]interface{}:
			if id, ok := v["id"].(string); ok {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

// GetActionValues extracts action values from actions array
func (p *Permission) GetActionValues() []string {
	var values []string
	for _, a := range p.Actions {
		switch v := a.(type) {
		case string:
			values = append(values, v)
		case map[string]interface{}:
			if val, ok := v["value"].(string); ok {
				values = append(values, val)
			}
		}
	}
	return values
}

// CreatePermissionRequest defines the request to create a permission
type CreatePermissionRequest struct {
	Name        string   `json:"name"`
	Users       []string `json:"users,omitempty"`
	UserGroups  []string `json:"user_groups,omitempty"`
	Assets      []string `json:"assets,omitempty"`
	AssetGroups []string `json:"asset_groups,omitempty"`
	Actions     []string `json:"actions"`
	Comment     string   `json:"comment,omitempty"`
}

// UpdatePermissionRequest defines the request to update a permission
type UpdatePermissionRequest struct {
	Name        string   `json:"name,omitempty"`
	Users       []string `json:"users,omitempty"`
	UserGroups  []string `json:"user_groups,omitempty"`
	Assets      []string `json:"assets,omitempty"`
	AssetGroups []string `json:"asset_groups,omitempty"`
	Actions     []string `json:"actions,omitempty"`
	Comment     string   `json:"comment,omitempty"`
}

// PermissionListResponse represents a paginated list of permissions
type PermissionListResponse struct {
	Count    int          `json:"count"`
	Next     *string      `json:"next,omitempty"`
	Previous *string      `json:"previous,omitempty"`
	Results  []Permission `json:"results"`
}

// PermissionCreateResponse represents the response from creating a permission
type PermissionCreateResponse struct {
	ID string `json:"id"`
	Permission
}

// CreatePermission creates a new permission
func (c *Client) CreatePermission(req *CreatePermissionRequest) (*Permission, error) {
	// Try direct Permission response first
	var directResult Permission
	err := c.Post("/api/v1/perms/asset-permissions/", req, &directResult)
	if err == nil && directResult.ID != "" {
		return &directResult, nil
	}

	// Try wrapped response
	var wrappedResult PermissionCreateResponse
	err = c.Post("/api/v1/perms/asset-permissions/", req, &wrappedResult)
	if err != nil {
		return nil, err
	}
	// If ID is returned in separate field, copy it
	if wrappedResult.ID != "" {
		wrappedResult.Permission.ID = wrappedResult.ID
	}
	return &wrappedResult.Permission, err
}

// GetPermission retrieves a permission by ID
func (c *Client) GetPermission(id string) (*Permission, error) {
	var result Permission
	err := c.Get(fmt.Sprintf("/api/v1/perms/asset-permissions/%s/", id), &result)
	return &result, err
}

// ListPermissions retrieves a list of permissions
func (c *Client) ListPermissions() ([]Permission, error) {
	var result PermissionListResponse
	err := c.Get("/api/v1/perms/asset-permissions/", &result)
	return result.Results, err
}

// UpdatePermission updates an existing permission
func (c *Client) UpdatePermission(id string, req *UpdatePermissionRequest) (*Permission, error) {
	var result Permission
	err := c.Put(fmt.Sprintf("/api/v1/perms/asset-permissions/%s/", id), req, &result)
	return &result, err
}

// DeletePermission deletes a permission
func (c *Client) DeletePermission(id string) error {
	return c.Delete(fmt.Sprintf("/api/v1/perms/asset-permissions/%s/", id), nil)
}
