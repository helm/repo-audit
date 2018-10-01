package audit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// An example with two repositories
var cs = `[
	{
		"name": "foo",
		"location": "https://example.com/foo"
	},
	{
		"name": "bar",
		"location": "https://bar.example.com/"
	}
]`

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig([]byte(cs))
	if err != nil {
		t.Errorf("unable to load configuration: %s", err)
	}

	if c[0].Name != "foo" {
		t.Error("loading config parsing wrong name")
	}

	if c[1].Location != "https://bar.example.com/" {
		t.Error("loading config parsing wrong location")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Write config to temp file and clean it up after test
	baseDir, err := ioutil.TempDir(os.TempDir(), "helm-repo-audit")
	if err != nil {
		t.Fatalf("unable to create temp directory: %s", err)
	}
	defer os.RemoveAll(baseDir)

	tmpcf := filepath.Join(baseDir, "tpconf")
	if err = ioutil.WriteFile(tmpcf, []byte(cs), 0666); err != nil {
		t.Fatalf("unable to write testing temp file: %s", err)
	}

	c, err := LoadConfigFromFile(tmpcf)
	if err != nil {
		t.Errorf("unable to load configuration: %s", err)
	}

	if c[0].Name != "foo" {
		t.Error("loading config parsing wrong name")
	}

	if c[1].Location != "https://bar.example.com/" {
		t.Error("loading config parsing wrong location")
	}
}
