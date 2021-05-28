package dao

import (
	"errors"
	"file_server/global"
	"file_server/model"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)


func UserSignUp(username string,passwd string,email,phone string)bool{
	hashedPassword,err :=bcrypt.GenerateFromPassword([]byte(passwd),bcrypt.DefaultCost)
	if err != nil{
		fmt.Println("get hash password failed")
		return false
	}
	user := model.TblUser{
		UserName: username,
		UserPwd: string(hashedPassword),
		Email: email,
		Phone: phone,
	}
	if global.DBEngine.Create(&user).Error!=nil{
		return false
	}
	return true

}


func UserValidate(phone string)bool {
	var user model.TblUser
	global.DBEngine.Find(&user).Where("phone =?", phone)
	user.EmailValidated = 1
	global.DBEngine.Save(&user)
	return true
}

func UserSignin(phone,password string)bool{
	var user model.TblUser
    global.DBEngine.Where("phone =?",phone).Find(&user)
	if err:=bcrypt.CompareHashAndPassword([]byte(user.UserPwd),[]byte(password));err!=nil{
		//c.JSON(http.StatusUnprocessableEntity,gin.H{
		//	"code":422,
		//	"msg":"密码错误",
		//})f
		return false
	}
	return true
}

func GetUserInfoDB(phone string)(model.TblUser,error){
	var user model.TblUser
	if global.DBEngine.Where("phone=?",phone).Find(&user).Error!=nil{
		return model.TblUser{},errors.New("fail to load UserInfo from DB")
	}
	return user,nil

}

func CheckUserValidate(phone string)bool{
	var user model.TblUser
	global.DBEngine.Where("phone = ?",phone).Find(&user)
	if user.EmailValidated==1 || user.PhoneValidated==1{
		return true
	}else {
		return false
	}

}