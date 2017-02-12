package gistory

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestParseVersion(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/Gemfile.lock")
	if err != nil {
		log.Fatal(err)
	}
	gemfileLockContent := string(data)

	testCases := []struct {
		gemName string
		source  string
		version string
	}{
		{"jasmine", "git", "1.3.2"},
		{"json-schema", "rubygems", "2.6.2"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s from %s", tc.gemName, tc.source), func(t *testing.T) {
			t.Parallel()
			changelog := NewChangelog(tc.gemName, nil)
			version, err := changelog.parseVersion(gemfileLockContent)
			if err != nil {
				t.Errorf("Could parse version for gem %s in lockfile", tc.gemName)
			}
			if version != tc.version {
				t.Errorf("Expected version for gem %s is %s but parsed %s instead.", tc.gemName, tc.version, version)
			}
		})
	}
}
