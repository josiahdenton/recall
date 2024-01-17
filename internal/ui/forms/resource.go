package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	name = iota
	source
)

var (
	selectedOptionStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF")).Width(25)
	defaultOptionStyle  = styles.SecondaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(25)
	paginationStyle     = list.DefaultStyles().PaginationStyle
	fadedTitleStyle     = styles.SecondaryGray.Copy()
)

type ResourceFormMsg struct {
	Resource domain.Resource
}

type ResourceModel struct {
	inputs        []textinput.Model
	options       []createResourceOption
	selectFrom    list.Model
	existing      list.Model
	existingReady bool
	choice        domain.ResourceType
	active        int
	status        string
}

type createResourceOption struct {
	Title    string
	AttachBy attachMethod
}

func (r *createResourceOption) FilterValue() string {
	return ""
}

func NewResourceForm() ResourceModel {
	inputName := textinput.New()
	inputName.Focus()
	inputName.Width = 60
	inputName.CharLimit = 60
	inputName.Prompt = "Name: "
	inputName.PromptStyle = formLabelStyle
	inputName.Placeholder = "..."

	inputName.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputSource := textinput.New()
	inputSource.Width = 60
	inputSource.CharLimit = 300
	inputSource.Prompt = "Source: "
	inputSource.PromptStyle = formLabelStyle
	inputSource.Placeholder = "..."

	inputSource.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 2)
	inputs[name] = inputName
	inputs[source] = inputSource

	options := make([]createResourceOption, 2)
	options[0] = createResourceOption{
		Title:    "New",
		AttachBy: newItem,
	}
	options[1] = createResourceOption{
		Title:    "Existing",
		AttachBy: existingItem,
	}
	items := make([]list.Item, len(options))
	for i := range options {
		item := &options[i]
		items[i] = item
	}

	selectFrom := list.New(items, createResourceOptionDelegate{}, 50, 20)
	selectFrom.Title = "Resource AttachBy"
	selectFrom.SetShowStatusBar(false)
	selectFrom.SetFilteringEnabled(false)
	selectFrom.Styles.PaginationStyle = paginationStyle
	selectFrom.Styles.Title = fadedTitleStyle
	selectFrom.SetShowHelp(false)
	selectFrom.KeyMap.Quit.Unbind()

	return ResourceModel{
		inputs:     inputs,
		options:    options,
		choice:     domain.NoneType,
		selectFrom: selectFrom,
	}
}

func (m ResourceModel) Init() tea.Cmd {
	return nil
}

func (m ResourceModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Resource"))
	b.WriteString("\n\n")
	if m.choice == newItem {
		b.WriteString(m.inputs[name].View())
		b.WriteString("\n")
		b.WriteString(m.inputs[source].View())
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(m.status))
	} else if m.choice == existingItem && m.existingReady {
		b.WriteString(m.existing.View())
	} else if m.choice == none {
		b.WriteString(m.selectFrom.View())
	}
	return b.String()
}

func (m ResourceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case shared.LoadedStateMsg:
		resources := msg.State.([]domain.Resource)
		m.existing = list.New(resourcesToItemList(resources), resourceDelegate{}, 50, 10)
		m.existing.Title = "attach one of the following types"
		m.existing.Styles.PaginationStyle = paginationStyle
		m.existing.Styles.Title = fadedTitleStyle
		m.existing.SetShowHelp(false)
		m.existing.KeyMap.Quit.Unbind()
		m.existingReady = true
	}

	if m.choice == newItem {
		m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
		cmds = append(cmds, cmd)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.active == name {
					m.inputs[m.active%len(m.inputs)].Blur()
					m.active++
					m.inputs[m.active%len(m.inputs)].Focus()
					break
				}

				// TODO fix the <nil>
				if m.inputs[name].Err != nil || m.inputs[source].Err != nil {
					m.status = fmt.Sprintf("%v, %v", m.inputs[name].Err, m.inputs[source].Err)
				} else {
					cmds = append(cmds, addResourceToTask(domain.Resource{
						Name:   m.inputs[name].Value(),
						Source: m.inputs[source].Value(),
						Type:   domain.WebResource,
					}))

					m.inputs[name].Reset()
					m.inputs[name].Focus()
					m.inputs[source].Reset()
					m.inputs[source].Blur()
					m.active = name
					m.choice = domain.NoneType
				}
			case tea.KeyTab:
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
			case tea.KeyShiftTab:
				if m.active > 0 {
					m.inputs[m.active%len(m.inputs)].Blur()
					m.active--
					m.inputs[m.active%len(m.inputs)].Focus()
				}
			}
			if len(m.inputs[name].Value()) > 0 || len(m.inputs[source].Value()) > 0 {
				m.status = ""
			}
		}
	} else if m.choice == existingItem && m.existingReady {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				selected := m.existing.SelectedItem().(*domain.Resource)
				cmds = append(cmds, addResourceToTask(*selected))
			}
		}
		m.existing, cmd = m.existing.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.choice == none {
		m.selectFrom, cmd = m.selectFrom.Update(msg)
		cmds = append(cmds, cmd)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				choice := m.selectFrom.SelectedItem().(*createResourceOption)
				m.choice = choice.AttachBy
				if m.choice == existingItem {
					cmds = append(cmds, shared.RequestState(shared.LoadResource, 0))
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func addResourceToTask(resource domain.Resource) tea.Cmd {
	return func() tea.Msg {
		return ResourceFormMsg{
			Resource: resource,
		}
	}
}

func resourcesToItemList(resources []domain.Resource) []list.Item {
	items := make([]list.Item, len(resources))
	for i := range resources {
		item := &resources[i]
		items[i] = item
	}
	return items
}
