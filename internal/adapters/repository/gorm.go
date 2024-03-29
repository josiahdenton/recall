package repository

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"reflect"
	"sort"
)

func NewGormInstance(path string) (*GormInstance, error) {
	db, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		return nil, err
	}

	return &GormInstance{db: db}, err
}

type GormInstance struct {
	db *gorm.DB
}

func (g GormInstance) ModifyStatus(status domain.Status) {
	err := g.db.Save(&status).Commit().Error
	if err != nil {
		log.Printf("failed to save status: %v", err)
	}
}

func (g GormInstance) DeleteAccomplishment(id uint) {
	err := g.db.Delete(&domain.Accomplishment{}, id).Error
	if err != nil {
		log.Printf("failed to delete accomplisment (%d), for reason: %v", id, err)
	}
}

func (g GormInstance) UndoDeleteTask(id uint) {
	err := g.db.Unscoped().Model(&domain.Task{}).Where("id", id).Update("deleted_at", nil).Error
	if err != nil {
		log.Printf("failed to undo delete for task (%d) for reason: %v", id, err)
	}
}

func (g GormInstance) UndoDeleteAccomplishment(id uint) {
	err := g.db.Unscoped().Model(&domain.Accomplishment{}).Where("id", id).Update("deleted_at", nil).Error
	if err != nil {
		log.Printf("failed to undo delete for accomplishment (%d) for reason: %v", id, err)
	}
}

func (g GormInstance) UnlinkTaskResource(task *domain.Task, resource *domain.Resource) {
	err := g.db.Model(task).Association("Resources").Delete(resource)
	if err != nil {
		log.Printf("failed to delete resource (%d) associated with task (%d) due to: %+v", task.ID, resource.ID, err)
	}
}

func (g GormInstance) UnlinkTaskStep(task *domain.Task, step *domain.Step) {
	err := g.db.Model(task).Association("Steps").Delete(step)
	if err != nil {
		log.Printf("failed to delete step (%d) associated with task (%d) due to: %+v", task.ID, step.ID, err)
	}
	err = g.db.Unscoped().Delete(step).Error
	if err != nil {
		log.Printf("failed to delete step (%d) associated with task (%d) due to: %+v", task.ID, step.ID, err)
	}
}

func (g GormInstance) UnlinkTaskStatus(task *domain.Task, status *domain.Status) {
	err := g.db.Model(task).Association("Status").Delete(status)
	if err != nil {
		log.Printf("failed to delete status (%d) associated with task (%d) due to: %+v", task.ID, status.ID, err)
	}
	err = g.db.Unscoped().Delete(status).Error
	if err != nil {
		log.Printf("failed to delete status (%d) associated with task (%d) due to: %+v", task.ID, status.ID, err)
	}
}

func (g GormInstance) ModifyStep(step domain.Step) {
	err := g.db.Save(&step).Commit().Error
	if err != nil {
		log.Printf("failed to save step: %v", err)
	}
}

func (g GormInstance) Task(id uint) *domain.Task {
	task := &domain.Task{}
	err := g.db.Preload(clause.Associations).First(task, id).Error
	if err != nil {
		log.Printf("failed to get task (%d): %v", id, err)
	}
	return task
}

func (g GormInstance) DeleteTask(id uint) {
	err := g.db.Delete(&domain.Task{}, id).Error
	if err != nil {
		log.Printf("failed to delete task: %v", err)
	}
}

func (g GormInstance) Cycle(id uint) *domain.Cycle {
	cycle := &domain.Cycle{}
	err := g.db.Preload(clause.Associations).First(cycle, id)
	if err != nil {
		log.Printf("failed to get cycle: %v", err)
	}
	return cycle
}

func (g GormInstance) Accomplishment(id uint) *domain.Accomplishment {
	accomplishment := &domain.Accomplishment{}
	err := g.db.Preload(clause.Associations).First(accomplishment, id).Error
	if err != nil {
		log.Printf("failed to get accomplishment: %v", err)
	}
	return accomplishment
}

