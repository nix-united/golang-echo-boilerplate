package helpers

import (
	mockRedis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
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

func InitRedis() (*mockRedis.Miniredis, *redis.Client) {
	mr, err := mockRedis.Run()
	if err != nil {
		panic(err)
	}

	return mr, redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
}
