package analytics

import "time"

// Analytic represents a single analytics entry
type Analytic struct {
	Device   string    `json:"device"`
	Command  string    `json:"command"`
	LastRun  time.Time `json:"last_run"`
	Count    int       `json:"count"`
}

// AnalyticsResponse is the response from GET /analytics.json
type AnalyticsResponse []Analytic
