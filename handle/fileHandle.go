package handle

import (
	"encoding/json"
	"file_server/dao"
	"file_server/meta"
	"file_server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method =="GET"{
		data,err :=ioutil.ReadFile("./static/view/index.html")
		if err!=nil{
			fmt.Println("file not found")
			return
		}
		io.WriteString(w, string(data))

	}else if r.Method=="POST"{
		file,head,err :=r.FormFile("file")
		if err!=nil{
			fmt.Println("fail to get data from file")
			return
		}
		fileMeta :=meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/"+head.Filename,
		}


		newFile,err :=os.Create(fileMeta.Location)
		if err!=nil{
			fmt.Printf("fail to create file %v",err)
			return
		}
		fileMeta.FileSize,err =io.Copy(newFile,file)
		//内存中的文件流拷贝到创建os的文件流
		if err!=nil{
			fmt.Printf("copy file failed %v",err)
			return
		}
		fmt.Println("file upload finished")
		defer file.Close()
		defer newFile.Close()
		newFile.Seek(0,0)
		fileMeta.FileSha =util.FileSha1(newFile)
		fmt.Println(fileMeta)
		meta.UpdateFileMeta(fileMeta)
		r.ParseForm()
		phone :=r.Form.Get("phone")
		result := dao.OnUserFileUploadFinished(phone,fileMeta.FileSha,fileMeta.FileName,fileMeta.FileSize)
		if result {
			http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
		}else {
			w.Write([]byte("upload failed"))
		}

	}
}

func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload success")
}

//获取函数原信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filehash :=r.Form["filehash"][0]
	filemeta,err :=meta.GetFileMetaFromDB(filehash)
	if err!=nil{
		fmt.Printf("error to get the record %v",err)
	}
	byte,err :=json.Marshal(filemeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.Write(byte)

}
func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filehash :=r.Form.Get("filehash")
	fmeta,err :=meta.GetFileMetaFromDB(filehash)
	if err!=nil{
		fmt.Printf("error to get the record %v",err)
	}
	f,er:=os.Open(fmeta.Location)
	defer f.Close()
	if er!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	data,err :=ioutil.ReadAll(f)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octet-stream")
	w.Header().Set("Content-Description","attachment;filename=\""+fmeta.FileName+"\"")
	//让浏览器知道是下载文件
	w.Write(data)

}
func UpdateHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	//flag
	//filesha
	//filename
	opType :=r.Form.Get("op")
	filesha := r.Form.Get("filehash")
	filename := r.Form.Get("filename")

	if opType!="0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method!="POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta :=meta.GetFileMeta(filesha)
	curFileMeta.FileName =filename
	meta.UpdateFileMeta(curFileMeta)
	data,err :=json.Marshal(curFileMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(data)

}


func DeleteHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filesha := r.Form.Get("fileHash")
	fmeta :=meta.GetFileMeta(filesha)
	location := fmeta.Location
	err :=os.Remove(location)
	if err!=nil{
		fmt.Printf("fail to delete local file %v",err)
	}
	meta.RemoveFileMeta(filesha)
	w.WriteHeader(200)
}

func FileQueryHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	limitCnt,_:=strconv.Atoi(r.Form.Get("limit"))
	fileMetas,err:=dao.QueryUserFileMetas(r.Form.Get("phone"),limitCnt)
	if err!=nil{
		w.WriteHeader(500)
		w.Write([]byte("fail to get file"))
	}
	data,err := json.Marshal(fileMetas)
	w.Write(data)
}

