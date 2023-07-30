package entities

import (
	"gorm.io/gorm"
)

type Step struct {
	gorm.Model
	Title      string
	IsFinished bool
	TaskID     uint
}

func (s Step) Save() Step {
	DB.Save(&s)
	return s
}

func (s Step) Delete() {
	DB.Delete(&s)
}
