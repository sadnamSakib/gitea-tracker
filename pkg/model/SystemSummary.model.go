package model

import "time"

type SystemSummary struct {
	Username              string    `json:"username"`
	Total_users           int       `json:"total_users"`
	Total_repos           int       `json:"total_repos"`
	Total_orgs            int       `json:"total_orgs"`
	Last_synced           time.Time `json:"last_synced"`
	SyncProgress          int       `json:"sync_progress"`
	Is_synced             bool      `json:"is_synced"`
	Estimated_update_time int64     `json:"estimated_update_time"`
}
