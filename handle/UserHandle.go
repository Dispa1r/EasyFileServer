package handle

import (
	rpool "file_server/cache"
	"file_server/dao"
	"file_server/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter,r *http.Request){
	if r.Method == http.MethodGet{
		data,err :=ioutil.ReadFile("./static/view/signup.html")
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username :=r.Form.Get("username")
	password :=r.Form.Get("password")
	email :=r.Form.Get("email")
	phone :=r.Form.Get("phone")
	if len(username)<3||len(password)<5 || len(email)==0||len(phone)!=11{
		w.Write([]byte("invalid information"))
		return
	}
	suc :=dao.UserSignUp(username,password,email,phone)
	code :=util.RandomCode()
	result := util.SendEmail(email,code)
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	rConn.Do("SETEX",email,300,code)
	//将email和验证码存入redis
	if result{
		fmt.Println("validate email has been send")
	}
	if suc{
		w.Write([]byte("SUCCESS"))
		w.WriteHeader(200)
	}else {
		w.Write([]byte("fail"))
		w.WriteHeader(500)
	}

}

func ValidateUserByEmail(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	code := r.Form.Get("code")
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	codeRedis, err := redis.String(rConn.Do("GET", email))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", codeRedis)
	}
	if code==codeRedis{
		rConn.Do("del",email)
		result :=dao.UserValidate(phone)
		if result{
			resp := util.RespMsg{
				Code: 0,
				Msg:  "OK",
				Data: struct {
					Location string
				}{
					Location: "http://" + r.Host + "/static/view/signin.html",
				},
			}
			w.Write(resp.JSONBytes())
		}
	}else {
		w.WriteHeader(500)
		w.Write([]byte("fail to validate"))
	}
}

func UserSignin(w http.ResponseWriter,r *http.Request){
	if r.Method == http.MethodGet{
		data,err :=ioutil.ReadFile("./static/view/signin.html")
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	phone := r.Form.Get("phone")
	password := r.Form.Get("password")
	result := dao.UserSignin(phone,password)
	if result!=true{
		w.Write([]byte("FAILED"))
		return
	}
	suc := dao.CheckUserValidate(phone)
	if !suc{
		w.Write([]byte("please validate your phone or email"))
		return
	}
	token,err := util.ReleaseToken(phone)
	if err!=nil{
		fmt.Println("fail to generate the token")
	}
	//dao.UpdateToken(phone,token)
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	rConn.Do("SET", phone,token,"EX",300*12*5)
	//redis存储token
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: phone,
			Token:    token,
		},
	}

	w.Write(resp.JSONBytes())

}
func HTTPInteceptor(h http.HandlerFunc)http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		phone := r.Form.Get("phone")
		token := r.Form.Get("token")
		isValidToken := util.AuthMiddleware(token,phone)
		if !isValidToken {
			w.WriteHeader(403)
			return
		}
		h(w,r)
	})
}
