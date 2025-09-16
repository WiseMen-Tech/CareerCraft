package models

type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Education string `json:"education" bson:"education"`
	Skills    string `json:"skills" bson:"skills"`
	Interests string `json:"interests" bson:"interests"`
	Location  string `json:"location" bson:"location"`
	Phone     string `json:"phone" bson:"phone"`
}
