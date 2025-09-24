package analyzer

import "time"

const (
	StatusOK     = "OK"
	StatusFailed = "FAILED"
)

type LogTarget struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type ReportEntry struct {
	LogID        string    `json:"log_id"`
	FilePath     string    `json:"file_path"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	ErrorDetails string    `json:"error_details"`
	Err          error     `json:"-"`
	AnalyzedAt   time.Time `json:"-"`
}
