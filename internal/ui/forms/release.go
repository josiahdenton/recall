package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

const (
	releaseChangeUrl = iota
	releaseOwner
	releaseDate
)

type editReleaseMsg struct {
	release *domain.Release
}

func EditRelease(release *domain.Release) tea.Cmd {
	return func() tea.Msg {
		return editReleaseMsg{release: release}
	}
}

type ReleaseFormMsg struct {
	Release domain.Release
	Edit    bool
}

func NewReleaseForm() ReleaseFormModel {
	inputChangeUrl := textinput.New()
	inputChangeUrl.Focus()
	inputChangeUrl.Width = 60
	inputChangeUrl.CharLimit = 120
	inputChangeUrl.Prompt = "Change Source: "
	inputChangeUrl.PromptStyle = styles.FormLabelStyle
	inputChangeUrl.Placeholder = "https://www.example.com/..."

	inputChangeUrl.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("release change url missing")
		}
		return nil
	}

	inputOwner := textinput.New()
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

	inputs := make([]textinput.Model, 3)
	inputs[releaseChangeUrl] = inputChangeUrl
	inputs[releaseOwner] = inputOwner
	inputs[releaseDate] = inputDate

	return ReleaseFormModel{
		inputs:  inputs,
		release: &domain.Release{},
	}
}

type ReleaseFormModel struct {
	inputs  []textinput.Model
	release *domain.Release
	active  int
}

func (m ReleaseFormModel) Init() tea.Cmd {
	return nil
}

func (m ReleaseFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Release"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[releaseChangeUrl].View())
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
	case editReleaseMsg:
		m.release = msg.release
		m.inputs[releaseChangeUrl].SetValue(m.release.ReleaseChange)
		m.inputs[releaseOwner].SetValue(m.release.Owner)
		m.inputs[releaseDate].SetValue(formatDate(m.release.Date))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.active = releaseChangeUrl
			m.inputs[releaseChangeUrl].Reset()
			m.inputs[releaseChangeUrl].Focus()
			m.inputs[releaseOwner].Reset()
			m.inputs[releaseOwner].Blur()
			m.inputs[releaseDate].Reset()
			m.inputs[releaseDate].Blur()
			m.release = &domain.Release{}
		case tea.KeyEnter:
			if m.active < releaseDate {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			if cmd := validateForm(m.inputs[releaseOwner].Err, m.inputs[releaseDate].Err, m.inputs[releaseChangeUrl].Err); cmd != nil {
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			releaseDateParsed, cmd := parseDate(m.inputs[releaseDate].Value())
			if cmd != nil {
				return m, tea.Batch(cmds...)
			}
			m.release.Date = releaseDateParsed
			m.release.Owner = m.inputs[releaseOwner].Value()
			m.release.ReleaseChange = m.inputs[releaseChangeUrl].Value()
			cmds = append(cmds, addRelease(*m.release))
			m.release = &domain.Release{}
			// Reset Form
			m.active = releaseChangeUrl
			m.inputs[releaseChangeUrl].Reset()
			m.inputs[releaseChangeUrl].Focus()
			m.inputs[releaseOwner].Reset()
			m.inputs[releaseOwner].Blur()
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

func addRelease(release domain.Release) tea.Cmd {
	return func() tea.Msg {
		return ReleaseFormMsg{Release: release, Edit: release.ID != 0}
	}
}
