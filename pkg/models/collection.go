package models

type Collection struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	GroupID string `json:"group_id" bson:"group_id"`
}
