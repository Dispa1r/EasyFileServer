package main

import (
	"file_server/DB"
	"file_server/cache"
	"file_server/global"
	"file_server/handle"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func init(){
	InitConfig()
	global.DBEngine = DB.InitDB()
}

func InitConfig(){
	workDir,_ :=os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir+"/config")
	err :=viper.ReadInConfig()
	if err!=nil{
		log.Println("read config failed")
	}
}

func main(){
	rConn := cache.RedisPool().Get()
	fmt.Println(rConn)
	http.Handle("/static/",http.StripPrefix("/static",http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload",handle.HTTPInteceptor(handle.UploadHandler))
	http.HandleFunc("/file/upload/suc",handle.UploadSucHandler)
	http.HandleFunc("/file/meta",handle.HTTPInteceptor(handle.GetFileMetaHandler))
	http.HandleFunc("/file/download",handle.HTTPInteceptor(handle.DownloadHandler))
	http.HandleFunc("/file/update",handle.HTTPInteceptor(handle.UpdateHandler))
	http.HandleFunc("/file/delete",handle.HTTPInteceptor(handle.DeleteHandler))
	http.HandleFunc("/user/signup",handle.SignUpHandler)
	http.HandleFunc("/user/code",handle.ValidateUserByEmail)
	http.HandleFunc("/user/signin",handle.UserSignin)
	http.HandleFunc("/user/info",handle.HTTPInteceptor(handle.UserInfoHandler))
	http.HandleFunc("/file/query",handle.HTTPInteceptor(handle.FileQueryHandler))
	http.HandleFunc("/file/fastupload",handle.HTTPInteceptor(handle.TryFastUpload))
	err :=http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Println("fail to start server")
	}

}
