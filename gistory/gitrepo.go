package gistory

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/serch/gistory-go/utils"
)

const gitCliDateLayout string = "Mon, 2 Jan 2006 15:04:05 -0700"

type gitRepo struct {
	Repo
	path string
	cli  gitClier
}

func NewGitRepo(path string, cli gitClier) *gitRepo {
	checkGitFolderExistsOrExit(path)
	if cli == nil {
		cli = NewGitCli()
	}
	repo := new(gitRepo)
	repo.path = path
	repo.cli = cli
	return repo
}

func (repo *gitRepo) ChangesToFile(filename string) []Commit {
	commitsAndDates := repo.cli.git("log", "--pretty=format:%h|%cD", "--max-count=100", "--follow", filename)
	if commitsAndDates == "" {
		return []Commit{}
	}
	commits := repo.parseCommitsAndDates(strings.Split(commitsAndDates, "\n"))
	return commits
}

func (repo *gitRepo) FileContentAtCommit(filename string, commitHash string) string {
	return repo.cli.git("show", fmt.Sprintf("%s:%s", commitHash, filename))
}

func (repo *gitRepo) parseCommitsAndDates(commitsAndDates []string) []Commit {
	commits := []Commit{}
	for _, commitAndDate := range commitsAndDates {
		split := strings.Split(commitAndDate, "|")
		commitHash := split[0]
		date, err := time.Parse(gitCliDateLayout, split[1])
		if err != nil {
			log.Fatalf("Couldn't parse git commit's date, error: %s", err)
		}
		commits = append(commits, Commit{ShortHash: commitHash, Date: date})
	}
	return commits
}

func checkGitFolderExistsOrExit(repoPath string) {
	gitDir := path.Join(repoPath, ".git")
	exists, err := utils.FileOrDirExists(gitDir)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		fmt.Printf("%s is not a git repository\n", repoPath)
		os.Exit(1)
	}
}
