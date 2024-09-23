package model

import "time"

type Repo struct {
	Name    string    `json:"name"`
	Owner   Owner     `json:"owner"`
	Created time.Time `json:"created_at"`
}
