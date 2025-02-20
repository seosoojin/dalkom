package models

type Binder struct {
	ID          string   `json:"id" bson:"_id"`
	ImageURL    string   `json:"image_url" bson:"image_url,omitempty"`
	Name        string   `json:"name" bson:"name,omitempty"`
	Description string   `json:"description" bson:"description,omitempty"`
	UserID      string   `json:"user_id" bson:"user_id,omitempty" indexed:"true"`
	IsFavorite  bool     `json:"is_favorite" bson:"is_favorite,omitempty" indexed:"true"`
	Type        string   `json:"type" bson:"type,omitempty"`
	CardIDs     []string `json:"card_ids" bson:"card_ids,omitempty"`
}
