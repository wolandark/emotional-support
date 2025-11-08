package main

import (
	"fmt"
	"log"
	"time"
)

// NotificationTiming holds all timing configuration for notifications
type NotificationTiming struct {
	// WindowCheckInterval is how often to check for active window changes
	WindowCheckInterval time.Duration

	// TimeBasedNotifications configures time-based coding notifications
	TimeBasedNotifications struct {
		// Intervals are the time milestones to trigger notifications (e.g., 30min, 1hr, 2hr)
		Intervals []time.Duration
		// Cooldown prevents duplicate notifications within this period
		Cooldown time.Duration
		// MinDuration is the minimum time before showing any time-based message
		MinDuration time.Duration
	}

	// LanguageNotifications configures language-specific encouragement
	LanguageNotifications struct {
		// Cooldown between language-specific notifications
		Cooldown time.Duration
	}

	// HealthReminders configures health and wellness reminders
	HealthReminders struct {
		// Interval between health reminder notifications
		Interval time.Duration
	}
}

// DefaultNotificationTiming returns sensible default timing configuration
func DefaultNotificationTiming() *NotificationTiming {
	nt := &NotificationTiming{
		WindowCheckInterval: 5 * time.Second,
	}

	// Time-based: notify at 30min, 1hr, 2hr, 3hr, etc.
	nt.TimeBasedNotifications.Intervals = []time.Duration{
		30 * time.Minute,
		1 * time.Hour,
		2 * time.Hour,
		3 * time.Hour,
		4 * time.Hour,
	}
	nt.TimeBasedNotifications.Cooldown = 1 * time.Minute
	nt.TimeBasedNotifications.MinDuration = 2 * time.Minute

	// Language notifications: every 30 minutes
	nt.LanguageNotifications.Cooldown = 3 * time.Minute

	// Health reminders: every 20 minutes
	nt.HealthReminders.Interval = 2 * time.Minute

	return nt
}

type EmotionalSupportApp struct {
	tracker   *WindowTracker
	detector  *ContextDetector
	messenger *MessageGenerator
	notifier  *Notifier
	state     *AppState
	timing    *NotificationTiming
	database  *Database
}

func NewEmotionalSupportApp() *EmotionalSupportApp {
	state, err := LoadState()
	if err != nil {
		log.Printf("Warning: Could not load state: %v", err)
		state = NewAppState()
	}

	database, err := NewDatabase()
	if err != nil {
		log.Printf("Warning: Could not initialize database: %v", err)
		database = nil
	}

	return &EmotionalSupportApp{
		tracker:   NewWindowTracker(),
		detector:  NewContextDetector(),
		messenger: NewMessageGenerator(),
		notifier:  NewNotifier(),
		state:     state,
		timing:    DefaultNotificationTiming(),
		database:  database,
	}
}

