package tasks

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	SelectedTaskStyle = styles.SelectedItemStyle.Copy()
	DefaultTaskStyle  = styles.DefaultItemStyle.Copy()
)

const MaxTagSz = 20

func renderTask(t *domain.Task, selected bool) string {
	cursor := "  "
	focusMarker := "\U000F0EFF"
	cursorStyle := styles.CursorStyle
	itemStyle := DefaultTaskStyle
	focusMarkerStyle := styles.SecondaryGrayStyle
	if selected {
		cursor = "> "
		itemStyle = SelectedTaskStyle
	}

	if t.Active {
		focusMarker = "\U000F0EFF"
		focusMarkerStyle = styles.SecondaryColorStyle
	}

	title := itemStyle.Width(60).Render(t.Title)
	tags := itemStyle.Width(MaxTagSz + 8).Render(styles.Summary(t.Tags, MaxTagSz))

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(cursor), title, tags, focusMarkerStyle.Render(focusMarker))
}
