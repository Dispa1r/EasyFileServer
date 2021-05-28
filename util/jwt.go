package util

import (
	"file_server/global"
	"file_server/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var jwtKey = []byte("this_is_the_key")

type Claims struct
{
	Phone string
	jwt.StandardClaims
}
//token就是通过一些密码学手段记录了用户的身份凭证
func ReleaseToken(phone string)(string,error){
	expirarionTime := time.Now().Add(7*2*time.Hour)
	claims :=&Claims{
		Phone: phone,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: expirarionTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "Dispa1r",
			Subject: "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)//生成token
	tokenString,err :=token.SignedString(jwtKey)
	if err!=nil{
		return "", err
	}
	return tokenString,err
}
func ParseToken(tokenString string)(*jwt.Token,*Claims,error){
	Claims :=&Claims{}
	token,err :=jwt.ParseWithClaims(tokenString,Claims,func(token *jwt.Token)(i interface{},err error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if ok {
			return jwtKey,nil
		}else {
			return
		}
	})
	return token,Claims,err
}

func AuthMiddleware(token1 string,phone string) bool{
	token,claims,err :=ParseToken(token1)
	log.Println(claims.Phone)
	if err!=nil ||!token.Valid{//valid函数会自动判断token是否超时，也就是硕那个是无效字段。。。
		log.Println("授权失败")
		return false

	}

	//根据token中的phone和数据库中的phone以及传入的phone对比
	userId :=claims.Phone
	if userId != phone{
		return false
	}
	var user model.TblUser
	if global.DBEngine.Where("phone = ?",userId).First(&user).Error !=nil{
		return false
	}
	//遇到一个小坑 值得记录一下
	return true
}
