package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	// Get database path
	dbDir, err := getStateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "activity.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	if err := database.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

func (d *Database) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS window_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		window_key TEXT NOT NULL,
		program TEXT,
		window_title TEXT,
		process_name TEXT,
		pid TEXT,
		language TEXT,
		is_programming INTEGER DEFAULT 0,
		started_at TIMESTAMP NOT NULL,
		ended_at TIMESTAMP,
		duration_seconds INTEGER,
		project_path TEXT
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		notification_type TEXT NOT NULL,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		program TEXT,
		language TEXT,
		duration_seconds INTEGER,
		sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS window_checks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		window_key TEXT,
		program TEXT,
		window_title TEXT,
		process_name TEXT,
		pid TEXT,
		checked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_window_sessions_started ON window_sessions(started_at);
	CREATE INDEX IF NOT EXISTS idx_window_sessions_program ON window_sessions(program);
	CREATE INDEX IF NOT EXISTS idx_notifications_sent_at ON notifications(sent_at);
	CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(notification_type);
	CREATE INDEX IF NOT EXISTS idx_window_checks_checked_at ON window_checks(checked_at);
	`

	_, err := d.db.Exec(schema)
	return err
}

func (d *Database) LogWindowSession(session *WindowSession) error {
	query := `
		INSERT INTO window_sessions (
			window_key, program, window_title, process_name, pid,
			language, is_programming, started_at, ended_at,
			duration_seconds, project_path
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var endedAt interface{}
	if !session.EndedAt.IsZero() {
		endedAt = session.EndedAt
	}

	var durationSeconds interface{}
	if session.Duration > 0 {
		durationSeconds = int(session.Duration.Seconds())
	}

	_, err := d.db.Exec(query,
		session.WindowKey,
		session.Program,
		session.WindowTitle,
		session.ProcessName,
		session.PID,
		session.Language,
		session.IsProgramming,
		session.StartedAt,
		endedAt,
		durationSeconds,
		session.ProjectPath,
	)

	return err
}

func (d *Database) LogNotification(notif *NotificationLog) error {
	query := `
		INSERT INTO notifications (
			notification_type, title, message, program, language, duration_seconds
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	var durationSeconds interface{}
	if notif.DurationSeconds > 0 {
		durationSeconds = notif.DurationSeconds
	}

	_, err := d.db.Exec(query,
		notif.Type,
		notif.Title,
		notif.Message,
		notif.Program,
		notif.Language,
		durationSeconds,
	)

	return err
}

func (d *Database) LogWindowCheck(check *WindowCheck) error {
	query := `
		INSERT INTO window_checks (
			window_key, program, window_title, process_name, pid
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		check.WindowKey,
		check.Program,
		check.WindowTitle,
		check.ProcessName,
		check.PID,
	)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

// Data structures for logging
type WindowSession struct {
	WindowKey      string
	Program        string
	WindowTitle    string
	ProcessName    string
	PID            string
	Language       string
	IsProgramming  bool
	StartedAt      time.Time
	EndedAt        time.Time
	Duration       time.Duration
	ProjectPath    string
}

type NotificationLog struct {
	Type            string
	Title           string
	Message         string
	Program         string
	Language        string
	DurationSeconds int
}

type WindowCheck struct {
	WindowKey   string
	Program     string
	WindowTitle string
	ProcessName string
	PID         string
}

