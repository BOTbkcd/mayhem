package entities

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Task struct {
	gorm.Model
	Title              string `gorm:"notnull"`
	Description        string
	Steps              []Step
	Deadline           time.Time
	Priority           int //3: High, 2: Mid, 1: Low, 0: No Priority
	IsFinished         bool
	IsRecurring        bool
	StartTime          time.Time //Applicable only for recurring tasks
	RecurrenceInterval int       // in days
	RecurChildren      []RecurTask
	StackID            uint
}

func (t Task) Save() Entity {
	DB.Save(&t)
	return t
}

func (t Task) Delete() {
	//Unscoped() is used to ensure hard delete, where task will be removed from db instead of being just marked as "deleted"
	DB.Unscoped().Select(clause.Associations).Delete(&t)
}

func (t Task) LatestRecurTask() (RecurTask, int64) {
	recurTask := RecurTask{}
	//localtime modifier has to be added to DATE other wise UTC time would be used
	result := DB.Last(&recurTask, "task_id = ? AND deadline <  DATE('now', 'localtime', 'start of day', '+1 day')", t.ID)

	// if t.IsFinished != recurTask.IsFinished {
	// 	t.IsFinished = recurTask.IsFinished
	// 	t.Save()
	// }
	return recurTask, result.RowsAffected
}

func (t Task) RemoveFutureRecurTasks() {
	DB.Unscoped().Where("deadline >=  DATE('now', 'start of day') AND task_id = ?", t.ID).Delete(&RecurTask{})
}
