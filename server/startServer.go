package server

import (
	"VisitSiteProject/db"
	"VisitSiteProject/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	"net/http"
	"os"
)

var authorisedId bson.ObjectId

func isAuth() bool {
	if string(authorisedId) != "" {
		return true
	}
	return false
}
func StartServer() {
	fmt.Println("authorisedId=", authorisedId)
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.LoadHTMLGlob("./templates/*")

	router.GET("/", startPageHandle)
	router.GET("/login", loginPageHandle)
	router.GET("/edit", editPageHandle)
	router.GET("/delete", deleteHandle)
	router.GET("/logout", logoutHandle)
	router.GET("/user/:id", homePageHandle)
	router.POST("/auth", authHandler)
	router.POST("/addUser", addUserHandler)
	router.POST("/loadFile", loadFileHandler)
	router.POST("/updateUser", updateUserHandler)

	router.Run(":8080")
}

func updateUserHandler(context *gin.Context) {
	user := model.UserStruct{}
	user.Id = authorisedId
	user.Firstname = context.PostForm("firstname")
	user.Lastname = context.PostForm("lastname")
	user.Email = context.PostForm("email")
	user.Tel = context.PostForm("tel")
	user.Password = context.PostForm("password")
	user.Description = context.PostForm("description")
	user.Profession = context.PostForm("profession")
	user.Auth = true
	fmt.Println("user: ", user)
	if db.SetToUsers(&user) {
		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": true, "Id": authorisedId})
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")
	}
}

func homePageHandle(context *gin.Context) {
	id := context.Param("id")

	fmt.Printf("id tipe %t", id)
	fmt.Println("id=", id)
	fmt.Println("id=", authorisedId)

	user, _ := db.GetUser(id)
	fmt.Println("user=", user)
	if string(authorisedId) != "" {
		context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": isAuth(), "User": user})
	} else {
		context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": isAuth(), "User": user})

	}
}

func logoutHandle(context *gin.Context) {
	user, _ := db.GetUser(authorisedId.String())
	fmt.Println("user=", user)
	db.SetAuthToFalse(user.Email, user.Password)
	authorisedId = ""
	context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": false})
}

func loadFileHandler(context *gin.Context) {
	// single file
	file, _ := context.FormFile("file")
	fmt.Println(file.Filename)

	// Upload the file to specific dst.
	err := context.SaveUploadedFile(file, "./assets/images/tmp.jpg")
	if err != nil {
		fmt.Errorf("Can't load file : %v ", err)
	}
	if string(authorisedId) != "" {
		user, _ := db.GetUser(string(authorisedId))
		context.HTML(http.StatusOK, "EditPage.html", bson.M{"IsAuthorized": isAuth(), "User": user})
	} else {
		context.HTML(http.StatusOK, "EditPage.html", bson.M{"IsAuthorized": isAuth()})

	}
}

func deleteHandle(context *gin.Context) { //todo delete user

}

func editPageHandle(context *gin.Context) {
	fmt.Println("editPageHandle authorisedId", authorisedId)
	if string(authorisedId) != "" {

		fmt.Println("User from editPageHandle")
		user, _ := db.GetUser(authorisedId.String())
		authorisedId = user.Id
		fmt.Println("User from editPageHandle", user)
		context.HTML(http.StatusOK, "EditPage.html", bson.M{"IsAuthorized": isAuth(), "FirstName": user.Firstname,
			"LastName": user.Lastname, "Profession": user.Profession, "Id": user.Id.String(), "Desc": user.Description,
			"Tel": user.Tel, "Password": user.Password, "Email": user.Email})
	} else {
		context.HTML(http.StatusOK, "EditPage.html", bson.M{"IsAuthorized": isAuth()})

	}

}

func addUserHandler(context *gin.Context) {
	user := model.UserStruct{}
	user.Id = bson.NewObjectId()
	user.Firstname = context.PostForm("firstname")
	user.Lastname = context.PostForm("lastname")
	user.Email = context.PostForm("email")
	user.Tel = context.PostForm("tel")
	user.Password = context.PostForm("password")
	user.Description = context.PostForm("description")
	user.Profession = context.PostForm("profession")
	fmt.Println("user: ", user)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")
	if db.SaveToUsers(&user) {
		authorisedId = user.Id
		fmt.Println("add user authorisedid ", authorisedId)

		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": true, "Id": authorisedId})
		fmt.Println("update auth ", authorisedId)

		//Путь начального  и конечного файлов
		err := os.Rename("/assets/images/tmp.jpg", "/assets/images/"+string(user.Id)+".jpg")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func authHandler(context *gin.Context) {
	email := context.PostForm("email")
	password := context.PostForm("password")

	user, ok := db.SetAuthToUsers(email, password)
	if ok {
		authorisedId = user.Id
		context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": user.Auth, "FirstName": user.Firstname,
			"LastName": user.Lastname, "Profession": user.Profession, "Id": user.Id.String(), "Desc": user.Description}) //
	} else {
		context.HTML(http.StatusOK, "notLogin.html", bson.M{"IsAuthorized": user.Auth}) //

	}
}

func loginPageHandle(context *gin.Context) {
	context.HTML(http.StatusOK, "loginForm.html", bson.M{})
}

func startPageHandle(context *gin.Context) {
	if string(authorisedId) != "" {
		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": true, "Id": authorisedId})
	} else {
		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": false})
		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": false})
	}
}
