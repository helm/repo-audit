package stderr

import (
	"fmt"
	"os"
)

// Report is a reporter that prints the output to StdErr
type Report struct {
	msgs []string
}

// New creates a new instance of Report
func New() *Report {
	return &Report{}
}

// Add appends a message to the list of messages in the report
func (r *Report) Add(ms string) {
	r.msgs = append(r.msgs, ms)
}

// Send prints the report to Stderr
func (r *Report) Send() error {
	for _, r := range r.msgs {
		_, err := fmt.Fprint(os.Stderr, r)
		if err != nil {
			return err
		}
	}

	return nil
}
