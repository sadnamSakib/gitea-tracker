package model

// User represents a user in the system
type Heatmap struct {
	Username       string         `json:"username"`
	HeatmapEntries []HeatmapEntry `json:"heatmap_entry"`
}
