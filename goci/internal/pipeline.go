package internal

import (
	"encoding/json"
	"fmt"
	"goci/internal/steps"

	"os"
	"strings"
	"time"
)

var StepsConfigPath = "../steps.json"

type Executer interface {
	Execute() (string, error)
}

func LoadPipeline(project, branch string) ([]Executer, error) {
	data, err := os.ReadFile(StepsConfigPath)
	if err != nil {
		return nil, err
	}

	var configs []StepConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}

	var pipeline []Executer

	for _, cfg := range configs {

		// Replace {{BRANCH}} inside args
		for i, a := range cfg.Args {
			cfg.Args[i] = strings.ReplaceAll(a, "{{BRANCH}}", branch)
		}

		switch cfg.Type {
		case "step":
			pipeline = append(pipeline,
				steps.NewStep(cfg.Name, cfg.Exe, cfg.Message, project, cfg.Args),
			)

		case "exception":
			pipeline = append(pipeline,
				steps.NewExceptionStep(cfg.Name, cfg.Exe, cfg.Message, project, cfg.Args),
			)

		case "timeout":
			timeout := time.Duration(cfg.TimeoutSec) * time.Second
			pipeline = append(pipeline,
				steps.NewTimeoutStep(cfg.Name, cfg.Exe, cfg.Message, project, cfg.Args, timeout),
			)

		default:
			return nil, fmt.Errorf("unknown step type: %s", cfg.Type)
		}
	}

	return pipeline, nil
}
