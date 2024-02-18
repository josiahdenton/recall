package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"io"
	"strings"
)

func RenderTaskHeader(task *domain.Task) string {
	/**
	Task Name
	Due 11/11/23 (tag1, tag2, tag3, tag4, ...)
	Active Since 10/10/23

	Steps 3/5
	*/
	var b strings.Builder
	b.WriteString(task.Title)
	b.WriteString("\n")
	b.WriteString(task.Due.String())
	b.WriteString("\t")
	b.WriteString(task.Tags)

	return b.String()
}

// RenderTask will render the task
func renderTask(task *domain.Task, selected bool) string {
	style := styles.DefaultStyle

	if selected {
		style = styles.SelectedStyle
	}

	return style.Render(task.Title)
}

func TasksToListItems(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i := range tasks {
		item := &tasks[i]
		items[i] = item
	}
	return items
}

type TaskDelegate struct{}

func (d TaskDelegate) Height() int  { return 1 }
func (d TaskDelegate) Spacing() int { return 0 }
func (d TaskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d TaskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	task, ok := item.(*domain.Task)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderTask(task, index == m.Index()))
}
