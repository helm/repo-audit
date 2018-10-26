package report

// Report is the interface for generating and sending reports
type Reporter interface {
	// Add a message to a report
	Add(message string)

	// Send the messages
	Send() error
}

type Config struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config,omitempty"`
}
