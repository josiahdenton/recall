package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"io"
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
	inputs     []textinput.Model
	options    []resourceTypeOption
	selectFrom list.Model
	choice     domain.ResourceType
	active     int
	status     string
}

type resourceTypeOption struct {
	Title string
	Type  domain.ResourceType
}

func (r *resourceTypeOption) FilterValue() string {
	return ""
}

type typeOptionDelegate struct{}

func (d typeOptionDelegate) Height() int  { return 1 }
func (d typeOptionDelegate) Spacing() int { return 1 }
func (d typeOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d typeOptionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	option, ok := item.(*resourceTypeOption)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderOption(option, index == m.Index()))
}

func renderOption(cycle *resourceTypeOption, selected bool) string {
	var s string
	if selected {
		s = selectedOptionStyle.Render(cycle.Title)
	} else {
		s = defaultOptionStyle.Render(cycle.Title)
	}
	return s
}

func NewStepResourceForm() ResourceModel {
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

	// TODO - for now, only web is supported
	options := make([]resourceTypeOption, 1)
	options[0] = resourceTypeOption{
		Title: "Web",
		Type:  domain.WebResource,
	}
	//options[1] = resourceTypeOption{
	//	Title: "Zettel",
	//	Type:  domain.ZettelResource,
	//}
	//options[2] = resourceTypeOption{
	//	Title: "File",
	//	Type:  domain.FilePathResource,
	//}
	items := make([]list.Item, 1)
	for i := range options {
		item := &options[i]
		items[i] = item
	}

	selectFrom := list.New(items, typeOptionDelegate{}, 50, 20)
	selectFrom.Title = "Resource Type"
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
	if m.choice != domain.NoneType {
		b.WriteString(m.inputs[name].View())
		b.WriteString("\n")
		b.WriteString(m.inputs[source].View())
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(m.status))
	} else {
		b.WriteString(m.selectFrom.View())
	}
	return b.String()
}

func (m ResourceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.choice != domain.NoneType {
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
					m.status = errorStyle.Render(fmt.Sprintf("%v, %v", m.inputs[name].Err, m.inputs[source].Err))
				} else {
					cmds = append(cmds, addResourceToTask(domain.Resource{
						Name:   m.inputs[name].Value(),
						Source: m.inputs[source].Value(),
						Type:   m.choice,
					}))

					m.inputs[name].Reset()
					m.inputs[source].Reset()
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
	} else {
		m.selectFrom, cmd = m.selectFrom.Update(msg)
		cmds = append(cmds, cmd)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				choice := m.selectFrom.SelectedItem().(*resourceTypeOption)
				m.choice = choice.Type
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
