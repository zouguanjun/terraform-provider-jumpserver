package jumpserver

import "fmt"

// User represents a JumpServer user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Comment  string `json:"comment,omitempty"`
	IsActive bool   `json:"is_active"`
	Created  string `json:"date_created,omitempty"`
	Updated  string `json:"date_updated,omitempty"`
}

// CreateUserRequest defines the request to create a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Comment  string `json:"comment,omitempty"`
	IsActive bool   `json:"is_active"`
}

// UpdateUserRequest defines the request to update a user
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Comment  string `json:"comment,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
	Results  []User  `json:"results"`
}

// CreateUser creates a new user
func (c *Client) CreateUser(req *CreateUserRequest) (*User, error) {
	var result User
	err := c.Post("/api/v1/users/users/", req, &result)
	return &result, err
}

// GetUser retrieves a user by ID
func (c *Client) GetUser(id string) (*User, error) {
	var result User
	err := c.Get(fmt.Sprintf("/api/v1/users/users/%s/", id), &result)
	return &result, err
}

// GetUserByUsername retrieves a user by username
func (c *Client) GetUserByUsername(username string) (*User, error) {
	users, err := c.ListUsers()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", username)
}

// ListUsers retrieves a list of users
func (c *Client) ListUsers() ([]User, error) {
	var result UserListResponse
	err := c.Get("/api/v1/users/users/", &result)
	return result.Results, err
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(id string, req *UpdateUserRequest) (*User, error) {
	var result User
	err := c.Put(fmt.Sprintf("/api/v1/users/users/%s/", id), req, &result)
	return &result, err
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(id string) error {
	return c.Delete(fmt.Sprintf("/api/v1/users/users/%s/", id), nil)
}
