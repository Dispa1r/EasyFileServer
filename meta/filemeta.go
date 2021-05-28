package meta

import (
	"file_server/dao"
	"fmt"
	"time"
)
//文件信息结构
type FileMeta struct {
	FileSha string
	FileName string
	FileSize int64
	Location string
	UploadAt time.Time
}

var FileMetas map[string]FileMeta

func init(){
	FileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta2(fmeta FileMeta){
	FileMetas[fmeta.FileSha] = fmeta

}
func UpdateFileMeta(fmeta FileMeta){
	FileMetas[fmeta.FileSha] = fmeta
	result :=dao.UploadFileToDB(fmeta.FileSha,fmeta.FileName,fmeta.Location,fmeta.FileSize)
	if !result{
		fmt.Println("save filemeta to database failed")
	}

}

func GetFileMetaFromDB(filesha string)(FileMeta,error){
	tfile,err :=dao.GetFileMeta(filesha)
	if err!=nil{
		return FileMeta{},err
	}
	fmeta :=FileMeta{
		FileSize: tfile.FileSize,
		FileName: tfile.FileName,
		FileSha: tfile.FileSha,
		Location: tfile.Location,
		UploadAt: tfile.UploadAt,
	}
	return fmeta,nil

}

func GetFileMeta(filesha string)FileMeta{
	return FileMetas[filesha]
}
func RemoveFileMeta(filesha string){
	delete(FileMetas,filesha)

}





func isFileExist(filesha string)bool{
	if _,ok :=FileMetas[filesha];ok{
		return true
	}else {
		return false
	}
}
