package entities

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Skipping priority field, just sort them alphabetically
type Stack struct {
	gorm.Model
	Title            string `gorm:"notnull"`
	PendingTaskCount int
	Tasks            []Task
}

func InitializeStacks() (Stack, error) {
	stack := Stack{Title: "New Stack"}
	result := DB.Create(&stack)
	return stack, result.Error
}

func FetchAllStacks() ([]Stack, error) {
	var stacks []Stack
	result := DB.Model(&Stack{}).Preload("Tasks").Preload("Tasks.Steps").Find(&stacks)

	if len(stacks) == 0 {
		stack, err := InitializeStacks()
		return []Stack{stack}, err
	}

	return stacks, result.Error
}

func IncPendingCount(id uint) {
	stack := Stack{}
	DB.Find(&stack, id)
	stack.PendingTaskCount++
	stack.Save()
}

func (s Stack) PendingRecurringCount() int {
	recurTasks := []RecurTask{}
	//localtime modifier has to be added to DATE other wise UTC time would be used
	result := DB.Find(&recurTasks, "deadline >= DATE('now', 'localtime', 'start of day') AND deadline < DATE('now', 'localtime', 'start of day', '+1 day') AND stack_id = ? AND is_finished = false", s.ID)
	return int(result.RowsAffected)
}

func (s Stack) Save() Entity {
	DB.Save(&s)
	return s
}

func (s Stack) Delete() {
	//Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
	// DB.Unscoped().Delete(&s)
	DB.Unscoped().Select(clause.Associations).Delete(&s)
}
