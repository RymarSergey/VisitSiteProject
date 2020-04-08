package model

import (
	"gopkg.in/mgo.v2/bson"
)

type UserStruct struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	Firstname   string        `json:"firstname" bson:"firstname" `
	Lastname    string        `json:"lastname" bson:"lastname" `
	Email       string        `json:"email" bson:"email" `
	Tel         string        `json:"tel" bson:"tel" `
	Password    string        `json:"password" bson:"password" `
	Profession  string        `json:"profession" bson:"profession" `
	Description string        `json:"description" bson:"description"`
	Auth        bool          `json:"auth" bson:"auth"`
}


