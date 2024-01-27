package artifact

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/toast"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

const (
	releases = iota
	resources
	header
)

func New() Model {
	models := make([]tea.Model, 3)
	models[releases] = forms.NewReleaseForm()
	models[resources] = forms.NewResourceForm()
	models[header] = forms.NewArtifactForm()
	return Model{
		forms:  models,
		active: header,
	}
}

type Model struct {
	artifact   *domain.Artifact
	releases   list.Model
	resources  list.Model
	active     int
	showForm   bool
	forms      []tea.Model
	activeForm int
	ready      bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		m.artifact = msg.State.(*domain.Artifact)
		// releases
		m.releases = list.New(releasesToItemList(m.artifact.Releases), releaseDelegate{}, 80, 20)
		m.releases.SetShowStatusBar(false)
		m.releases.SetShowHelp(false)
		m.releases.SetFilteringEnabled(false)
		m.releases.Title = "Releases"
		m.releases.Styles.PaginationStyle = paginationStyle
		m.releases.Styles.Title = titleStyle
		m.releases.KeyMap.Quit.Unbind()
		// resources
		m.resources = list.New(resourcesToItemList(m.artifact.Resources), resourceDelegate{}, 80, 8)
		m.resources.SetShowStatusBar(false)
		m.resources.SetShowHelp(false)
		m.resources.SetFilteringEnabled(false)
		m.resources.Title = "Resources"
		m.resources.Styles.PaginationStyle = paginationStyle
		m.resources.Styles.Title = titleStyle
		m.resources.KeyMap.Quit.Unbind()
		m.ready = true
	case state.SaveStateMsg:
		m.showForm = false
	case forms.ReleaseFormMsg:
		// attach release
	case forms.ResourceFormMsg:
		// attach resource
		m.artifact.Resources = append(m.artifact.Resources, msg.Resource)
		m.resources.InsertItem(len(m.artifact.Resources), &m.artifact.Resources[len(m.artifact.Resources)-1])
		cmds = append(cmds, updateArtifact(m.artifact))
	}

	if !m.ready {
		return m, tea.Batch(cmds...)
	}

	if m.showForm {
		m.forms[m.active%len(m.forms)], cmd = m.forms[m.active%len(m.forms)].Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.showForm {
				m.showForm = false
			} else {
				cmds = append(cmds, router.GotoPreviousPage())
				m.active = header
			}
		}
	}

	if m.showForm {
		return m, tea.Batch(cmds...)
	}

	if m.active == releases {
		m.releases, cmd = m.releases.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.active == resources {
		m.resources, cmd = m.resources.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// depends on the thing...
			switch m.active {
			case releases:
				if len(m.releases.Items()) > 0 {
					selected := m.releases.SelectedItem().(*domain.Release)
					switch selected.Outcome {
					case domain.AwaitingRelease:
						selected.Outcome = domain.SuccessfulRelease
					case domain.SuccessfulRelease:
						selected.Outcome = domain.FailedRelease
					case domain.FailedRelease:
						selected.Outcome = domain.AwaitingRelease
					}
					cmds = append(cmds, updateRelease(selected))
				}
			case resources:
				if len(m.resources.Items()) > 0 {
					selected := m.resources.SelectedItem().(*domain.Resource)
					selected.Open()
					cmds = append(cmds, toast.ShowToast("opened resource!"))
				}
			case header:
				if !m.showForm {
					m.showForm = true
					cmds = append(cmds, forms.EditArtifact(m.artifact))
				}
			}
		}

		switch msg.String() {
		case "a":
			if m.active == releases || m.active == resources {
				m.showForm = true
			}
		case "d":
			// TODO - support deleting releases / resources
			//if m.active == links {
			//	selected := m.links.SelectedItem().(*domain.Zettel)
			//	m.links.RemoveItem(m.links.Index())
			//	cmds = append(cmds, unlinkZettel(m.zettel, selected), toast.ShowToast("unlinked zettel!"))
			//} else if
		}
	}

	return m, tea.Batch(cmds...)
}

func updateArtifact(artifact *domain.Artifact) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *artifact,
			Type:   state.ModifyArtifact,
		}
	}
}

func updateRelease(release *domain.Release) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *release,
			Type:   state.ModifyRelease,
		}
	}
}

func releasesToItemList(releases []domain.Release) []list.Item {
	items := make([]list.Item, len(releases))
	for i := range releases {
		item := &releases[i]
		items[i] = item
	}
	return items
}

func resourcesToItemList(resources []domain.Resource) []list.Item {
	items := make([]list.Item, len(resources))
	for i := range resources {
		item := &resources[i]
		items[i] = item
	}
	return items
}
