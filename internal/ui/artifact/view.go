package artifact

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"
)

var (
	paginationStyle     = list.DefaultStyles().PaginationStyle
	defaultSectionStyle = styles.SecondaryGray.Copy()
	activeSectionStyle  = styles.SecondaryColor.Copy()
	contextStyle        = styles.SecondaryGray.Copy()
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
	var b strings.Builder
	if m.showForm {
		b.WriteString(m.forms[m.active%len(m.forms)].View())
	} else {
		if m.active == header {
			b.WriteString(activeSectionStyle.Render(m.artifact.Name))
		} else {
			b.WriteString(defaultSectionStyle.Render(m.artifact.Name))
		}
		b.WriteString("\n")
		b.WriteString(contextStyle.Render("Tags "))
		b.WriteString(contextStyle.Render(m.artifact.Tags))
		b.WriteString("\n\n")
		b.WriteString(m.releases.View())
		b.WriteString("\n")
		b.WriteString(m.resources.View())
	}
	return styles.WindowStyle.Render(b.String())
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
		m.releases.Styles.Title = defaultSectionStyle
		m.releases.KeyMap.Quit.Unbind()
		// resources
		m.resources = list.New(resourcesToItemList(m.artifact.Resources), resourceDelegate{}, 80, 8)
		m.resources.SetShowStatusBar(false)
		m.resources.SetShowHelp(false)
		m.resources.SetFilteringEnabled(false)
		m.resources.Title = "Resources"
		m.resources.Styles.PaginationStyle = paginationStyle
		m.resources.Styles.Title = defaultSectionStyle
		m.resources.KeyMap.Quit.Unbind()
		m.ready = true
	case state.SaveStateMsg:
		m.showForm = false
	case forms.ReleaseFormMsg:
		// attach release
		m.artifact.Releases = append(m.artifact.Releases, msg.Release)
		m.releases.InsertItem(len(m.artifact.Releases), &m.artifact.Releases[len(m.artifact.Releases)-1])
		cmds = append(cmds, updateArtifact(m.artifact))
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
		case tea.KeyTab:
			m.active = nextSection(m.active)
			if m.active == releases {
				m.releases.Styles.Title = activeSectionStyle
			} else if m.active == resources {
				m.releases.Styles.Title = defaultSectionStyle
				m.resources.Styles.Title = activeSectionStyle
			} else {
				m.resources.Styles.Title = defaultSectionStyle
			}
		}

		switch msg.String() {
		case "a":
			if m.active == releases || m.active == resources {
				m.showForm = true
			}
		case "d":
			//TODO - support deleting releases / resources
			// will need to have my type modifications changed a bit.
			// should make an "Unlink"
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

func nextSection(section int) int {
	switch section {
	case header:
		return releases
	case releases:
		return resources
	case resources:
		return header
	}
	return header
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
