package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

type ActivityRecord struct {
	WindowKey string        `json:"window_key"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Language  string        `json:"language,omitempty"`
	Program   string        `json:"program,omitempty"`
}

type AppState struct {
	Activities []ActivityRecord `json:"activities"`
	TotalTime  time.Duration    `json:"total_time"`
	LastSave   time.Time        `json:"last_save"`
}

func NewAppState() *AppState {
	return &AppState{
		Activities: make([]ActivityRecord, 0),
		TotalTime:  0,
		LastSave:   time.Now(),
	}
}

func (as *AppState) RecordActivity(windowKey string, duration time.Duration) {
	record := ActivityRecord{
		WindowKey: windowKey,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	as.Activities = append(as.Activities, record)
	as.TotalTime += duration
}

func (as *AppState) Save() error {
	as.LastSave = time.Now()

	stateDir, err := getStateDir()
	if err != nil {
		return fmt.Errorf("failed to get state directory: %w", err)
	}

	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	stateFile := filepath.Join(stateDir, "state.json")

	// Convert durations to strings for JSON
	type StateJSON struct {
		Activities []struct {
			WindowKey string `json:"window_key"`
			Duration  string `json:"duration"`
			Timestamp string `json:"timestamp"`
			Language  string `json:"language,omitempty"`
			Program   string `json:"program,omitempty"`
		} `json:"activities"`
		TotalTime string `json:"total_time"`
		LastSave  string `json:"last_save"`
	}

	stateJSON := StateJSON{
		Activities: make([]struct {
			WindowKey string `json:"window_key"`
			Duration  string `json:"duration"`
			Timestamp string `json:"timestamp"`
			Language  string `json:"language,omitempty"`
			Program   string `json:"program,omitempty"`
		}, len(as.Activities)),
		TotalTime: as.TotalTime.String(),
		LastSave:  as.LastSave.Format(time.RFC3339),
	}

	for i, act := range as.Activities {
		stateJSON.Activities[i] = struct {
			WindowKey string `json:"window_key"`
			Duration  string `json:"duration"`
			Timestamp string `json:"timestamp"`
			Language  string `json:"language,omitempty"`
			Program   string `json:"program,omitempty"`
		}{
			WindowKey: act.WindowKey,
			Duration:  act.Duration.String(),
			Timestamp: act.Timestamp.Format(time.RFC3339),
			Language:  act.Language,
			Program:   act.Program,
		}
	}

	data, err := json.MarshalIndent(stateJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func LoadState() (*AppState, error) {
	stateDir, err := getStateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	stateFile := filepath.Join(stateDir, "state.json")
	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return NewAppState(), nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	type StateJSON struct {
		Activities []struct {
			WindowKey string `json:"window_key"`
			Duration  string `json:"duration"`
			Timestamp string `json:"timestamp"`
			Language  string `json:"language,omitempty"`
			Program   string `json:"program,omitempty"`
		} `json:"activities"`
		TotalTime string `json:"total_time"`
		LastSave  string `json:"last_save"`
	}

	var stateJSON StateJSON
	if err := json.Unmarshal(data, &stateJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	state := &AppState{
		Activities: make([]ActivityRecord, len(stateJSON.Activities)),
	}

	for i, act := range stateJSON.Activities {
		duration, err := time.ParseDuration(act.Duration)
		if err != nil {
			duration = 0
		}
		timestamp, err := time.Parse(time.RFC3339, act.Timestamp)
		if err != nil {
			timestamp = time.Now()
		}

		state.Activities[i] = ActivityRecord{
			WindowKey: act.WindowKey,
			Duration:  duration,
			Timestamp: timestamp,
			Language:  act.Language,
			Program:   act.Program,
		}
		state.TotalTime += duration
	}

	if stateJSON.LastSave != "" {
		if lastSave, err := time.Parse(time.RFC3339, stateJSON.LastSave); err == nil {
			state.LastSave = lastSave
		}
	}

	return state, nil
}

func getStateDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".config", "emotional-support"), nil
}
