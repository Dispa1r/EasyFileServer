package dao

import (
	"errors"
	"file_server/global"
	"file_server/model"
	"fmt"
	"time"
)

type TableFile struct {
	FileName string
	FileSha string
	FileSize int64
	Location string
	UploadAt time.Time
}

func UploadFileToDB(file_sha1,file_name,file_addr string,file_size int64)bool{
	fileNew := model.TblFile{
		FileAddr: file_addr,
		FileName: file_name,
		FileSha1: file_sha1,
		FileSize: file_size,
		Status: 1,
	}
	if global.DBEngine.Create(&fileNew).Error !=nil{
		return false
	}
	return true
}

func GetFileMeta(filehash string)(* TableFile,error){
	 
	file:= model.TblFile{}
	if global.DBEngine.Where("file_sha1=?", filehash).Where("status=?",1).Find(&file).Error !=nil{
		fmt.Printf("cannot find record in DB")
		return &TableFile{},errors.New("error to find")
	}
	tfile := TableFile{
		FileName:file.FileName,
		FileSha:  file.FileSha1,
		FileSize: file.FileSize,
		Location: file.FileAddr,
		UploadAt: file.UpdatedAt,
	}
	return &tfile,nil
}
