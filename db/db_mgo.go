package db

import (
	"VisitSiteProject/model"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func ConnectToDb() *mgo.Session {
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		fmt.Errorf("Error to connect db %v ", err)
	}
	return session
}
func IsUserAuth() {
	//todo
}
func SaveToUsers(user *model.UserStruct) bool {
	fmt.Println("to userSave db")
	n, _ := ConnectToDb().DB("businesscard").C("users").Find(bson.M{"email": user.Email, "password": user.Password}).Count()
	fmt.Println("to userSave db n=", n)

	if n == 0 {
		fmt.Println("to userSave db in 1 if")

		err := ConnectToDb().DB("businesscard").C("users").Insert(user)
		if err != nil {
			fmt.Errorf("Error to connect db %v ", err)
			return false
		}
		fmt.Println("to userSave db set auth true")
		SetAuthToUsers(user.Email, user.Password)
		return true
	}
	return false
}

func GetUser(id string) (*model.UserStruct, bool) {
	fmt.Println("id from get user id=" + id)
	result := model.UserStruct{}
	users := []model.UserStruct{}
	err := ConnectToDb().DB("businesscard").C("users").Find(bson.M{}).All(&users)
	if err != nil {
		fmt.Errorf("Error to connect db %v ", err)
		return nil, false
	}
	fmt.Println("id from get user &user[0]", users)
	if len(users) != 0 {
		for _, user := range users {
			fmt.Println("id from get user ", user.Id.String())
			if user.Id.String() == id {
				fmt.Println("id from get user 4  ", user)
				result = user
			}
		}
	}
	return &result, true
}
func GetAllFromUsers(fname string) []model.UserStruct { //todo for start page
	// критерий выборки
	query := bson.M{}
	if len(fname) != 0 {
		query = bson.M{
			"fname": fname,
		}
	}
	// объект для сохранения результата
	users := []model.UserStruct{}
	ConnectToDb().DB("businesscard").C("users").Find(query).All(&users)
	return users
}

func SetToUsers(user *model.UserStruct) bool {

	fmt.Println("user in SetToUsers ", user)
	err := ConnectToDb().DB("businesscard").C("users").Update(bson.M{"_id": user.Id},
		bson.M{"%set": bson.M{"firstname": user.Firstname, "lastname": user.Lastname, "email": user.Email, "tel": user.Tel,
			"password": user.Password, "profession": user.Profession, "description": user.Description}})
	if err != nil {
		fmt.Errorf("Error to connect db %v ", err)
		return false
	}
	return true
}
func SetAuthToUsers(email, password string) (*model.UserStruct, bool) {
	user := &model.UserStruct{}
	err := ConnectToDb().DB("businesscard").C("users").Update(bson.M{"email": email},
		bson.M{"$set": bson.M{"auth": true}})
	if err != nil {
		fmt.Errorf("Error to connect db %v ", err)
		return nil, false
	}
	err = ConnectToDb().DB("businesscard").C("users").Find(bson.M{"email": email}).One(&user)
	if err != nil || user.Id == "" || user.Password != password {
		fmt.Errorf("Error to connect db %v ", err)
		return nil, false
	}
	return user, true
}
func SetAuthToFalse(email, password string) bool {
	err := ConnectToDb().DB("businesscard").C("users").Update(bson.M{"email": email},
		bson.M{"$set": bson.M{"auth": false}})
	if err != nil {
		fmt.Errorf("Error to connect db %v ", err)
		return false
	}

	return true
}
func DeleteUser(id string) bool {
	// критерий выборки
	query := bson.M{
		"_id": id,
	}
	err := ConnectToDb().DB("businesscard").C("users").Remove(query)
	if err != nil {
		fmt.Errorf("Error to connect db %v ", err)
		return false
	}
	return true
}
