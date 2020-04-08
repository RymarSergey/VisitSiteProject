package db

import (
	"VisitSiteProject/model"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


func ConnectToDb() *mgo.Session {
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		fmt.Errorf("Error to connect db %v ",err)
	}
	return session
}
func SaveToUsers(user *model.UserStruct)  bool{
	fmt.Println("to userSave ")
	if n,_:=ConnectToDb().DB("businesscard").C("users").FindId(user.Id).Count();n==0 {
		fmt.Println("to userSave  in 1 if")

		err := ConnectToDb().DB("businesscard").C("users").Insert(user)
		if err != nil {
			fmt.Errorf("Error to connect db %v ", err)
			return false
		}
		return true
	}
	return false
}
func SaveToComments(comment *model.Comment)  bool{
	// добавляем один объект
	err := ConnectToDb().DB("businesscard").C("comments").Insert(comment)
	if err != nil{
		fmt.Errorf("Error to connect db %v ",err)
		return false
	}
	return true
}
func GetUser(id string) (*model.UserStruct,bool) {
	user:=model.UserStruct{}
	err:=ConnectToDb().DB("businesscard").C("users").Find(bson.M{"_id": bson.ObjectId(id)}).One(&user)
	if err != nil{
		fmt.Errorf("Error to connect db %v ",err)
		return nil,false
	}
	return &user,true
}
func GetAllFromUsers(fname string) []model.UserStruct{//todo for start page
	// критерий выборки
	query:=bson.M{}
	if len(fname)!=0 {
		query = bson.M{
			"fname": fname,
		}
	}
	// объект для сохранения результата
	users :=[] model.UserStruct{}
	ConnectToDb().DB("businesscard").C("users").Find(query).All(&users)
	return users
}
func GetAllFromComments(toId string) []model.Comment{//todo for home page
	// критерий выборки
	query:=bson.M{}
	if len(toId)!=0 {
		query = bson.M{
			"toId": toId,
		}
	}
	// объект для сохранения результата
	comments :=[] model.Comment{}
	err:=ConnectToDb().DB("businesscard").C("comments").Find(query).All(&comments)
	if err != nil {
		fmt.Errorf("Error to connect db %v ",err)
	}
	return comments
}

func SetToUsers(user *model.UserStruct) bool {
	err := ConnectToDb().DB("businesscard").C("users").Update(bson.M{"_id": user.Id},user)
	if err != nil{
		fmt.Errorf("Error to connect db %v ",err)
		return false
	}
	return true
}
func SetAuthToUsers(email,password string) (*model.UserStruct,bool) {
	user :=&model.UserStruct{}
	err:= ConnectToDb().DB("businesscard").C("users").Update(bson.M{"email":email},
	bson.M{"$set":bson.M{"auth":true}})
	if err != nil{
		fmt.Errorf("Error to connect db %v ",err)
		return nil,false
	}
	err = ConnectToDb().DB("businesscard").C("users").Find(bson.M{"email":email}).One(&user)
	if err != nil || user.Id=="" || user.Password!=password{
		fmt.Errorf("Error to connect db %v ",err)
		return nil,false
	}
	return user,true
}
func DeleteUser(id string) bool {
	// критерий выборки
	query := bson.M{
		"_id" : id,
	}
	err :=ConnectToDb().DB("businesscard").C("users").Remove(query)
	if err != nil{
		fmt.Errorf("Error to connect db %v ",err)
		return false
	}
	return true
}
