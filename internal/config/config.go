package config

import (
	"encoding/json"
	"fmt"
	"loganalyzer/internal/analyzer"
	"os"
)

func LoadConfig(path string) ([]analyzer.LogTarget, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lecture config: %w", err)
	}
	var targets []analyzer.LogTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("parse config JSON: %w", err)
	}
	return targets, nil
}
