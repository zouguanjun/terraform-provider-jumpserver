package jumpserver

import (
	"fmt"
)

type CommandExecution struct {
	Asset   string `json:"asset"`
	Account string `json:"account"`
	Command string `json:"command"`
}

type FileUpload struct {
	Asset   string `json:"asset"`
	Account string `json:"account"`
	Command string `json:"command"`
}

type CreateCommandExecutionRequest struct {
	Assets      []string `json:"assets"`
	Nodes       []string `json:"nodes"`
	Module      string   `json:"module"`
	Args        string   `json:"args"`
	Runas       string   `json:"runas"`
	RunasPolicy string   `json:"runas_policy"`
	Instant     bool     `json:"instant"`
	IsPeriodic  bool     `json:"is_periodic"`
	Timeout     int      `json:"timeout"`
}

type CommandExecutionResponse struct {
	Data string `json:"data"`
	End  bool   `json:"end"`
	Mark string `json:"mark"`
}

func (c *Client) CreateExecuteCommand(req *CreateCommandExecutionRequest) (*CommandExecutionResponse, error) {
	var result CommandExecutionResponse
	err := c.Post("/api/v1/ops/jobs/", req, &result)
	return &result, err
}

func (c *Client) GetCommandExecution(id string) (*CommandExecutionResponse, error) {
	var result CommandExecutionResponse
	err := c.Get(fmt.Sprintf("/api/v1/ops/celery/task/00000000-0000-0000-0000-000000000002/task-execution/%s/log/", id), &result)
	return &result, err
}
