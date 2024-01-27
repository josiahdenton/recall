package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"
)

const (
	releaseOwner = iota
	releaseDate
)

type ReleaseFormMsg struct {
	Release domain.Release
}

func NewReleaseForm() ReleaseFormModel {
	inputOwner := textinput.New()
	inputOwner.Focus()
	inputOwner.Width = 60
	inputOwner.CharLimit = 60
	inputOwner.Prompt = "Owner: "
	inputOwner.PromptStyle = styles.FormLabelStyle
	inputOwner.Placeholder = "..."

	inputOwner.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("release owner missing")
		}
		return nil
	}

	inputDate := textinput.New()
	inputDate.Width = 60
	inputDate.CharLimit = 120
	inputDate.Prompt = "Release Date: "
	inputDate.PromptStyle = styles.FormLabelStyle
	inputDate.Placeholder = "Jan 5, 2013"

	inputDate.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("release date missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 2)
	inputs[releaseOwner] = inputOwner
	inputs[releaseDate] = inputDate

	return ReleaseFormModel{
		inputs: inputs,
	}
}

type ReleaseFormModel struct {
	inputs []textinput.Model
	active int
}

func (m ReleaseFormModel) Init() tea.Cmd {
	return nil
}

func (m ReleaseFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Release"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[releaseOwner].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[releaseDate].View())
	return b.String()
}

func (m ReleaseFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active < releaseDate {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			if cmd := validateReleaseForm(m.inputs[releaseOwner].Err, m.inputs[releaseDate].Err); cmd != nil {
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			releaseDateParsed, cmd := parseDate(m.inputs[releaseDate].Value())
			if cmd != nil {
				return m, tea.Batch(cmds...)
			}
			cmds = append(cmds, addRelease(domain.Release{
				Date:  releaseDateParsed,
				Owner: m.inputs[releaseOwner].Value(),
			}))
			// Reset Form
			m.inputs[releaseOwner].Reset()
			m.inputs[releaseOwner].Focus()
			m.inputs[releaseDate].Reset()
			m.inputs[releaseDate].Blur()
		case tea.KeyTab:
			m.inputs[m.active%len(m.inputs)].Blur()
			m.active++
			m.inputs[m.active%len(m.inputs)].Focus()
		}
	}

	m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func validateReleaseForm(a, b error) tea.Cmd {
	if a != nil {
		return toast.ShowToast(fmt.Sprintf("%v", a))
	}
	if b != nil {
		return toast.ShowToast(fmt.Sprintf("%v", b))
	}
	return nil
}

func addRelease(release domain.Release) tea.Cmd {
	return func() tea.Msg {
		return ReleaseFormMsg{Release: release}
	}
}
