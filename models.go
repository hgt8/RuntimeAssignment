package main

import (
	"encoding/json"
	"time"
)

type Policy struct {
	ID          int             `json:"id"`
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
}

type FullPolicy struct {
	ID          int             `json:"id"`
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreatePolicyRequest struct {
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
}

type UpdatePolicyRequest struct {
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
