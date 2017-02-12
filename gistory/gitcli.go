package gistory

import (
	"log"
	"os/exec"
)

type gitClier interface {
	git(args ...string) string
}

type gitCli struct {
}

func NewGitCli() *gitCli {
	checkGitCliExistsOrExit()
	return new(gitCli)
}

func (cli *gitCli) git(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func checkGitCliExistsOrExit() {
	_, err := exec.Command("git", "--version").Output()
	if err != nil {
		log.Fatalln("git cli is not available, please install it")
	}
}
