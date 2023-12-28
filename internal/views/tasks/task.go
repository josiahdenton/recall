package tasks

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"time"
)

const (
	None = iota
	High
	Low
)

type Priority int

type Task struct {
	Title             string
	Due               string
	Priority          Priority
	Active            bool
	Complete          bool
	TotalActiveTime   time.Duration
	LastActivatedTime time.Time
	Steps             []Step
}

func (t *Task) Activate() {
	t.Active = true
	t.LastActivatedTime = time.Now()
}

func (t *Task) Deactivate() {
	t.Active = false
	//delta := time.Now().Sub(t.LastActivatedTime)
	//t.TotalActiveTime += delta
	t.TotalActiveTime += time.Now().Sub(t.LastActivatedTime)
}

func (t *Task) Runtime() string {
	seconds := int((time.Now().Add(t.TotalActiveTime).Sub(t.LastActivatedTime)).Seconds())
	switch {
	case seconds < 60:
		return fmt.Sprintf("%d sec", seconds)
	case seconds < 3600:
		return fmt.Sprintf("%d min %d sec", seconds/60, seconds%60)
	case seconds < 86400:
		return fmt.Sprintf("%d hr %d min %d sec", seconds/60/60, (seconds/60)%60, seconds%60)
	default:
		return fmt.Sprintf("%d days %d hr %d min %d sec", seconds/60/60/24, (seconds/60/60)%24, (seconds/60)%60, seconds%60)
	}
}

var (
	taskStyle                 = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	activeStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("#2dd4bf"))
	lowPriorityStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	highPriorityStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	selectedStyle             = lipgloss.NewStyle().Background(lipgloss.Color("#1f2937"))
	selectedTaskStyle         = selectedStyle.Copy().Foreground(lipgloss.Color("#cbd5e1"))
	selectedActiveStyle       = selectedStyle.Copy().Foreground(lipgloss.Color("#99f6e4"))
	selectedLowPriorityStyle  = selectedStyle.Copy().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	selectedHighPriorityStyle = selectedStyle.Copy().Foreground(lipgloss.Color("#ef4444")).Bold(true)
)

func (t *Task) DetailedView() string {
	// title
	// meta section (runtime, due date, difficulty
	// subtasks
	// resources
	// status updates
	return ""
}

func (t *Task) Render(selected bool) string {
	var style lipgloss.Style
	var priorityStyle lipgloss.Style
	var priorityMarker string
	// TODO I need to clean all this logic up
	activeMarker := "\uF4C3"

	switch {
	case selected && t.Active:
		style = selectedActiveStyle
		activeMarker = "\uF444"
	case selected:
		style = selectedTaskStyle
	case t.Active:
		style = activeStyle
		activeMarker = "\uF444"
	default:
		style = taskStyle
	}

	switch {
	case t.Priority == None && selected:
		priorityStyle = selectedStyle
	case t.Priority == Low && selected:
		priorityStyle = selectedLowPriorityStyle
		priorityMarker = " *"
	case t.Priority == High && selected:
		priorityStyle = selectedHighPriorityStyle
		priorityMarker = " ***"
	case t.Priority == Low:
		priorityStyle = lowPriorityStyle
		priorityMarker = " *"
	case t.Priority == High:
		priorityStyle = highPriorityStyle
		priorityMarker = " ***"
	}

	content := style.Width(30).Render(activeMarker, t.Title)
	date := style.Width(10).Italic(true).Render(t.Due)
	priority := priorityStyle.Width(5).Render(priorityMarker)
	return lipgloss.JoinHorizontal(lipgloss.Top, content, priority, date)
}

func activeTask(tasks []Task) int {
	for i, task := range tasks {
		if task.Active {
			return i
		}
	}
	return -1
}

func hasActiveTask(tasks []Task) bool {
	for _, task := range tasks {
		if task.Active {
			return true
		}
	}
	return false
}
