package artifact

import (
	"fmt"
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

func New(keyBinds domain.Keybindings) Model {
	models := make([]tea.Model, 3)
	models[releases] = forms.NewReleaseForm()
	models[resources] = forms.NewResourceForm()
	models[header] = forms.NewArtifactForm()
	return Model{
		keyBinds: keyBinds,
		forms:    models,
		active:   header,
	}
}

type Model struct {
	keyBinds   domain.Keybindings
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
		if !msg.Edit {
			m.artifact.Releases = append(m.artifact.Releases, msg.Release)
			m.releases.InsertItem(len(m.artifact.Releases), &m.artifact.Releases[len(m.artifact.Releases)-1])
			cmds = append(cmds, updateArtifact(*m.artifact))
		} else {
			m.showForm = false
			cmds = append(cmds, updateRelease(msg.Release))
		}

	case forms.ResourceFormMsg:
		// attach resource
		m.artifact.Resources = append(m.artifact.Resources, msg.Resource)
		m.resources.InsertItem(len(m.artifact.Resources), &m.artifact.Resources[len(m.artifact.Resources)-1])
		cmds = append(cmds, updateArtifact(*m.artifact))
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
		case tea.KeySpace:
			if m.active == releases && len(m.releases.Items()) > 0 {
				selected := m.releases.SelectedItem().(*domain.Release)
				switch selected.Outcome {
				case domain.AwaitingRelease:
					selected.Outcome = domain.SuccessfulRelease
				case domain.SuccessfulRelease:
					selected.Outcome = domain.FailedRelease
				case domain.FailedRelease:
					selected.Outcome = domain.AwaitingRelease
				}
				cmds = append(cmds, updateRelease(*selected))
			}

		case tea.KeyEnter:
			switch m.active {
			case releases:
				selected, ok := m.releases.SelectedItem().(*domain.Release)
				if ok {
					selected.OpenChange()
					cmds = append(cmds, toast.ShowToast("opened release change!", toast.Info))
				}
			case resources:
				if len(m.resources.Items()) > 0 {
					selected := m.resources.SelectedItem().(*domain.Resource)
					selected.Open()
					cmds = append(cmds, toast.ShowToast("opened resource!", toast.Info))
				}
			case header:
				if editor := m.artifact.Open(); editor != nil {
					cmd = tea.ExecProcess(editor, func(err error) tea.Msg {
						if err != nil {
							return toast.ShowToastMsg{
								Message: fmt.Sprintf("failed to open artifact %v", err),
								Toast:   toast.Warn,
							}
						}
						return nil
					})
					cmds = append(cmds, cmd, toast.ShowToast("opening artifact!", toast.Info))
				} else {
					cmds = append(cmds, toast.ShowToast("failed to open artifact", toast.Warn))
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
		case "e":
			if !m.showForm && m.active == header {
				m.showForm = true
				cmds = append(cmds, forms.EditArtifact(m.artifact))
			} else if !m.showForm && m.active == releases {
				selected := m.releases.SelectedItem().(*domain.Release)
				cmds = append(cmds, forms.EditRelease(selected))
				m.showForm = true
			}
		case "d":
			if m.active == releases && len(m.releases.Items()) > 0 {
				selected := m.releases.SelectedItem().(*domain.Release)
				index := m.releases.Index()
				m.artifact.Releases = append(m.artifact.Releases[:index], m.artifact.Releases[index+1:]...)
				m.releases.SetItems(releasesToItemList(m.artifact.Releases))
				cmds = append(cmds, removeReleaseFromArtifact(m.artifact, selected))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func removeReleaseFromArtifact(parent *domain.Artifact, child *domain.Release) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type:   state.ModifyRelease,
			Parent: parent,
			Child:  child,
		}
	}
}

func updateArtifact(artifact domain.Artifact) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: artifact,
			Type:   state.ModifyArtifact,
		}
	}
}

func updateRelease(release domain.Release) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: release,
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
