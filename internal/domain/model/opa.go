package model

type OpaResult struct {
	Result []struct {
		DecisionID string `json:"decision_id"`
		Path       string `json:"path"`
		Result     bool   `json:"result"`
	} `json:"result"`
}
