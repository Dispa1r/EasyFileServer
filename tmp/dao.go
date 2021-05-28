package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"goWeb/model"
	"log"
)
var DB *gorm.DB
func GetDB() *gorm.DB{
	return DB
}
func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username :=viper.GetString("datasource.username")
	//password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	dsn := fmt.Sprintf("%s:001228@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		host,
		port,
		database,
		charset,
		)
	log.Println(dsn)
	//dsn := "root:001228@tcp(127.0.0.1:3306)/ginVueTest?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(driverName, dsn)
	//
	//if err != nil {
	//	log.Println("data base connect failed")
	//}
	db.LogMode(true)
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}
