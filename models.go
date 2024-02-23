package main

type Policy struct {
	ID          int    `json:"id"`
	PolicyName  string `json:"policyName"`
	Author      string `json:"author"`
	ControlData string `json:"controlData"`
}
