package model

type Repo struct {
	Name  string `json:"name"`
	Owner Owner  `json:"owner"`
}
