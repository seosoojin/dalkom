package models

type Idol struct {
	ID        string `json:"id" bson:"_id"`
	StageName string `json:"stage_name" bson:"stage_name"`
	Name      string `json:"name" bson:"name"`
	GroupID   string `json:"group_id" bson:"group_id"`
}
