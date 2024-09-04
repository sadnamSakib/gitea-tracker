package model

import "time"

// User represents a user in the system
type User struct {
	Id            int            `json:"id"`
	Username      string         `json:"username"`
	Full_name     string         `json:"full_name"`
	Avatar_url    string         `json:"avatar_url"`
	Email         string         `json:"email"`
	Last_updated  time.Time      `json:"last_updated"`
	Total_commits int            `json:"total_commits"`
	Repos         []Repo         `json:"repos"`
	Heatmap       []HeatmapEntry `json:"heatmap"`
}
