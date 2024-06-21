package models

type ChatRoom struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatorID int    `json:"creator_id"`
}
