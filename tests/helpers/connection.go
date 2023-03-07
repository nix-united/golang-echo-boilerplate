package helpers

import (
	mocket "github.com/selvatico/go-mocket"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, err := gorm.Open(mocket.DriverName, "connection_string")
	if err != nil {
		panic(err.Error())
	}
	return db
}
