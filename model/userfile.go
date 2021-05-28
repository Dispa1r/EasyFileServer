package model

import (
	"time"
)

type TblUserFile struct {
	ID  uint `gorm:"column:id"`
	LastUpdate time.Time `gorm:"column:last_update;default:null"`
	Status int `gorm:"column:status"`
	UploadAt time.Time `gorm:"column:upload_at;default:null"`
	FileSize int64 `gorm:"column:file_size;default:'0'"`
	FileSha1 string `gorm:"column:file_sha1;default:''"`
	FileName string `gorm:"column:file_name;default:''"`
	Phone string `gorm:"column:phone;default:''"`
}
