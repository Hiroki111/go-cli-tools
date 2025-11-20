package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var stepsConfigPath = "steps.json"

func loadPipeline(project, branch string) ([]executer, error) {
	data, err := os.ReadFile(stepsConfigPath)
	if err != nil {
		return nil, err
	}

	var configs []StepConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}

	var pipeline []executer

	for _, cfg := range configs {

		// Replace {{BRANCH}} inside args
		for i, a := range cfg.Args {
			cfg.Args[i] = strings.ReplaceAll(a, "{{BRANCH}}", branch)
		}

		switch cfg.Type {
		case "step":
			pipeline = append(pipeline,
				newStep(cfg.Name, cfg.Exe, cfg.Message, project, cfg.Args),
			)

		case "exception":
			pipeline = append(pipeline,
				newExceptionStep(cfg.Name, cfg.Exe, cfg.Message, project, cfg.Args),
			)

		case "timeout":
			timeout := time.Duration(cfg.TimeoutSec) * time.Second
			pipeline = append(pipeline,
				newTimeoutStep(cfg.Name, cfg.Exe, cfg.Message, project, cfg.Args, timeout),
			)

		default:
			return nil, fmt.Errorf("unknown step type: %s", cfg.Type)
		}
	}

	return pipeline, nil
}
