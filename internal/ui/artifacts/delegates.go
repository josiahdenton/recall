package artifacts

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type artifactDelegate struct{}

func (d artifactDelegate) Height() int  { return 1 }
func (d artifactDelegate) Spacing() int { return 1 }
func (d artifactDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d artifactDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	artifact, ok := item.(*domain.Artifact)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderArtifact(artifact, index == m.Index()))
}
