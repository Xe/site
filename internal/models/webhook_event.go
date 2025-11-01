package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// WebhookEvent represents an incoming GitHub Sponsors webhook event.
// This helps with auditing and debugging webhook processing.
type WebhookEvent struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	GitHubID  string          `json:"github_id" gorm:"uniqueIndex;not null"` // GitHub's delivery ID
	Action    string          `json:"action" gorm:"not null"`                // "created", "edited", "cancelled", etc.
	SenderID  uint            `json:"sender_id"`                             // The account that triggered the event
	Sender    Account         `json:"sender" gorm:"foreignKey:SenderID"`
	SponsorshipID *uint       `json:"sponsorship_id"`                        // The sponsorship this event relates to
	Sponsorship *Sponsorship  `json:"sponsorship" gorm:"foreignKey:SponsorshipID"`

	// Event details
	EventType    string     `json:"event_type" gorm:"not null"` // Always "sponsorship" for our use case
	ProcessedAt  time.Time  `json:"processed_at"`              // When we processed the event
	Success      bool       `json:"success" gorm:"default:true"` // Whether processing was successful
	ErrorMessage string     `json:"error_message"`               // Error message if processing failed

	// Raw webhook payload for debugging/reprocessing
	PayloadJSON string `json:"-" gorm:"type:text"` // JSON encoded webhook payload
	PayloadSize  int    `json:"payload_size"`       // Size of the payload in bytes

	// Request details
	RemoteAddr string    `json:"remote_addr"` // Client IP address
	UserAgent  string    `json:"user_agent"`  // Client user agent
	Timestamp  time.Time `json:"timestamp"`   // When the webhook was received

	// Local timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the WebhookEvent model.
func (WebhookEvent) TableName() string {
	return "webhook_events"
}

// Payload returns the parsed webhook payload.
func (w *WebhookEvent) Payload() map[string]interface{} {
	if w.PayloadJSON == "" {
		return make(map[string]interface{})
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(w.PayloadJSON), &payload); err != nil {
		return make(map[string]interface{})
	}
	return payload
}

// SetPayload sets the webhook payload from a map.
func (w *WebhookEvent) SetPayload(payload map[string]interface{}) error {
	if payload == nil {
		w.PayloadJSON = ""
		w.PayloadSize = 0
		return nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.PayloadJSON = string(data)
	w.PayloadSize = len(data)
	return nil
}

// BeforeSave marshals the payload before saving.
func (w *WebhookEvent) BeforeSave(tx *gorm.DB) error {
	return w.SetPayload(w.Payload())
}

// BeforeCreate ensures data consistency before creating records.
func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if w.GitHubID == "" {
		return gorm.ErrInvalidField
	}
	if w.Action == "" {
		return gorm.ErrInvalidField
	}
	if w.Timestamp.IsZero() {
		w.Timestamp = time.Now()
	}
	if w.ProcessedAt.IsZero() {
		w.ProcessedAt = time.Now()
	}
	return nil
}

// IsRecent returns true if the webhook was received within the last hour.
func (w *WebhookEvent) IsRecent() bool {
	return time.Since(w.Timestamp) < time.Hour
}

// GetProcessingDuration returns how long it took to process the webhook.
func (w *WebhookEvent) GetProcessingDuration() time.Duration {
	return w.ProcessedAt.Sub(w.Timestamp)
}