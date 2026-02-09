package events

import "time"

// Event represents a single event entry
type Event struct {
	ID               string                 `json:"id"`
	EventType        string                 `json:"event_type"`
	OccurredAt       time.Time              `json:"occurred_at"`
	ActorID          *string                `json:"actor_id"`
	ActorType        *string                `json:"actor_type"`
	TargetID         *string                `json:"target_id"`
	TargetType       *string                `json:"target_type"`
	TargetIdentifier *string                `json:"target_identifier"`
	TargetSnapshot   map[string]interface{} `json:"target_snapshot,omitempty"`
	Changes          map[string]interface{} `json:"changes,omitempty"`
}

// EventsResponse is the response from GET /events.json
type EventsResponse []Event

// RequestQueryOptions represents optional query parameters for events requests
type RequestQueryOptions struct {
	Filter   string // Filter by actor type: user, system, or all
	Download bool   // For CSV: force download as attachment (set to true to add download=1)
}
