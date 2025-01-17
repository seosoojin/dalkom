package models

type Binder struct {
	ID          string   `json:"id"`
	Image       string   `json:"image"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UserID      string   `json:"user_id"`
	Type        string   `json:"type"`
	CardIDs     []string `json:"card_ids"`
}
