package main

import "encoding/json"

type Policy struct {
	ID          int             `json:"id"`
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
}

type CreatePolicyRequest struct {
	PolicyName  string          `json:"policyName"`
	Author      string          `json:"author"`
	ControlData json.RawMessage `json:"controlData"`
}