func (app *EmotionalSupportApp) Run() error {
	log.Println("Starting Emotional Support Activity Tracker...")

	// Send initial welcome message
	welcomeMsg := "I'm here to support you! Let's have a great coding session! 💚"
	if err := app.notifier.Send("Emotional Support", welcomeMsg, ""); err != nil {
		log.Printf("Warning: Could not send welcome notification: %v", err)
	} else if app.database != nil {
		// Log welcome notification
		notif := &NotificationLog{
			Type:    "welcome",
			Title:   "Emotional Support",
			Message: welcomeMsg,
		}
		if err := app.database.LogNotification(notif); err != nil {
			log.Printf("Error logging welcome notification: %v", err)
		}
	}

	// Ensure database is closed on exit
	defer func() {
		if app.database != nil {
			if err := app.database.Close(); err != nil {
				log.Printf("Error closing database: %v", err)
			}
		}
	}()

	ticker := time.NewTicker(app.timing.WindowCheckInterval)
	defer ticker.Stop()

	lastWindow := ""
	lastWindowTime := time.Now()
	lastContext := &Context{}
	lastWindowInfo := &WindowInfo{}
	lastNotificationTime := make(map[string]time.Time)

	for {
		select {
		case <-ticker.C:
			windowInfo, err := app.tracker.GetActiveWindow()
			if err != nil {
				// Log but don't spam - only log occasionally
				// This can happen in i3 when switching windows quickly
				log.Printf("Warning: Could not get active window: %v", err)
				continue
			}

			// Log window check to database
			if app.database != nil {
				check := &WindowCheck{
					WindowKey:   fmt.Sprintf("%s|%s", windowInfo.Process, windowInfo.Title),
					Program:     windowInfo.Process,
					WindowTitle: windowInfo.Title,
					ProcessName: windowInfo.Process,
					PID:         windowInfo.PID,
				}
				if err := app.database.LogWindowCheck(check); err != nil {
					log.Printf("Error logging window check: %v", err)
				}
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

					// Log window session to database (using previous window's context)
					if app.database != nil {
						session := &WindowSession{
							WindowKey:     lastWindow,
							Program:       lastContext.Program,
							WindowTitle:    lastContext.WindowTitle,
							ProcessName:    lastWindowInfo.Process,
							PID:           lastWindowInfo.PID,
							Language:      lastContext.Language,
							IsProgramming: lastContext.IsProgramming,
							StartedAt:      lastWindowTime,
							EndedAt:        time.Now(),
							Duration:       duration,
							ProjectPath:    lastContext.ProjectPath,
						}
						if err := app.database.LogWindowSession(session); err != nil {
							log.Printf("Error logging window session: %v", err)
						}
					}
				}

				// Update tracking variables
				lastWindow = windowKey
				lastWindowTime = time.Now()
				lastContext = context
				lastWindowInfo = windowInfo
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
	timing := app.timing

	// Check for time-based notifications (for all programs, not just programming)
	if context.Program != "" && duration >= timing.TimeBasedNotifications.MinDuration {
		// Check if we've reached any of the configured intervals
		// We check within a small window (2x check interval) to account for timing variations
		checkWindow := 2 * timing.WindowCheckInterval
		for _, interval := range timing.TimeBasedNotifications.Intervals {
			// Check if we're at or just past the interval (within the check window)
			if duration >= interval && duration <= interval+checkWindow {
				key := fmt.Sprintf("time_%s_%s", interval.String(), context.Program)
				if lastNotif, ok := lastNotificationTime[key]; !ok || now.Sub(lastNotif) > timing.TimeBasedNotifications.Cooldown {
					message := app.messenger.GetTimeBasedMessage(context, duration)
					if message != "" {
						if err := app.notifier.Send("Emotional Support", message, ""); err != nil {
							log.Printf("Error sending notification: %v", err)
						} else {
							lastNotificationTime[key] = now
							// Log notification to database
							if app.database != nil {
								notif := &NotificationLog{
									Type:            "time_based",
									Title:           "Emotional Support",
									Message:         message,
									Program:         context.Program,
									Language:        context.Language,
									DurationSeconds: int(duration.Seconds()),
								}
								if err := app.database.LogNotification(notif); err != nil {
									log.Printf("Error logging notification: %v", err)
								}
							}
						}
					}
				}
				break // Only trigger one notification per check
			}
		}
	}

	// Check for language-specific encouragement
	if context.Language != "" {
		key := fmt.Sprintf("lang_%s", context.Language)
		if lastNotif, ok := lastNotificationTime[key]; !ok || now.Sub(lastNotif) > timing.LanguageNotifications.Cooldown {
			message := app.messenger.GetLanguageMessage(context.Language)
			if message != "" {
				if err := app.notifier.Send("Emotional Support", message, ""); err != nil {
					log.Printf("Error sending notification: %v", err)
				} else {
					lastNotificationTime[key] = now
					// Log notification to database
					if app.database != nil {
						notif := &NotificationLog{
							Type:     "language",
							Title:    "Emotional Support",
							Message:  message,
							Program:  context.Program,
							Language: context.Language,
						}
						if err := app.database.LogNotification(notif); err != nil {
							log.Printf("Error logging notification: %v", err)
						}
					}
				}
			}
		}
	}

	// Check for health reminders
	key := "health_reminder"
	if lastNotif, ok := lastNotificationTime[key]; !ok || now.Sub(lastNotif) > timing.HealthReminders.Interval {
		message := app.messenger.GetHealthReminder()
		if message != "" {
			if err := app.notifier.Send("Emotional Support", message, ""); err != nil {
				log.Printf("Error sending notification: %v", err)
			} else {
				lastNotificationTime[key] = now
				// Log notification to database
				if app.database != nil {
					notif := &NotificationLog{
						Type:    "health",
						Title:   "Emotional Support",
						Message: message,
					}
					if err := app.database.LogNotification(notif); err != nil {
						log.Printf("Error logging notification: %v", err)
					}
				}
			}
		}
	}
}
