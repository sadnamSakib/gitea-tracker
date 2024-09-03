package model

import "time"

// User represents a user in the system
type User struct {
	Id              int       `json:"id"`
	Username        string    `json:"username"`
	Full_name       string    `json:"full_name"`
	Avatar_url      string    `json:"avatar_url"`
	Email           string    `json:"email"`
	Last_updated    time.Time `json:"last_updated"`
	Total_commits   int       `json:"total_commits"`
	Weekly_commits  int       `json:"weekly_commits"`
	Monthly_commits int       `json:"monthly_commits"`
	Yearly_commits  int       `json:"yearly_commits"`
}
