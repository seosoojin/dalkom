package models

type Card struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ImageUrl     string   `json:"image_url"`
	GroupID      string   `json:"group_id"`
	CollectionID string   `json:"collection_id"`
	IdolIDs      []string `json:"idol_ids"`
}
