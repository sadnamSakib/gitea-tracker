package model

import "time"

// User represents a user in the system
type Activity struct {
	Id          int         `json:"id"`
	PerformedBy PerformedBy `json:"act_user"`
	OpType      string      `json:"op_type"`
	Repo        Repo        `json:"repo"`
	Date        time.Time   `json:"created"`
}
