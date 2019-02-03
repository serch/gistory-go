package gistory

import (
	"errors"
	"fmt"
	"log"
	"regexp"
)

const lockfile string = "Gemfile.lock"

type Changelog struct {
	gemName      string
	repo         Repo
	versionRegex *regexp.Regexp
}

func NewChangelog(gemName string, repo Repo) *Changelog {
	cl := new(Changelog)
	cl.gemName = gemName
	cl.repo = repo
	cl.versionRegex = compileRegexForGemVersion(gemName)
	return cl
}

func (cl *Changelog) Changelog() []VersionChange {
	commits := cl.repo.ChangesToFile(lockfile)
	if len(commits) == 0 {
		log.Fatalf("%s not found in git history", lockfile)
	}
	return cl.versionChangesForCommits(commits)
}

func (cl *Changelog) versionChangesForCommits(commits []Commit) []VersionChange {
	previousVersion := ""
	versionChanges := []VersionChange{}

	// no lockfile found or no changes to the lockfile found
	if len(commits) == 0 {
		return versionChanges
	}

	previousCommit := commits[0]
	previousVersion, err := cl.gemVersionAtCommit(lockfile, previousCommit.ShortHash)
	if err != nil {
		// only one change to the lockfile was found and the gem was not there
		return versionChanges
	}

	for _, currentCommit := range commits[1:] {
		currentVersion, err := cl.gemVersionAtCommit(lockfile, currentCommit.ShortHash)

		if err != nil {
			// gem not found any more, it wasn't in the lockfile back then
			break
		}

		if currentVersion != previousVersion {
			versionChange := VersionChange{Version: previousVersion, Commit: previousCommit}
			versionChanges = append(versionChanges, versionChange)
		}
		previousVersion = currentVersion
		previousCommit = currentCommit
	}

	versionChange := VersionChange{Version: previousVersion, Commit: previousCommit}
	versionChanges = append(versionChanges, versionChange)

	return versionChanges
}

func (cl *Changelog) gemVersionAtCommit(lockfile string, commitHash string) (string, error) {
	fileContent := cl.repo.FileContentAtCommit(lockfile, commitHash)
	version, err := cl.parseVersion(fileContent)
	return version, err
}

func (cl *Changelog) parseVersion(fileContent string) (string, error) {
	matched := cl.versionRegex.FindStringSubmatch(fileContent)
	if len(matched) == 0 {
		return "", errors.New("Couldn't find gem in lockfile")
	}
	return matched[1], nil
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
