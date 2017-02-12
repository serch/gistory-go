package gistory

import "time"

type Commit struct {
	ShortHash string
	Date      time.Time
}