func (g GormInstance) ModifyAccomplishment(accomplishment domain.Accomplishment) {
	err := g.db.Save(&accomplishment).Commit().Error
	if err != nil {
		log.Printf("failed to save accomplishment: %v", err)
	}
}

func (g GormInstance) ModifyTask(task domain.Task) {
	err := g.db.Save(&task).Commit().Error
	if err != nil {
		log.Printf("failed to save task: %v", err)
	}
}

func (g GormInstance) AllTasks() []domain.Task {
	var tasks []domain.Task
	err := g.db.Where("archive = ?", false).Find(&tasks).Error
	if err != nil {
		log.Printf("failed to get all tasks: %v", err)
	}
	// sort them...
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Active && !tasks[j].Active
	})

	lastActiveIndex := 0
	for i, task := range tasks {
		if task.Active {
			lastActiveIndex = i
		}
	}

	// TODO: - could push this into a shared func...
	activeTasks := tasks[:lastActiveIndex+1]
	sort.Slice(activeTasks, func(i, j int) bool {
		if (reflect.ValueOf(activeTasks[i].Due).IsZero() && reflect.ValueOf(activeTasks[j].Due).IsZero()) || activeTasks[i].Due.Equal(activeTasks[j].Due) {
			return activeTasks[i].Title < activeTasks[j].Title
		}
		return activeTasks[i].Due.Before(activeTasks[j].Due)
	})

	if len(tasks) > 0 {
		inactiveTasks := tasks[lastActiveIndex+1:]
		sort.Slice(inactiveTasks, func(i, j int) bool {
			if (reflect.ValueOf(inactiveTasks[i].Due).IsZero() && reflect.ValueOf(inactiveTasks[j].Due).IsZero()) || inactiveTasks[i].Due.Equal(inactiveTasks[j].Due) {
				return inactiveTasks[i].Title < inactiveTasks[j].Title
			}
			return inactiveTasks[i].Due.Before(inactiveTasks[j].Due)
		})
	}

	return tasks
}

func (g GormInstance) ArchivedTasks() []domain.Task {
	var tasks []domain.Task
	err := g.db.Where(&domain.Task{Archive: true}).Find(&tasks)
	if err != nil {
		log.Printf("failed to get archived tasks: %v", err)
	}
	return tasks
}

func (g GormInstance) ModifyCycle(cycle domain.Cycle) {
	err := g.db.Save(&cycle).Commit().Error
	if err != nil {
		log.Printf("failed to save cycle: %v", err)
	}
}

func (g GormInstance) AllCycles() []domain.Cycle {
	var cycles []domain.Cycle
	err := g.db.Find(&cycles).Error
	if err != nil {
		log.Printf("failed to find all cycles: %v", err)
	}
	return cycles
}

func (g GormInstance) ModifyResource(resource domain.Resource) {
	err := g.db.Save(&resource).Commit().Error
	if err != nil {
		log.Printf("failed to save resource: %v", err)
	}
}

func (g GormInstance) AllResources() []domain.Resource {
	var resources []domain.Resource
	err := g.db.Find(&resources).Error
	if err != nil {
		log.Printf("failed to find all resources: %v", err)
	}
	return resources
}

func (g GormInstance) LoadRepository() error {
	// auto migrate all schemas
	err := g.db.AutoMigrate(&domain.Accomplishment{})
	if err != nil {
		return fmt.Errorf("failed to migrate accomplishment: %w", err)
	}
	err = g.db.AutoMigrate(&domain.Cycle{})
	if err != nil {
		return fmt.Errorf("failed to migrate cycle: %w", err)
	}
	err = g.db.AutoMigrate(&domain.Resource{})
	if err != nil {
		return fmt.Errorf("failed to migrate resource: %w", err)
	}
	err = g.db.AutoMigrate(&domain.Status{})
	if err != nil {
		return fmt.Errorf("failed to migrate status: %w", err)
	}
	err = g.db.AutoMigrate(&domain.Step{})
	if err != nil {
		return fmt.Errorf("failed to migrate step: %w", err)
	}
	err = g.db.AutoMigrate(&domain.Task{})
	if err != nil {
		return fmt.Errorf("failed to migrate task: %w", err)
	}

	return nil
}
