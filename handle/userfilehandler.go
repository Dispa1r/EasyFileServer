package handle

import (
	"file_server/dao"
	"file_server/meta"
	"file_server/util"
	"fmt"
	"net/http"
	"strconv"
)

func UserInfoHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	phone := r.Form.Get("phone")
	//token := r.Form.Get("token")
	//fmt.Println(phone, "      ", token)
	//isValidToken := util.AuthMiddleware(token,phone)
	//if !isValidToken {
	//	w.WriteHeader(403)
	//	return
	//}
	user, err := dao.GetUserInfoDB(phone)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func TryFastUpload(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	phone := r.Form.Get("phone")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize := r.Form.Get("filesize")
	filemeta,err := meta.GetFileMetaFromDB(filehash)
	if err!=nil {
		fmt.Println("fail to get file meta from db")
		w.WriteHeader(500)
	}
	var zero meta.FileMeta
	if filemeta==zero{
		resp :=util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，转用普通上传接口",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
	}
	fileSize,_ := strconv.Atoi(filesize)

	suc :=dao.OnUserFileUploadFinished(phone,filehash,filename,int64(fileSize))
	if suc{
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
	}else{
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传失败",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
	}
}