package handle

import (
	"file_server/util"
	"fmt"
	"math"
	"net/http"
	"strconv"
	rpool "file_server/cache"
	"time"
)

type MulPartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}
func InitMultiPartUpload(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	phone := r.Form.Get("phone")
	filehash := r.Form.Get("filehash")
	filesize,_ := strconv.Atoi(r.Form.Get("filesize"))
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	upInfo := MulPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   phone+fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize:  5*1024*1024,
		ChunkCount: int(math.Ceil(float64(filesize)/(5*1024*1024))),
	}
	//写入缓存
	rConn.Do("HSET","MP_"+upInfo.UploadID,"chunkcount",upInfo.ChunkCount)
	rConn.Do("HSET","MP_"+upInfo.UploadID,"filehash",upInfo.FileHash)
	rConn.Do("HSET","MP_"+upInfo.UploadID,"filesize",upInfo.FileSize)
	w.Write(util.NewRespMsg(0,"OK",upInfo).JSONBytes())
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request){
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	
}
