package model

import "time"

type Repo struct {
	Name               string            `json:"name"`
	Owner              Owner             `json:"owner"`
	Created            time.Time         `json:"created_at"`
	Updated_at         time.Time         `json:"updated_at"`
	Aggregated_Commits AggregatedCommits `json:"aggregated_commits"`
	Following          bool              `json:"following"`
}
