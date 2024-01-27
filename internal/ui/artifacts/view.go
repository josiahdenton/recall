package artifacts

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = styles.SecondaryGray.Copy()
)

func New() tea.Model {
	return Model{
		form: forms.NewArtifactForm(),
	}
}

type Model struct {
	artifacts list.Model
	showForm  bool
	form      tea.Model
	ready     bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var s string
	if m.showForm && m.ready {
		s = styles.WindowStyle.Render(m.form.View())
	} else if m.ready {
		s = styles.WindowStyle.Render(m.artifacts.View())
	}
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		artifacts := msg.State.([]domain.Artifact)
		m.artifacts = list.New(toItemList(artifacts), artifactDelegate{}, 50, 30)
		m.artifacts.Title = "Artifacts"
		m.artifacts.Styles.PaginationStyle = paginationStyle
		m.artifacts.Styles.Title = titleStyle
		m.artifacts.SetShowHelp(false)
		m.artifacts.KeyMap.Quit.Unbind()
		m.ready = true
	case state.SaveStateMsg:
		m.showForm = false
		if msg.Type == state.ModifyArtifact {
			artifact := msg.Update.(domain.Artifact)
			m.artifacts.InsertItem(len(m.artifacts.Items()), &artifact)
		}
	}

	if m.showForm && m.ready {
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.ready {
		m.artifacts, cmd = m.artifacts.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && m.showForm {
			m.showForm = false
		} else if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		}
	}

	if m.showForm {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if len(m.artifacts.Items()) < 1 {
				return m, tea.Batch(cmds...)
			}
			artifact := m.artifacts.SelectedItem().(*domain.Artifact)
			cmd = router.GotoPage(domain.ArtifactPage, artifact.ID)
			cmds = append(cmds, cmd)
		}
	}

	if m.artifacts.FilterState() == list.Filtering {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.showForm = true
		case "d":
			selected := m.artifacts.SelectedItem().(*domain.Artifact)
			m.artifacts.RemoveItem(m.artifacts.Index())
			cmds = append(cmds, deleteArtifact(selected.ID), toast.ShowToast("removed artifact!"))
		case "u":
			cmds = append(cmds, state.UndoDeleteState(), toast.ShowToast("undo!"))
		}
	}

	return m, tea.Batch(cmds...)
}

func toItemList(artifacts []domain.Artifact) []list.Item {
	items := make([]list.Item, len(artifacts))
	for i := range artifacts {
		item := &artifacts[i]
		items[i] = item
	}
	return items
}

func deleteArtifact(id uint) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type: state.ModifyArtifact,
			ID:   id,
		}
	}
}
