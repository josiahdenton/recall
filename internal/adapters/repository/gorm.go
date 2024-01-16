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

func (g GormInstance) ModifyZettel(zettel domain.Zettel) {
	err := g.db.Save(&zettel).Commit().Error
	if err != nil {
		log.Printf("failed to save zettel: %v", err)
	}
}

func (g GormInstance) AllZettels() []domain.Zettel {
	var zettels []domain.Zettel
	err := g.db.Find(&zettels).Error
	if err != nil {
		log.Printf("failed to find all zettels: %v", err)
	}
	return zettels
}

func (g GormInstance) Zettel(id uint) *domain.Zettel {
	zettel := &domain.Zettel{}
	err := &g.db.Preload(clause.Associations).First(zettel, id).Error
	if err != nil {
		log.Printf("failed to get zettel (%d): %v", id, err)
	}
	return zettel
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
	err := g.db.Where("archive = ?", false).Find(&tasks)
	if err != nil {
		log.Printf("failed to get all tasks: %v", err)
	}
	// sort them...
	sort.Slice(tasks, func(i, j int) bool {
		if (reflect.ValueOf(tasks[i].Due).IsZero() && reflect.ValueOf(tasks[j].Due).IsZero()) || tasks[i].Due.Equal(tasks[j].Due) {
			return tasks[i].Title < tasks[j].Title
		}
		return tasks[i].Due.Before(tasks[j].Due)
	})
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

func (g GormInstance) ModifySettings(settings domain.Settings) {
	err := g.db.Save(&settings).Commit().Error
	if err != nil {
		log.Printf("failed to save settings: %v", err)
	}
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
	err = g.db.AutoMigrate(&domain.Settings{})
	if err != nil {
		return fmt.Errorf("failed to migrate settings: %w", err)
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
	err = g.db.AutoMigrate(&domain.Zettel{})
	if err != nil {
		return fmt.Errorf("failed to migrate zettel: %w", err)
	}

	return nil
}
