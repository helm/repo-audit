package email

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

// Report is a reporter that prints the output to email
type Report struct {
	msgs []string
	auth smtp.Auth
	to   []string
	from string
	host string
}

// New creates a new instance of Report
// Map fields need to include:
// - username: For SMTP server
// - password: For SMTP server
// - host: For SMTP server
// - port: For SMTP server
// - to: comma separated list of email addresses to send report to
// - from: email address this is from
// TODO(mattfarina): Add error handling if missing values
func New(conf map[string]string) *Report {
	r := &Report{}
	r.from = conf["from"]
	r.auth = smtp.PlainAuth("", conf["username"], conf["password"], conf["host"])
	r.host = conf["host"] + ":" + conf["port"]
	r.to = strings.Split(conf["to"], ",")

	return r
}

// Add appends a message to the list of messages in the report
func (r *Report) Add(ms string) {
	r.msgs = append(r.msgs, ms)
}

type templateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

// RFC 822 message format
const emailTemplate = `From: {{ .From }};
To: {{ .To }};
Subject: {{ .Subject }};

{{ .Body }}

`

// Send prints the report to Stderr
func (r *Report) Send() error {
	t := time.Now()
	td := &templateData{
		From:    r.from,
		To:      strings.Join(r.to, ", "),
		Subject: "Report from Helm Repo Audit at " + t.Format(time.RFC822),
		Body:    strings.Join(r.msgs, ""),
	}

	temp := template.Must(template.New("email").Parse(emailTemplate))
	var body bytes.Buffer
	err := temp.Execute(&body, td)
	if err != nil {
		return err
	}

	err = smtp.SendMail(r.host,
		r.auth,
		r.from,
		r.to,
		body.Bytes())

	return err
}
