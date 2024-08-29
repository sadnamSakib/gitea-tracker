package model

// ChatMessage represents a chat message
type Org struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	FullName  string `json:"full_name"`
}
