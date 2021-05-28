package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"video_server/api/common"
	"video_server/api/model"
)

func AddVideo(c *gin.Context){
	DB := common.InitDB()
	aid,err := strconv.Atoi(c.PostForm("author_id"))
	if err!=nil{
		fmt.Println("author id wrong")
		return
	}
	name := c.PostForm("video_name")
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")
	newVideo := model.VideoInfo{
		AuthorId:aid,
		VideoName: name,
		DisplayCtime: ctime,
	}
	DB.Create(&newVideo)
	c.JSON(200,gin.H{
		"code":"200",
		"msg":"create video success",
	})


}

