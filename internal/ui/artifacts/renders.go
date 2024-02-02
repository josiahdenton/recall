package artifacts

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedArtifactStyle  = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(90)
	defaultArtifactStyle   = styles.SecondaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(90)
	nameStyle              = styles.PrimaryGray.Copy().Width(35)
	tagsStyle              = styles.SecondaryGray.Copy().Width(35).Align(lipgloss.Right)
	upcomingReleaseStyle   = styles.AccentColor.Copy().Width(2)
	successfulReleaseStyle = styles.SecondaryColor.Copy().Width(2)
	failedReleaseStyle     = styles.PrimaryColor.Copy().Width(2)
)

func renderArtifact(artifact *domain.Artifact, selected bool) string {
	style := defaultArtifactStyle
	releaseMarker := "\uF4C9"
	releaseStyle := upcomingReleaseStyle
	if selected {
		style = selectedArtifactStyle
	}
	if len(artifact.Releases) > 0 {
		switch artifact.Releases[len(artifact.Releases)-1].Outcome {
		case domain.AwaitingRelease:
			releaseStyle = upcomingReleaseStyle
		case domain.SuccessfulRelease:
			releaseStyle = successfulReleaseStyle
		case domain.FailedRelease:
			releaseStyle = failedReleaseStyle
		}
	}
	return style.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		releaseStyle.Render(releaseMarker),
		nameStyle.Render(artifact.Name),
		tagsStyle.Render(artifact.Tags),
	))
}
