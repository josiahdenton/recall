package artifacts

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedArtifactStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(70)
	defaultArtifactStyle  = styles.SecondaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(70)
	nameStyle             = styles.PrimaryGray.Copy().Width(35)
	tagsStyle             = styles.SecondaryGray.Copy().Width(35).Align(lipgloss.Right)
)

func renderArtifact(artifact *domain.Artifact, selected bool) string {
	style := defaultArtifactStyle
	if selected {
		style = selectedArtifactStyle
	}
	return style.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		nameStyle.Render(artifact.Name),
		tagsStyle.Render(artifact.Tags),
	))
}
