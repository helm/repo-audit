package report

// Report is the interface for generating and sending reports
type Report interface {
	// Add a message to a report
	Add(message string)

	// Send the messages
	Send() error
}
