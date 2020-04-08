package server

import (
	"VisitSiteProject/db"
	"VisitSiteProject/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
)

var authorisedId bson.ObjectId
func isAuth() bool{
	if string(authorisedId)!="" {
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

	router.Run(":8080")
}

func homePageHandle(context *gin.Context) {
	id := context.Param("id")
	user, _ := db.GetUser(id)

	if string(authorisedId) != "" {
		context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": isAuth(), "User": user})
	} else {
		context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": isAuth(), "User": user})

	}
}

func logoutHandle(context *gin.Context) { //todo

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
	if string(authorisedId) != "" {
		user, _ := db.GetUser(string(authorisedId))
		context.HTML(http.StatusOK, "EditPage.html", bson.M{"IsAuthorized": isAuth(), "User": user})
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
fmt.Println("user",user)
	if db.SaveToUsers(&user) {
		authorisedId = user.Id
		fmt.Println("authorisedid ",authorisedId)

		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": isAuth(), "Id": authorisedId})
		fmt.Println("update auth ",authorisedId)

		//Путь начального  и конечного файлов
		if _, err := os.Stat("./assets/images/tmp.jpg"); os.IsExist(err) {
			err := os.Rename("./assets/images/tmp.jpg", "./assets/images/"+string(user.Id)+".jpg")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func authHandler(context *gin.Context) {
	email := context.PostForm("email")
	password := context.PostForm("password")

	user, ok := db.SetAuthToUsers(email, password)
	if ok {
		authorisedId = user.Id
		context.HTML(http.StatusOK, "HomePage.html", bson.M{"IsAuthorized": user.Auth, "User": user}) //
	} else {
		context.HTML(http.StatusOK, "notLogin.html", bson.M{"IsAuthorized": user.Auth}) //

	}
}

func loginPageHandle(context *gin.Context) {
	context.HTML(http.StatusOK, "loginForm.html", bson.M{})
}

func startPageHandle(context *gin.Context) {
	if string(authorisedId) != "" {
		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": isAuth(), "Id": authorisedId})
	} else {
		context.HTML(http.StatusOK, "StartPage.html", bson.M{"IsAuthorized": isAuth()})
	}
}
