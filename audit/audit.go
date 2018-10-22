package audit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/mattfarina/helm-repo-audit/report"
	"github.com/mattfarina/helm-repo-audit/report/stderr"
	"k8s.io/helm/pkg/repo"
)

// Audit performs an audit of repos
// TODO(mattfarina): Make the process parallel to audit
func Audit(cfgs Configs, store string) error {

	// Where we will store elements we download durning the process that can
	// be deleted after the fact.
	// TODO(mattfarina): make this configurable
	baseDir, err := ioutil.TempDir(os.TempDir(), "helm-repo-audit")
	if err != nil {
		return err
	}
	defer os.RemoveAll(baseDir)

	// Ensure the store exists
	os.MkdirAll(store, 0755)

	// Iterate over each of the repos and audit it
	for _, cfg := range cfgs {

		// Handle the reporter
		// TODO(mattfarina): Make this configurable and handle multiple
		reporter := stderr.New()

		// Making name safe to write to the file system
		sname := sanitize.BaseName(cfg.Name)
		spath := filepath.Join(store, sname+".json")
		icache := filepath.Join(baseDir, sname+"-index.yaml")

		// Fetch index and cache it
		err := fetchIndex(cfg.Location, icache)
		if err != nil {
			return err
		}
		rf, err := repo.LoadIndexFile(icache)
		if err != nil {
			return err
		}

		// Load the local details. If we don't have them than build them for
		// the first time.
		deets, err := ioutil.ReadFile(spath)
		if err != nil && strings.Contains(err.Error(), "no such file or directory") {
			// First time

			// TODO(mattfarina): Make the io writing alterable
			reporter.Add(fmt.Sprintf("First time fetching %s and %q\n", cfg.Name, cfg.Location))
			reporter.Add(fmt.Sprintln("On first run information is downloaded and cached. Future runs will look for changes"))

			err = repoToDeets(rf, spath, cfg.Name, cfg.Location)
			continue
		} else if err != nil {
			return err
		}

		ds := Deets{}
		err = json.Unmarshal(deets, &ds)
		if err != nil {
			return err
		}

		// Do some auditing
		changed, err := compareDigest(&ds, rf, reporter)
		if err != nil {
			return err
		}

		if changed {
			err = saveDeets(spath, ds)
			if err != nil {
				return err
			}
		}

		reporter.Send()
	}

	return nil
}

func fetchIndex(loc, icache string) error {
	// Some repo urls have query strings (for stuff like access controls)
	// Need to handle that
	u, err := url.Parse(loc)
	if err != nil {
		return err
	}
	u.Path = strings.TrimSuffix(u.Path, "/") + "/index.yaml"
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(icache, body, 0666)

	return err
}

func repoToDeets(rf *repo.IndexFile, pth, name, loc string) error {
	deets := Deets{
		Name:     name,
		Location: loc,
		Charts:   make(map[string][]Deet),
	}

	for n, chart := range rf.Entries {
		deets.Charts[n] = []Deet{}
		for _, cv := range chart {
			d := Deet{
				Name:    cv.Name,
				Version: cv.Version,
				Digest:  cv.Digest,
			}
			deets.Charts[n] = append(deets.Charts[n], d)
		}
	}

	return saveDeets(pth, deets)
}

func saveDeets(pth string, deets Deets) error {
	jd, err := json.Marshal(deets)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pth, jd, 0655)
}

func compareDigest(deets *Deets, index *repo.IndexFile, reporter report.Report) (bool, error) {
	changed := false
	fmt.Printf("Auditing %s at %q\n", deets.Name, deets.Location)
	for n, chart := range index.Entries {
		for _, cv := range chart {
			d, err := deets.Get(n, cv.Version)
			if err != nil {

				// Cannot find details for a chart at that version so we need
				// to store those details.
				if err.Error() == "Cannot find details for chart" {
					d := Deet{
						Name:    n,
						Version: cv.Version,
						Digest:  cv.Digest,
					}
					deets.Add(n, d)
					changed = true
				} else {
					reporter.Add(fmt.Sprintf("Error handling chart %q: %s\n", n, err))
				}
			} else {
				if cv.Digest != d.Digest {
					reporter.Add(fmt.Sprintf("PROBLEM: The digest for %q at version %q has changed from %q to %q\n", d.Name, d.Version, d.Digest, cv.Digest))
				}
			}
		}
	}

	return changed, nil
}
