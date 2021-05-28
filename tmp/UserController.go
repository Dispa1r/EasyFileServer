package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"video_server/api/common"
	"video_server/api/model"
	"video_server/api/session"
)

func CreateUser(c *gin.Context){
	DB := common.InitDB()
    name := c.PostForm("username")
	password := c.PostForm("password")
	if len(password)<6{
		log.Print("password error\n")
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"code":422,
			"message":"密码不能小于6",
		})
		return
	}
	if len(name)==0{
		log.Print("password error\n")
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"code":422,
			"message":"姓名不能为空",
		})
		return
	}
	hashedPassword,err :=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"code":500,
			"message":"系统错误",
		})
		return
	}
	newUser := model.UserInfo{
		Name:name,
		Password: string(hashedPassword),
	}
	session,err1 :=session.GenerateNewSessionId(name)
	if err1!=nil{
		c.JSON(200,gin.H{
			"code":"500",
			"msg":"get session failed",
		})
	}
	DB.Create(&newUser)
	c.JSON(200,gin.H{
		"code":"200",
		"msg":"create user success",
		"session":session,
	})
}
func Login(c *gin.Context){
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println(username,password)
	c.JSON(200,gin.H{
		"code":"200",
		"msg":"ok",
	})

}