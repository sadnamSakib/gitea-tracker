package model

// User represents a user in the system
type HeatmapEntry struct {
	Timestamp     int64 `json:"timestamp"`
	Contributions int   `json:"contributions"`
}
