package db

import (
	"echo-demo-project/server/db/seeders"
	"echo-demo-project/server/models"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
)

func Init() *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	fmt.Println(dataSourceName)

	db, err := gorm.Open(os.Getenv("DB_DRIVER"), dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Post{})

	db.Model(&models.Post{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	userSeeder := seeders.NewUserSeeder(db)
	userSeeder.SetUsers()

	return db
}
