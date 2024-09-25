package model

type AggregatedCommits struct {
	Last_Week  int `json:"last_week"`
	Last_Month int `json:"last_month"`
	Last_Year  int `json:"last_year"`
	All_Time   int `json:"all_time"`
}
