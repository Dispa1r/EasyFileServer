package handle

import (
	"file_server/dao"
	"file_server/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	rpool "file_server/cache"
	"strings"
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

// UploadPartHandler : 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	//	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 2. 获得redis连接池中的一个连接
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := "./data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}
	// 4. 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	// 5. 返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

//1todo : 根据用户名称存储文件
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request){
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	phone := r.Form.Get("phone")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	data,err :=redis.Values(rConn.Do("HGETALL","MP_"+upid))
	if err!=nil{
		w.Write(util.NewRespMsg(-1,"complete failed",nil).JSONBytes())
		return
	}
	totalCount :=0
	chunkCount :=0
	for i:=0;i<len(data);i+=2{
		k:=string(data[i].([]byte)) //data里有好几个键值对，都以string的形式拿出来,
		v :=string(data[i+1].([]byte))
		if k== "chunkcount"{
			totalCount,_ =strconv.Atoi(v)
		}else if strings.HasPrefix(k, "chkidx_") && v == "1"{
			chunkCount++
		}
	}
	rConn.Do("del","MP_"+upid)
	//上传完成后要删除
	if totalCount!=chunkCount{
		w.Write(util.NewRespMsg(-2,"file not vaild",nil).JSONBytes())
	}
	//2todo:合并分块
    mergeFile(upid,filename,totalCount)
	fsize,_ := strconv.Atoi(filesize)
	dao.UploadFileToDB(filehash,filename,"",int64(fsize))//更新文件表
	dao.OnUserFileUploadFinished(phone,filehash,filename,int64(fsize))//更新用户文件表
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())

}

func mergeFile(upid,filename string,num int){
	path := "./data/"+upid+"/"+filename+".merge"
	fii, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer fii.Close()
	if err != nil {
		panic(err)
		return
	}
	for i := 1; i <= int(num); i++ {
		f, err := os.OpenFile("./data/"+upid+"/"+strconv.Itoa(int(i)), os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fii.Write(b)
		f.Close()
		os.Remove("./data/"+upid+"/"+strconv.Itoa(int(i)))
	}
}

func CancelMPUpload(w http.ResponseWriter, r *http.Request){
	//如何取消 删除本地文件
	// 清除redis缓存
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	dirpath := "./data/"+upid
	os.RemoveAll(dirpath)
	rConn := rpool.RedisPool().Get()//连接池中获取
	defer rConn.Close()
	rConn.Do("del","MP_"+upid)
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}
