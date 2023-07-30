package entities

import (
	"time"

	"gorm.io/gorm"
)

type RecurTask struct {
	gorm.Model
	Deadline   time.Time `gorm:"index:idx_member"`
	IsFinished bool
	StackID    uint `gorm:"index:idx_member"`
	TaskID     uint
}

func (r RecurTask) Save() {
	DB.Save(&r)
}
