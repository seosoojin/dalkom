package models

type Card struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	GroupID      string   `json:"group_id"`
	CollectionID string   `json:"collection_id"`
	IdolIDs      []string `json:"idol_ids"`
}
