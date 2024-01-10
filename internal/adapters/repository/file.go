package repository

import (
	"encoding/json"
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	tasksChanged = iota
	accomplishmentsChanged
	cyclesChanged
	settingsChanged
	resourcesChanged
)

type Feature = int

func NewFileStorage(path string) *FileStorage {
	log.Printf("adding path as: %v", path)
	// use to track when changes are made
	changes := make(map[Feature]bool)
	changes[tasksChanged] = false
	changes[accomplishmentsChanged] = false
	changes[cyclesChanged] = false
	changes[settingsChanged] = false
	changes[resourcesChanged] = false

	return &FileStorage{path: path, changes: changes}
}

type FileStorage struct {
	path            string
	accomplishments map[string]domain.Accomplishment
	tasks           map[string]domain.Task
	taskArchive     map[string]domain.Task
	resources       map[string]domain.Resource
	changes         map[Feature]bool
	settings        domain.Settings
	cycles          []domain.Cycle
}

func (fl *FileStorage) LinkedTasks(ids []string) []domain.Task {
	tasks := make([]domain.Task, len(ids))
	// TODO - fix ordering problem...
	for i, id := range ids {
		task, ok := fl.tasks[id]
		if ok {
			tasks[i] = task
		}
	}
	return tasks
}

func (fl *FileStorage) DeleteTask(id string) {
	delete(fl.tasks, id)
	fl.changes[tasksChanged] = true
}

func (fl *FileStorage) LinkedResources(ids []string) []domain.Resource {
	resources := make([]domain.Resource, len(ids))
	// TODO - fix ordering problem...
	for i, id := range ids {
		resource, ok := fl.resources[id]
		if ok {
			resources[i] = resource
		}
	}
	return resources
}

func (fl *FileStorage) SaveResource(resource domain.Resource) {
	fl.resources[resource.Id] = resource
	fl.changes[resourcesChanged] = true
}

func (fl *FileStorage) AllResources() []domain.Resource {
	resources := make([]domain.Resource, len(fl.resources))
	// TODO - fix ordering problem...
	i := 0
	for _, resource := range fl.resources {
		resources[i] = resource
		i++
	}
	return resources
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
	fl.changes[accomplishmentsChanged] = true
	fl.accomplishments[accomplishment.Id] = accomplishment
	// mark all tasks as Completed
	// TODO - this doesn't belong here... only here for now to work
	for _, taskId := range accomplishment.AssociatedTaskIds {
		if task, ok := fl.tasks[taskId]; ok {
			task.Archive = true
			fl.tasks[taskId] = task // TODO - may need this...
		}
	}
	fl.changes[tasksChanged] = true
}

func (fl *FileStorage) LinkedAccomplishments(ids []string) []domain.Accomplishment {
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
	fl.changes[tasksChanged] = true
	fl.tasks[task.Id] = task
}

func (fl *FileStorage) ArchivedTasks() []domain.Task {
	tasks := make([]domain.Task, 0)
	i := 0
	for _, task := range fl.tasks {
		if task.Archive {
			tasks = append(tasks, task)
		}
		i++
	}
	return tasks
}

func (fl *FileStorage) AllTasks() []domain.Task {
	tasks := make([]domain.Task, 0)
	i := 0
	for _, task := range fl.tasks {
		if !task.Archive {
			tasks = append(tasks, task)
		}
		i++
	}
	return tasks
}

func (fl *FileStorage) SaveCycle(updated domain.Cycle) {
	fl.changes[cyclesChanged] = true
	// replace the cycle with the matching ID...
	replaced := false
	for i, cycle := range fl.cycles {
		if cycle.Id == updated.Id {
			fl.cycles[i] = updated
			replaced = true
		}
	}
	if !replaced {
		// this is a new cycle
		fl.cycles = append(fl.cycles, updated)
	}
	// any update to a cycle should trigger this check
	fl.setActiveCycle()
}

func (fl *FileStorage) setActiveCycle() {
	// find the cycle with the most recent date in the past
	if len(fl.cycles) < 1 {
		return
	}
	mostRecentPastTime := &fl.cycles[0]
	for _, cycle := range fl.cycles {
		if mostRecentPastTime.StartDate.Before(cycle.StartDate) && cycle.StartDate.Before(time.Now()) {
			mostRecentPastTime = &cycle
		}
	}
	mostRecentPastTime.Active = true
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
	for feature := range fl.changes {
		fl.changes[feature] = false
	}
	return nil
}

func (fl *FileStorage) saveFeatureChanges(feature Feature) error {
	switch feature {
	case accomplishmentsChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		return fl.writeToFile(f, accomplishmentsLayout{Accomplishments: fl.accomplishments})
	case cyclesChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		return fl.writeToFile(f, cyclesLayout{Cycles: fl.cycles})
	case tasksChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		return fl.writeToFile(f, tasksLayout{Tasks: fl.tasks})
	case settingsChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		return fl.writeToFile(f, settingsLayout{Settings: fl.settings})
	case resourcesChanged:
		f, err := os.OpenFile(fl.featureFilePath(feature), os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		return fl.writeToFile(f, resourcesLayout{Resources: fl.resources})
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
	case resourcesChanged:
		return fmt.Sprintf("%s/%s", fl.path, domain.ResourcesFileName)
	}
	return ""
}

func (fl *FileStorage) SaveSettings(settings domain.Settings) {
	fl.changes[settingsChanged] = true
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

	// resources
	path := fmt.Sprintf("%s/%s", fl.path, domain.ResourcesFileName)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var resources resourcesLayout
	err = json.Unmarshal(bytes, &resources)
	if err != nil {
		return err
	}
	fl.resources = resources.Resources

	// accomplishments
	path = fmt.Sprintf("%s/%s", fl.path, domain.AccomplishmentsFileName)
	bytes, err = os.ReadFile(path)
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
	existingFiles[domain.ResourcesFileName] = false

	err := filepath.WalkDir(fl.path, func(path string, d fs.DirEntry, err error) error {
		if _, ok := existingFiles[d.Name()]; ok {
			existingFiles[d.Name()] = true
		}
		return nil
	})
	if err != nil {
		return err
	}

	if !existingFiles[domain.ResourcesFileName] {
		var layout resourcesLayout
		layout.Resources = make(map[string]domain.Resource)
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, domain.ResourcesFileName))
		if err != nil {
			return err
		}
		err = fl.writeToFile(f, layout)
		if err != nil {
			return err
		}
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
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, domain.AccomplishmentsFileName))
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
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, domain.CyclesFileName))
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
		f, err := os.Create(fmt.Sprintf("%s/%s", fl.path, domain.SettingsFileName))
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
