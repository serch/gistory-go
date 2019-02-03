package gistory

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"testing"
	"time"
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

type repoStub struct {
	Repo
	hashToVersion map[string]string
}

func NewRepoStub(hashToVersion map[string]string) *repoStub {
	rs := new(repoStub)
	rs.hashToVersion = hashToVersion
	return rs
}

func (repo *repoStub) FileContentAtCommit(filename string, commitHash string) string {
	version, present := repo.hashToVersion[commitHash]
	if !present {
		version = ""
	}
	formattedVersion := fmt.Sprintf("\n    %s\n", version)
	return formattedVersion
}

func (repo *repoStub) ChangesToFile(filename string) []Commit {
	var commits []Commit
	return commits
}

func TestVersionChangesForCommits_NoChangesToTheLockfileAreFound(t *testing.T) {
	// when no changes to the lockfile are found
	changelog := NewChangelog("my-gem", nil)
	var commits []Commit
	versionChanges := changelog.versionChangesForCommits(commits)
	if len(versionChanges) != 0 {
		t.Errorf("Expected no version changes but found %d instead.", len(versionChanges))
	}
}

func TestVersionChangesForCommits_GemNotFoundInTheLockfile(t *testing.T) {
	// when the gem is not found in the lockfile
	var hashToVersion = map[string]string{
		"1234567": "another-gem (1.2.3)",
	}
	repo := NewRepoStub(hashToVersion)
	changelog := NewChangelog("my-gem", repo)

	var commits [1]Commit
	commits[0] = Commit{ShortHash: "1234567", Date: time.Now()}

	versionChanges := changelog.versionChangesForCommits(commits[:])
	if len(versionChanges) != 0 {
		t.Errorf("Expected no version changes but found %d instead.", len(versionChanges))
	}
}

func TestVersionChangesForCommits_MultipleCommits(t *testing.T) {
	// when there are multiple commits
	var hashToVersion = map[string]string{
		"1abcdef": "sidekiq (5.0.5)",
		"2abcdef": "sidekiq (5.0.2)",
		"3abcdef": "sidekiq (5.0.2)",
		"4abcdef": "sidekiq (5.0.1)",
		"5abcdef": "sidekiq (5.0.1)",
		"6abcdef": "sidekiq (5.0.0)",
		"7abcdef": "foobar (1.2.3)",
	}
	repo := NewRepoStub(hashToVersion)
	changelog := NewChangelog("sidekiq", repo)

	hashes := make([]string, 0)
	for hash, _ := range hashToVersion {
		hashes = append(hashes, hash)
	}
	sort.Strings(hashes)
	commits := []Commit{}
	for _, hash := range hashes {
		commits = append(commits, Commit{ShortHash: hash, Date: time.Now()})
	}

	versionChanges := changelog.versionChangesForCommits(commits[:])
	if len(versionChanges) != 4 {
		t.Errorf("Expected 4 version changes but found %d instead.", len(versionChanges))
	}

	if versionChanges[0].Commit.ShortHash != "1abcdef" {
		t.Errorf("Expected commit hash to be 1abcdef but it was %s instead.", versionChanges[0].Commit.ShortHash)
	}

	if versionChanges[0].Version != "5.0.5" {
		t.Errorf("Expected version to be 5.0.5 but it was %s instead.", versionChanges[0].Version)
	}

	if versionChanges[1].Commit.ShortHash != "3abcdef" {
		t.Errorf("Expected commit hash to be 3abcdef but it was %s instead.", versionChanges[1].Commit.ShortHash)
	}

	if versionChanges[1].Version != "5.0.2" {
		t.Errorf("Expected version to be 5.0.2 but it was %s instead.", versionChanges[1].Version)
	}

	if versionChanges[2].Commit.ShortHash != "5abcdef" {
		t.Errorf("Expected commit hash to be 5abcdef but it was %s instead.", versionChanges[2].Commit.ShortHash)
	}

	if versionChanges[2].Version != "5.0.1" {
		t.Errorf("Expected version to be 5.0.1 but it was %s instead.", versionChanges[2].Version)
	}

	if versionChanges[3].Commit.ShortHash != "6abcdef" {
		t.Errorf("Expected commit hash to be 6abcdef but it was %s instead.", versionChanges[3].Commit.ShortHash)
	}

	if versionChanges[3].Version != "5.0.0" {
		t.Errorf("Expected version to be 5.0.0 but it was %s instead.", versionChanges[3].Version)
	}
}
