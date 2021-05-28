package model

import "github.com/jinzhu/gorm"

type TblFile struct {
	*gorm.Model
	FileSha1 string `gorm:"type:char(40);not null;unique_index"`
	FileName string `gorm:"type:varchar(256);not null"`
	FileSize int64 `gorm:"type:bigint(20);default:'0'"`
	FileAddr string `gorm:"type:varchar(1024);not null"`
	Status int `gorm:"type:int(11);default:'0';key:'idx_status'"`
	Ext1 int `gorm:"type:int(11);default:'0'"`
	Ext2 string `gorm:"type:text;"`

}
