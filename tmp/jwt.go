package common

import (
	"github.com/dgrijalva/jwt-go"
	"goWeb/model"
	"time"
)

var jwtKey = []byte("this_is_the_key")

type Claims struct
{
	UserId uint
	jwt.StandardClaims
}
//token就是通过一些密码学手段记录了用户的身份凭证
func ReleaseToken(user model.User)(string,error){
	expirarionTime := time.Now().Add(7*2*time.Hour)
	claims :=&Claims{
		UserId: user.ID,
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
		return jwtKey,nil
	})
	return token,Claims,err
}