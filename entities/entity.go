package entities

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Entity interface {
	Save() Entity
	Delete()
}

var DB *gorm.DB

func InitializeDB() {
	db, err := gorm.Open(sqlite.Open(".todo.db"), &gorm.Config{
		//Silent mode ensures that errors logs don't interfere with the view
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Stack{}, &Task{}, &Step{}, &RecurTask{})

	DB = db
}
