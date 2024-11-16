package dbM

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func DbConnect(conf string) *gorm.DB {
	var (
		err error
		res *gorm.DB
	)
	res, err = gorm.Open(mysql.Open(conf), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	sqlDB, err := res.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return res
}
