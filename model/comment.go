package model

import "gopkg.in/mgo.v2/bson"

type Comment struct {
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	FromId bson.ObjectId `json:"fromId" bson:"fromId"`
	ToId   bson.ObjectId `json:"toId" bson:"toId"`
	Text   string        `json:"text" bson:"text"`
	Time   int64         `json:"time" bson:"time"`
}
