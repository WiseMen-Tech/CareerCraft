package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
}

type UserProfile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Education string             `bson:"education" json:"education"`
	Skills    []string           `bson:"skills" json:"skills"`
	Interests []string           `bson:"interests" json:"interests"`
	Location  string             `bson:"location" json:"location"`
	Phone     string             `bson:"phone" json:"phone"`
	ResumeURL []string             `bson:"resumes,omitempty" json:"resumes,omitempty"`
}
