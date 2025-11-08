package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type WindowInfo struct {
	Title   string
	Process string
	PID     string
}

type WindowTracker struct{}

func NewWindowTracker() *WindowTracker {
	return &WindowTracker{}
}

func (wt *WindowTracker) GetActiveWindow() (*WindowInfo, error) {
	// Get active window using xdotool
	cmd := exec.Command("xdotool", "getactivewindow", "getwindowname")
	titleBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get window title: %w", err)
	}
	title := strings.TrimSpace(string(titleBytes))

	// Get window PID
	cmd = exec.Command("xdotool", "getactivewindow", "getwindowpid")
	pidBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get window PID: %w", err)
	}
	pid := strings.TrimSpace(string(pidBytes))

	// Get process name from PID
	process := ""
	if pid != "" {
		cmd = exec.Command("ps", "-p", pid, "-o", "comm=")
		processBytes, err := cmd.Output()
		if err == nil {
			process = strings.TrimSpace(string(processBytes))
		}
	}

	return &WindowInfo{
		Title:   title,
		Process: process,
		PID:     pid,
	}, nil
}
