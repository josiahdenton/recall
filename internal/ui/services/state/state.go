package state

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"log"
)

var (
	FailedItemConversion = errors.New("failed item conversion")
)

type Repository interface {

	// Tasks

	Task(uint) *domain.Task
	AllTasks() []domain.Task
	ArchivedTasks() []domain.Task
	ModifyTask(domain.Task) domain.Task
	DeleteTask(uint)
	UnlinkTaskResource(*domain.Task, *domain.Resource)
	UnlinkTaskStep(*domain.Task, *domain.Step)
	UnlinkTaskStatus(*domain.Task, *domain.Status)
	UndoDeleteTask(uint)
	ModifyStep(step domain.Step) domain.Step

	// Cycles

	Cycle(uint) *domain.Cycle // etc...
	AllCycles() []domain.Cycle
	ModifyCycle(domain.Cycle) domain.Cycle

	// Accomplishments

	Accomplishment(uint) *domain.Accomplishment
	ModifyAccomplishment(domain.Accomplishment) domain.Accomplishment
	DeleteAccomplishment(uint)
	UndoDeleteAccomplishment(uint)

	// Resources

	ModifyResource(domain.Resource) domain.Resource
	AllResources() []domain.Resource

	LoadRepository() error
}

const (
	// Pairs //

	Task = iota
	Tasks
	Resource
	Resources
	Cycles
	Cycle

	// Singles //

	Accomplishment
	Step
	Status
)

type Type = int

const (
	View = iota
	// Edit mode is for showing the form
	Edit
	// Focus - idea - this is like Zen
	Focus
)

type Mode = int

type switchModeMsg struct {
	mode Mode
}

func SwitchMode(mode Mode) tea.Cmd {
	return func() tea.Msg {
		return switchModeMsg{mode: mode}
	}
}

type ModeSwitchMsg struct {
	Previous Mode
	Current  Mode
}

func New(path string, storage Repository) *State {
	// TODO - call and setup storage...
	return &State{
		Mode:    View,
		Storage: storage,
	}
}

// State represents the state of recall
type State struct {
	Mode    Mode
	Storage Repository
}

func (s *State) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case saveStateMsg:
		cmd = s.Save(Request{
			State: msg.State,
			Type:  msg.Type,
		})
	case deleteStateMsg:
		cmd = s.Delete(Request{
			ID:   msg.ID,
			Type: msg.Type,
		})
	case loadStateMsg:
		log.Printf("loadStateMsg")
		cmd = s.Load(Request{
			ID:   msg.ID,
			Type: msg.Type,
		})
	case switchModeMsg:
		cmd = s.ChangeMode(msg.mode)
	}

	return cmd
}

type SavedStateMsg struct {
	State any
	Type  Type
	Err   error
}

// Save is called after we finish
// adding or creating a new item
// use the repository
func (s *State) Save(r Request) tea.Cmd {
	return func() tea.Msg {
		var state any
		var err error
		switch r.Type {
		case Task:
			if item, ok := r.State.(domain.Task); ok {
				state = s.Storage.ModifyTask(item)
			} else {
				err = FailedItemConversion
			}
		case Tasks:
			// no mass edits supported yet!
		case Resource:
			if item, ok := r.State.(domain.Resource); ok {
				state = s.Storage.ModifyResource(item)
			} else {
				err = FailedItemConversion
			}
		case Resources:
			// no mass edits supported yet!
		case Cycle:
			if item, ok := r.State.(domain.Cycle); ok {
				state = s.Storage.ModifyCycle(item)
			} else {
				err = FailedItemConversion
			}
		case Cycles:
			// no mass edits supported yet!
		case Accomplishment:
			if item, ok := r.State.(domain.Accomplishment); ok {
				state = s.Storage.ModifyAccomplishment(item)
			} else {
				err = FailedItemConversion
			}
		case Step:
			if item, ok := r.State.(domain.Step); ok {
				state = s.Storage.ModifyStep(item)
			} else {
				err = FailedItemConversion
			}
		}

		return SavedStateMsg{
			State: state,
			Type:  r.Type,
			Err:   err,
		}
	}
}

type DeletedStateMsg struct {
	Type Type
	ID   uint
}

func (s *State) Delete(r Request) tea.Cmd {
	return func() tea.Msg {
		switch r.Type {
		case Task:
			s.Storage.DeleteTask(r.ID)
		case Tasks:
			// no mass edits supported yet!
		case Resource:
			s.Storage.UnlinkTaskResource(r.Parent.(*domain.Task), r.State.(*domain.Resource))
		case Resources:
			// no mass edits supported yet!
		case Cycle:
		case Cycles:
			// no mass edits supported yet!
		case Accomplishment:
		case Step:
		}

		return DeletedStateMsg{Type: r.Type, ID: r.ID}
	}
}

type LoadedStateMsg struct {
	State any
	Type  Type
}

func (s *State) Load(r Request) tea.Cmd {
	log.Printf("loading...")
	return func() tea.Msg {
		var state any
		switch r.Type {
		case Task:
			state = &mockTask
			log.Printf("loaded mock task!!!")
		case Tasks:
			state = mockTasks
			log.Printf("loaded mock tasks")
		case Resource:
		case Resources:
			state = mockResources
			log.Printf("loaded mock resources")
		case Cycle:
		case Cycles:
		case Accomplishment:
		}

		log.Printf("loaded state")
		return LoadedStateMsg{
			State: state,
			Type:  r.Type,
		}
	}
}

func (s *State) ChangeMode(mode Mode) tea.Cmd {
	return func() tea.Msg {
		prev := s.Mode
		s.Mode = mode
		return ModeSwitchMsg{
			Previous: prev,
			Current:  mode,
		}
	}
}

// Messages

type Request struct {
	ID    uint
	State any
	Type  Type
	// ParentID for removing associations in a delete request
	Parent     any
	ParentType Type
}

type deleteStateMsg struct {
	ID   uint
	Type Type
}

func Delete(r Request) tea.Cmd {
	return func() tea.Msg {
		return deleteStateMsg{
			ID:   r.ID,
			Type: r.Type,
		}
	}
}

type saveStateMsg struct {
	State any
	Type  Type
}

func Save(r Request) tea.Cmd {
	return func() tea.Msg {
		return saveStateMsg{
			State: r.State,
			Type:  r.Type,
		}
	}
}

type loadStateMsg struct {
	// ID of 0 is loading all
	ID   uint
	Type Type
}

func Load(r Request) tea.Cmd {
	return func() tea.Msg {
		return loadStateMsg{
			ID:   r.ID,
			Type: r.Type,
		}
	}
}
