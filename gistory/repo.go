package gistory

// Repo interface for all SVC repositories
type Repo interface {
	ChangesToFile(filename string) []Commit
	FileContentAtCommit(filename string, commitId string) string
}
