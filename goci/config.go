package main

type StepConfig struct {
	Type       string   `json:"type"` // "step", "exception", "timeout"
	Name       string   `json:"name"`
	Exe        string   `json:"exe"`
	Args       []string `json:"args"`
	Message    string   `json:"message"`
	TimeoutSec int      `json:"timeout_sec"`
}
