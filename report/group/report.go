package group

import (
	"github.com/helm/repo-audit/report"
	"github.com/helm/repo-audit/report/email"
	"github.com/helm/repo-audit/report/stderr"
)

func New(cfgs []report.Config) *GroupReporter {
	g := &GroupReporter{}
	for _, c := range cfgs {
		switch c.Type {
		case "stderr":
			re := stderr.New()
			g.AddReporter(re)
		case "email":
			re := email.New(c.Config)
			g.AddReporter(re)
		}
	}

	// stderr is the default if nothing is declared
	if g.Length() == 0 {
		re := stderr.New()
		g.AddReporter(re)
	}

	return g
}

type GroupReporter struct {
	reporters []report.Reporter
}

func (r *GroupReporter) AddReporter(re report.Reporter) {
	r.reporters = append(r.reporters, re)
}

func (r *GroupReporter) Length() int {
	return len(r.reporters)
}

func (r *GroupReporter) Add(ms string) {
	for _, r := range r.reporters {
		r.Add(ms)
	}
}

func (r *GroupReporter) Send() error {
	for _, r := range r.reporters {
		err := r.Send()

		if err != nil {
			return err
		}
	}

	return nil
}
