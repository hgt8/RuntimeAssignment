package main

import (
	"encoding/json"
	"time"
)

type BasePolicy struct {
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
}

type Policy struct {
	ID int `json:"id"`
	BasePolicy
}

type FullPolicy struct {
	Policy
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePolicyRequest struct {
	BasePolicy
}

type UpdatePolicyRequest struct {
	BasePolicy
	UpdatedAt time.Time `json:"updated_at"`
}
