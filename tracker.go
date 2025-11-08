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
	// First, try to get the active window ID
	cmd := exec.Command("xdotool", "getactivewindow")
	windowIDBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get active window ID: %w", err)
	}
	windowID := strings.TrimSpace(string(windowIDBytes))
	if windowID == "" || windowID == "0" {
		return nil, fmt.Errorf("no active window found")
	}

	// Get window title (this might fail for some windows, so we make it optional)
	title := ""
	cmd = exec.Command("xdotool", "getwindowname", windowID)
	if titleBytes, err := cmd.Output(); err == nil {
		title = strings.TrimSpace(string(titleBytes))
	}

	// Get window PID
	pid := ""
	cmd = exec.Command("xdotool", "getwindowpid", windowID)
	if pidBytes, err := cmd.Output(); err == nil {
		pid = strings.TrimSpace(string(pidBytes))
	}

	// Get process name from PID
	process := ""
	if pid != "" && pid != "0" {
		cmd = exec.Command("ps", "-p", pid, "-o", "comm=")
		if processBytes, err := cmd.Output(); err == nil {
			process = strings.TrimSpace(string(processBytes))
		}
	}

	// If we still don't have a process name, try xprop WM_CLASS as fallback
	// This works better with some window managers like i3
	if process == "" {
		cmd = exec.Command("xprop", "-id", windowID, "WM_CLASS")
		if classBytes, err := cmd.Output(); err == nil {
			classStr := string(classBytes)
			// WM_CLASS format: WM_CLASS(STRING) = "instance", "class"
			// We want the class name (second quoted string)
			if idx := strings.LastIndex(classStr, `"`); idx > 0 {
				if prevIdx := strings.LastIndex(classStr[:idx], `"`); prevIdx > 0 {
					process = strings.ToLower(strings.TrimSpace(classStr[prevIdx+1 : idx]))
				}
			}
		}
	}

	return &WindowInfo{
		Title:   title,
		Process: process,
		PID:     pid,
	}, nil
}
