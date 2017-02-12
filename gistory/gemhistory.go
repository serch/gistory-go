package gistory

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/serch/gistory-go/utils"
)

func Run(gemName string, repoPath string) {
	checkIfLockfileExistsOrExit(repoPath)
	repo := NewGitRepo(repoPath, nil)
	changelog := NewChangelog(gemName, repo)

	changes := changelog.Changelog()
	if len(changes) == 0 {
		log.Fatalf("Gem '%s' not found in lock file, maybe a typo?\n", gemName)
	}

	fmt.Printf("Gem: %s\n", gemName)
	fmt.Printf("Current version: %s\n\n", gemName)

	fmt.Println("Change history:")
	for _, change := range changes {
		prettyDate := change.Commit.Date.Format(time.RFC1123Z)
		fmt.Printf("%s on %s (commit %s)\n", change.Version, prettyDate, change.Commit.ShortHash)
	}
}

func checkIfLockfileExistsOrExit(repoPath string) {
	lockfilePath := path.Join(repoPath, Lockfile)
	exists, err := utils.FileOrDirExists(lockfilePath)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Fatalf("%s not found in current directory\n", Lockfile)
	}
}
