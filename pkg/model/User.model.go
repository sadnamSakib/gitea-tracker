package model

import "time"

// User represents a user in the system
type User struct {
	Id             int       `json:"id"`
	Username       string    `json:"login"`
	FullName       string    `json:"full_name"`
	AvatarURL      string    `json:"avatar_url"`
	Email          string    `json:"email"`
	LastUpdated    time.Time `json:"last_updated"`
	TotalCommits   int       `json:"total_commits"`
	WeeklyCommits  int       `json:"weekly_commits"`
	MonthlyCommits int       `json:"monthly_commits"`
	YearlyCommits  int       `json:"yearly_commits"`
}
