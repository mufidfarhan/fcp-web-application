package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(id int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *taskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) Store(task *model.Task) error {
	err := t.db.Create(task).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) Update(id int, task *model.Task) error {
	err := t.db.Model(&task).Where("id = ?", id).Updates(map[string]interface{}{
		"title":       task.Title,
		"deadline":    task.Deadline,
		"priority":    task.Priority,
		"category_id": task.CategoryID,
		"status":      task.Status,
	}).Error
	if err != nil {
		return err
	}

	return nil // TODO: replace this
}

func (t *taskRepository) Delete(id int) error {
	var task model.Task
	if err := t.db.Where("id = ?", id).Delete(&task).Error; err != nil {
		return err
	}

	return nil // TODO: replace this
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	var task model.Task
	err := t.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (t *taskRepository) GetList() ([]model.Task, error) {
	var task model.Task
	var tasks = []model.Task{}

	rows, err := t.db.Model(&task).Select("*").Rows()
	if err != nil {
		return tasks, err
	}

	defer rows.Close()
	for rows.Next() {
		t.db.ScanRows(rows, &tasks)
	}

	return tasks, nil // TODO: replace this
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	var taskCategory = []model.TaskCategory{}
	var task model.Task

	err := t.db.Model(&task).Select("tasks.id, tasks.title, categories.name as category").Joins("join categories on tasks.category_id = ? and categories.id = $1", id).Limit(1).Scan(&taskCategory).Error
	if err != nil {
		return taskCategory, err
	}

	return taskCategory, nil // TODO: replace this
}
