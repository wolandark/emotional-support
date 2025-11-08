package main

import (
	"fmt"
	"log"
	"time"
)

type EmotionalSupportApp struct {
	tracker       *WindowTracker
	detector      *ContextDetector
	messenger     *MessageGenerator
	notifier      *Notifier
	state         *AppState
	checkInterval time.Duration
}

func NewEmotionalSupportApp() *EmotionalSupportApp {
	state, err := LoadState()
	if err != nil {
		log.Printf("Warning: Could not load state: %v", err)
		state = NewAppState()
	}

	return &EmotionalSupportApp{
		tracker:       NewWindowTracker(),
		detector:      NewContextDetector(),
		messenger:     NewMessageGenerator(),
		notifier:      NewNotifier(),
		state:         state,
		checkInterval: 5 * time.Second, // Check every 5 seconds
	}
}

func (app *EmotionalSupportApp) Run() error {
	log.Println("Starting Emotional Support Activity Tracker...")

	// Send initial welcome message
	if err := app.notifier.Send("Emotional Support", "I'm here to support you! Let's have a great coding session! 💚", ""); err != nil {
		log.Printf("Warning: Could not send welcome notification: %v", err)
	}

	ticker := time.NewTicker(app.checkInterval)
	defer ticker.Stop()

	lastWindow := ""
	lastWindowTime := time.Now()
	lastNotificationTime := make(map[string]time.Time)

	for {
		select {
		case <-ticker.C:
			windowInfo, err := app.tracker.GetActiveWindow()
			if err != nil {
				log.Printf("Error getting active window: %v", err)
				continue
			}

			// Detect context
			context := app.detector.DetectContext(windowInfo)

			// Check if window changed
			windowKey := fmt.Sprintf("%s|%s", context.Program, context.WindowTitle)
			if windowKey != lastWindow {
				// Save time spent in previous window
				if lastWindow != "" {
					duration := time.Since(lastWindowTime)
					app.state.RecordActivity(lastWindow, duration)
					if err := app.state.Save(); err != nil {
						log.Printf("Error saving state: %v", err)
					}
				}

				lastWindow = windowKey
				lastWindowTime = time.Now()
			}

			// Calculate time spent in current window
			currentDuration := time.Since(lastWindowTime)

			// Generate and send notifications based on context and time
			app.checkAndNotify(context, currentDuration, lastNotificationTime)
		}
	}
}

func (app *EmotionalSupportApp) checkAndNotify(context *Context, duration time.Duration, lastNotificationTime map[string]time.Time) {
	now := time.Now()

	// Check for time-based notifications (e.g., every hour of coding)
	if context.IsProgramming {
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		// Notify at 30 minutes, 1 hour, 2 hours, etc.
		if (hours > 0 && minutes == 0) || (hours == 0 && minutes == 30) {
			key := fmt.Sprintf("time_%dh%dm_%s", hours, minutes, context.Program)
			if lastNotif, ok := lastNotificationTime[key]; !ok || now.Sub(lastNotif) > 50*time.Minute {
				message := app.messenger.GetTimeBasedMessage(context, duration)
				if message != "" {
					if err := app.notifier.Send("Emotional Support", message, ""); err != nil {
						log.Printf("Error sending notification: %v", err)
					} else {
						lastNotificationTime[key] = now
					}
				}
			}
		}
	}

	// Check for language-specific encouragement
	if context.Language != "" {
		key := fmt.Sprintf("lang_%s", context.Language)
		if lastNotif, ok := lastNotificationTime[key]; !ok || now.Sub(lastNotif) > 30*time.Minute {
			message := app.messenger.GetLanguageMessage(context.Language)
			if message != "" {
				if err := app.notifier.Send("Emotional Support", message, ""); err != nil {
					log.Printf("Error sending notification: %v", err)
				} else {
					lastNotificationTime[key] = now
				}
			}
		}
	}

	// Check for health reminders (every 20 minutes)
	key := "health_reminder"
	if lastNotif, ok := lastNotificationTime[key]; !ok || now.Sub(lastNotif) > 20*time.Minute {
		message := app.messenger.GetHealthReminder()
		if message != "" {
			if err := app.notifier.Send("Emotional Support", message, ""); err != nil {
				log.Printf("Error sending notification: %v", err)
			} else {
				lastNotificationTime[key] = now
			}
		}
	}
}
