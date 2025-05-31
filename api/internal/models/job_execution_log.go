package models

import (
	"time"
)

// JobExecutionLog represents a log entry for a job execution
type JobExecutionLog struct {
	LogID            int            `json:"log_id"`
	JobID            int            `json:"job_id"`
	JobType          string         `json:"job_type"`
	StartTime        time.Time      `json:"start_time"`
	End              time.Time      `json:"end_time"`
	Status           string         `json:"status"`
	RecordsProcessed int            `json:"records_processed"`
	ErrorMessage     string         `json:"error_message"`
	Details          map[string]any `json:"details"`
	CreatedAt        time.Time      `json:"created_at"`
}
