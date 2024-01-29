package artifact

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
	"time"
)

var (
	cursorStyle            = styles.PrimaryColor.Copy()
	resourceStyle          = styles.PrimaryGray.Copy().PaddingLeft(2)
	selectedResourceStyle  = styles.SecondaryColor.Copy().PaddingLeft(2)
	resourceMetaTitleStyle = styles.SecondaryGray.Copy().Width(5)
	// release
	selectedReleaseStyle   = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(70).Align(lipgloss.Center)
	defaultReleaseStyle    = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(70).Align(lipgloss.Center)
	upcomingReleaseStyle   = styles.AccentColor.Copy().PaddingRight(1)
	successfulReleaseStyle = styles.SecondaryColor.Copy().PaddingRight(1)
	failedReleaseStyle     = styles.PrimaryColor.Copy().PaddingRight(1)
	releaseMetaTitleStyle  = styles.SecondaryGray.Copy()
)

func renderResource(r *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := resourceStyle
	if selected {
		selectedMarker = ">"
		style = selectedResourceStyle
	}
	name := lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Width(2).Render(selectedMarker), style.Width(50).Render(r.Name))
	resourceType := lipgloss.JoinHorizontal(lipgloss.Left, resourceMetaTitleStyle.Render("Type"), style.Width(10).Render(r.StringType()))
	return lipgloss.JoinHorizontal(lipgloss.Left, name, resourceType)
}

func renderRelease(r *domain.Release, selected bool) string {
	var style lipgloss.Style
	var releaseMarker string
	var releaseKey string

	cardStyle := defaultReleaseStyle
	if selected {
		cardStyle = selectedReleaseStyle
	}

	switch r.Outcome {
	case domain.AwaitingRelease:
		style = upcomingReleaseStyle
		releaseKey = "Release "
		releaseMarker = "\uF056"
	case domain.SuccessfulRelease:
		style = successfulReleaseStyle
		releaseKey = "Released "
		releaseMarker = "\U000F05E0"
	case domain.FailedRelease:
		style = failedReleaseStyle
		releaseKey = "Released "
		releaseMarker = "\U000F0159"
	}
	releaseDate := formatDate(r.Date)

	var b strings.Builder
	b.WriteString(style.Render(releaseMarker))
	b.WriteString(releaseMetaTitleStyle.Render(releaseKey))
	b.WriteString(style.Render(releaseDate))
	b.WriteString(releaseMetaTitleStyle.Render("Owner "))
	b.WriteString(style.Render(r.Owner))

	return cardStyle.Render(b.String())
}

const longDateForm = "Jan 2, 2006 at 3:04pm (MST)"

func formatDate(date time.Time) string {
	return strings.Split(date.Format(longDateForm), "at")[0]
}
