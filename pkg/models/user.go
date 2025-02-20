package models

type User struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email" indexed:"true"`
	Username string `json:"username" bson:"username"`
	ImageURL string `json:"image_url" bson:"image_url"`
	Password string `json:"password" bson:"password"`
}
