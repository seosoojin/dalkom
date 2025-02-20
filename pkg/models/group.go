package models

type Group struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	ImageURL string `json:"image_url" bson:"image_url"`
}
