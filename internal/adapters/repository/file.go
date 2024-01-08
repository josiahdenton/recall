package repository

import (
	"encoding/json"
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	tasksChanged = iota
	accomplishmentsChanged
	cyclesChanged
	settingsChanged
)

type Feature = int

func NewFileStorage(path string) *FileStorage {
	// use to track when changes are made
	changes := make(map[Feature]bool)
	changes[tasksChanged] = false
	changes[accomplishmentsChanged] = false
	changes[cyclesChanged] = false
	changes[settingsChanged] = false

	return &FileStorage{path: path, changes: changes}
}

type FileStorage struct {
	path            string
	accomplishments map[string]domain.Accomplishment
	tasks           map[string]domain.Task
	changes         map[Feature]bool
	settings        domain.Settings
	cycles          []domain.Cycle
}

func (fl *FileStorage) Task(id string) *domain.Task {
	task, ok := fl.tasks[id]
	if !ok {
		return &domain.Task{}
	}
	return &task
}

func (fl *FileStorage) Cycle(id string) *domain.Cycle {
	for _, cycle := range fl.cycles {
		if cycle.Id == id {
			return &cycle
		}
	}
	return &domain.Cycle{}
}

func (fl *FileStorage) Accomplishment(id string) *domain.Accomplishment {
	accomplishment, ok := fl.accomplishments[id]
	if !ok {
		return &domain.Accomplishment{}
	}
	return &accomplishment
}

func (fl *FileStorage) SaveAccomplishment(accomplishment domain.Accomplishment) {
	fl.accomplishments[accomplishment.Id] = accomplishment
}

func (fl *FileStorage) AllAccomplishments(ids []string) []domain.Accomplishment {
	accomplishments := make([]domain.Accomplishment, len(ids))
	for i, id := range ids {
		accomplishment, ok := fl.accomplishments[id]
		if ok {
			accomplishments[i] = accomplishment
		}
	}
	return accomplishments
}

func (fl *FileStorage) SaveTask(task domain.Task) {
	fl.tasks[task.Id] = task
}

func (fl *FileStorage) AllTasks() []domain.Task {
	tasks := make([]domain.Task, len(fl.tasks))
	i := 0
	for _, task := range fl.tasks {
		tasks[i] = task
		i++
	}
	return tasks
}

func (fl *FileStorage) SaveCycle(cycle domain.Cycle) {
	// replace the cycle with the matching ID...
	replaced := false
	for i, cycle := range fl.cycles {
		if cycle.Id == cycle.Id {
			fl.cycles[i] = cycle
			replaced = true
		}
	}
	if !replaced {
		// this is a new cycle
		fl.cycles = append(fl.cycles, cycle)
	}
}

func (fl *FileStorage) AllCycles() []domain.Cycle {
	return fl.cycles
}

func (fl *FileStorage) SaveChanges() error {
	for feature, changed := range fl.changes {
		if changed {
			err := fl.saveFeatureChanges(feature)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (fl *FileStorage) saveFeatureChanges(feature Feature) error {
	switch feature {
	case accomplishmentsChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY, 0644)
		err = fl.writeToFile(f, accomplishmentsLayout{Accomplishments: fl.accomplishments})
		if err != nil {
			return err
		}
		return f.Close()
	case cyclesChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY, 0644)
		err = fl.writeToFile(f, cyclesLayout{Cycles: fl.cycles})
		if err != nil {
			return err
		}
		return f.Close()
	case tasksChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY, 0644)
		err = fl.writeToFile(f, tasksLayout{Tasks: fl.tasks})
		if err != nil {
			return err
		}
		return f.Close()
	case settingsChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY, 0644)
		err = fl.writeToFile(f, settingsLayout{Settings: fl.settings})
		if err != nil {
			return err
		}
		return f.Close()
	}
	return nil
}

func (fl *FileStorage) featureFilePath(feature Feature) string {
	switch feature {
	case accomplishmentsChanged:
		return fmt.Sprintf("%s/%s", fl.path, domain.AccomplishmentsFileName)
	case cyclesChanged:
		return fmt.Sprintf("%s/%s", fl.path, domain.CyclesFileName)
	case tasksChanged:
		return fmt.Sprintf("%s/%s", fl.path, domain.TasksFileName)
	case settingsChanged:
		return fmt.Sprintf("%s/%s", fl.path, domain.SettingsFileName)
	}
	return ""
}

func (fl *FileStorage) SaveSettings(settings domain.Settings) {
	fl.settings = settings
}

func (fl *FileStorage) LoadRepository() error {
	// TODO - add more context to errors...
	err := os.Mkdir(fl.path, 0755)
	if !os.IsExist(err) {
		return err
	}

	err = fl.createMissing()
	if err != nil {
		return err
	}
	// accomplishments
	path := fmt.Sprintf("%s/%s", fl.path, domain.AccomplishmentsFileName)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var accomplishments accomplishmentsLayout
	err = json.Unmarshal(bytes, &accomplishments)
	if err != nil {
		return err
	}
	fl.accomplishments = accomplishments.Accomplishments

	// tasks
	path = fmt.Sprintf("%s/%s", fl.path, domain.TasksFileName)
	bytes, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	var tasks tasksLayout
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		return err
	}
	fl.tasks = tasks.Tasks

	// cycles
	bytes, err = os.ReadFile(fmt.Sprintf("%s/%s", fl.path, domain.CyclesFileName))
	if err != nil {
		return err
	}
	var cycles cyclesLayout
	err = json.Unmarshal(bytes, &cycles)
	if err != nil {
		return err
	}
	fl.cycles = cycles.Cycles

	return nil
}

func (fl *FileStorage) createMissing() error {
	existingFiles := make(map[string]bool)
	existingFiles[domain.AccomplishmentsFileName] = false
	existingFiles[domain.TasksFileName] = false
	existingFiles[domain.CyclesFileName] = false

	err := filepath.WalkDir(fl.path, func(path string, d fs.DirEntry, err error) error {
		if _, ok := existingFiles[d.Name()]; ok {
			existingFiles[d.Name()] = true
		}
		return nil
	})
	if err != nil {
		return err
	}

	if !existingFiles[domain.TasksFileName] {
		var layout tasksLayout
		layout.Tasks = make(map[string]domain.Task)
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, domain.TasksFileName))
		if err != nil {
			return err
		}
		err = fl.writeToFile(f, layout)
		if err != nil {
			return err
		}
	}

	if !existingFiles[domain.AccomplishmentsFileName] {
		var layout accomplishmentsLayout
		layout.Accomplishments = make(map[string]domain.Accomplishment)
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, fmt.Sprintf("%s/%s", fl.path, domain.AccomplishmentsFileName)))
		if err != nil {
			return err
		}
		err = fl.writeToFile(f, layout)
		if err != nil {
			return err
		}
	}

	if !existingFiles[domain.CyclesFileName] {
		var layout cyclesLayout
		layout.Cycles = make([]domain.Cycle, 0)
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, fmt.Sprintf("%s/%s", fl.path, domain.CyclesFileName)))
		if err != nil {
			return err
		}
		err = fl.writeToFile(f, layout)
		if err != nil {
			return err
		}
	}

	if !existingFiles[domain.SettingsFileName] {
		var layout settingsLayout
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, fmt.Sprintf("%s/%s", fl.path, domain.SettingsFileName)))
		if err != nil {
			return err
		}
		err = fl.writeToFile(f, layout)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fl *FileStorage) writeToFile(f *os.File, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return f.Close()
}
