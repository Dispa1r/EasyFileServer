package main

import (
	"bufio"
	"encoding/json"
	"file_server/DB"
	"file_server/config"
	"file_server/global"
	"file_server/model"
	"file_server/mq"
	"file_server/oss"
	"fmt"
	"github.com/spf13/viper"
	"log"
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
func UpdateFileLocation(filesha1,Location string)bool{
	file:= model.TblFile{}
	fmt.Println(global.DBEngine)
	err := global.DBEngine.Where("file_sha1 = ?",filesha1).Find(&file).Error
	if err!=nil{
		log.Panicf("%v",err)
		return false
	}
	file.FileAddr = Location
	global.DBEngine.Save(&file)
	return true
}
// ProcessTransfer : 处理文件转移
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.Location)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_ = UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	return true
}

func main() {
	//if !config.AsyncTransferEnable {

	//	log.Println("异步转移文件功能目前被禁用，请检查相关配置")
	//	return
	//}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	mq.StartConsume(
		config.TransExchangeQueueName,
		"transfer_oss",
		ProcessTransfer)
}
