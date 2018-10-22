package audit

import "errors"

type Deets struct {
	Name     string            `json:"name"`
	Location string            `json:"location"`
	Charts   map[string][]Deet `json:"charts"`
}

type Deet struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Digest  string `json:"digest"`
}

func (d *Deets) Get(name, version string) (Deet, error) {
	for _, dd := range d.Charts[name] {
		if dd.Version == version {
			return dd, nil
		}
	}

	// TODO(mattfarina): Update index for next run
	return Deet{}, errors.New("Cannot find details for chart")
}

func (d *Deets) Add(name string, deet Deet) {
	if len(d.Charts[name]) == 0 {
		d.Charts[name] = []Deet{}
	}

	d.Charts[name] = append(d.Charts[name], deet)
}
