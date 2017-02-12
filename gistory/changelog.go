package gistory

import (
	"errors"
	"fmt"
	"log"
	"regexp"
)

const lockfile string = "Gemfile.lock"

type changelog struct {
	gemName      string
	repo         Repo
	versionRegex *regexp.Regexp
}

func NewChangelog(gemName string, repo Repo) *changelog {
	cl := new(changelog)
	cl.gemName = gemName
	cl.repo = repo
	cl.versionRegex = compileRegexForGemVersion(gemName)
	return cl
}

func (cl *changelog) Changelog() []VersionChange {
	commits := cl.repo.ChangesToFile(lockfile)
	if len(commits) == 0 {
		log.Fatalf("%s not found in git history", lockfile)
	}
	return cl.versionChangesForCommits(commits)
}

func (cl *changelog) versionChangesForCommits(commits []Commit) []VersionChange {
	previousVersion := ""
	versionChanges := []VersionChange{}

	for _, commit := range commits {
		fileContent := cl.repo.FileContentAtCommit(lockfile, commit.ShortHash)
		version, err := cl.parseVersion(fileContent)

		if err != nil {
			// gem not found any more, it wasn't in the lockfile back then
			break
		}

		if version != previousVersion {
			versionChange := VersionChange{Version: version, Commit: commit}
			versionChanges = append(versionChanges, versionChange)
			previousVersion = version
		}
	}
	return versionChanges
}

func (cl *changelog) parseVersion(fileContent string) (string, error) {
	matched := cl.versionRegex.FindStringSubmatch(fileContent)
	if len(matched) > 0 {
		return matched[1], nil
	} else {
		return "", errors.New("Couldn't find gem in lockfile")
	}
}

func compileRegexForGemVersion(gemName string) *regexp.Regexp {
	// gem version looks like "    byebug (9.0.6)"
	regexString := fmt.Sprintf("\n\\s{4}%s \\((.+)\\)\n", gemName)
	regex, err := regexp.Compile(regexString)
	if err != nil {
		log.Fatal(err)
	}
	return regex
}
