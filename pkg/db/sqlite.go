package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	instance *sql.DB
	once     sync.Once
)

// LogEntry represents a database log entry
type LogEntry struct {
	ID        int64
	Timestamp string
	Hostname  string
	Severity  string
	Service   string
	Message   string
	Context   string // JSON string of the context map
	Analysis  string // Store the AI analysis
}

// InitDB initializes the SQLite database connection
func InitDB(dbPath string) (*sql.DB, error) {
	var err error
	once.Do(func() {
		instance, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Printf("Error opening database: %v", err)
			return
		}

		// Create logs table if it doesn't exist
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS agent_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			hostname TEXT NOT NULL,
			severity TEXT NOT NULL,
			service TEXT NOT NULL,
			message TEXT NOT NULL,
			context TEXT,
			analysis TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`

		_, err = instance.Exec(createTableSQL)
		if err != nil {
			log.Printf("Error creating table: %v", err)
			return
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return instance, nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return instance
}

// InsertLogEntry inserts a new log entry into the database
func InsertLogEntry(entry LogEntry) error {
	if instance == nil {
		return fmt.Errorf("database not initialized")
	}

	query := `
	INSERT INTO agent_logs (timestamp, hostname, severity, service, message, context, analysis)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := instance.Exec(query,
		entry.Timestamp,
		entry.Hostname,
		entry.Severity,
		entry.Service,
		entry.Message,
		entry.Context,
		entry.Analysis,
	)
	if err != nil {
		return fmt.Errorf("failed to insert log: %v", err)
	}

	return nil
}

// GetLogEntries retrieves log entries with optional filters
func GetLogEntries(limit int, severity string) ([]LogEntry, error) {
	if instance == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `
	SELECT id, timestamp, hostname, severity, service, message, context, analysis
	FROM agent_logs
	WHERE (? = '' OR severity = ?)
	ORDER BY timestamp DESC
	LIMIT ?`

	rows, err := instance.Query(query, severity, severity, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %v", err)
	}
	defer rows.Close()

	var entries []LogEntry
	for rows.Next() {
		var entry LogEntry
		err := rows.Scan(
			&entry.ID,
			&entry.Timestamp,
			&entry.Hostname,
			&entry.Severity,
			&entry.Service,
			&entry.Message,
			&entry.Context,
			&entry.Analysis,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
