package dao

import (
	"errors"
	"file_server/global"
	"file_server/model"
	"time"
)
//TODO:qrcode login，手机号验证和登陆
func OnUserFileUploadFinished(phone,filehash,filename string,filesize int64)bool{
	UserFile := model.TblUserFile{
		Status:     0,
		FileSize:   filesize,
		FileSha1:   filehash,
		Phone:      phone,
		FileName: filename,
		UploadAt: time.Now(),
	}
	global.DBEngine.Create(&UserFile)
	return true
}
func QueryUserFileMetas(username string,limit1 int)([]model.TblFile,error){
	var fileList []model.TblFile
	var tokenList []model.TblUserFile
	if global.DBEngine.Where("phone =?",username).Limit(limit1).Find(&tokenList).Error !=nil{
		return nil,errors.New("fail to load file info")
	}else {
		for _,v := range tokenList{
			var temp model.TblFile
			global.DBEngine.Where("file_sha1 =?",v.FileSha1).Find(&temp)
			fileList = append(fileList, temp)
		}
		return fileList,nil
	}
}
