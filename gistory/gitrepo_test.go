package gistory

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

type gitCliStub struct {
}

func (cli *gitCliStub) git(args ...string) string {
	result := ""
	switch cmd := args[0]; cmd {
	case "log":
		result = "69ad2092|Wed, 1 Feb 2017 10:12:47 +0100\n1a616051|Tue, 31 Jan 2017 10:43:40 +0100"
	case "show":
		result = ""
	default:
		result = ""
	}
	return result
}

func TestChangesToFile(t *testing.T) {
	t.Parallel()
	tmpDir, _ := ioutil.TempDir("", "")
	os.Mkdir(path.Join(tmpDir, ".git"), 0666)
	cli := new(gitCliStub)
	repo := NewGitRepo(tmpDir, cli)
	commits := repo.ChangesToFile("somefile")

	if commits[0].ShortHash != "69ad2092" {
		t.Errorf("Expected commit hash to be 69ad2092 but it was %s instead.", commits[0].ShortHash)
	}
	if commits[0].Date.Format(time.RFC1123Z) != "Wed, 01 Feb 2017 10:12:47 +0100" {
		t.Errorf("Expected commit date to be Wed, 01 Feb 2017 10:12:47 +0100 but it was %s instead.", commits[0].Date.Format(time.RFC1123Z))
	}
	if commits[1].ShortHash != "1a616051" {
		t.Errorf("Expected commit hash to be n1a616051 but it was %s instead.", commits[1].ShortHash)
	}
	if commits[1].Date.Format(time.RFC1123Z) != "Tue, 31 Jan 2017 10:43:40 +0100" {
		t.Errorf("Expected commit date to be Tue, 31 Jan 2017 10:43:40 +0100 but it was %s instead.", commits[1].Date.Format(time.RFC1123Z))
	}
}
