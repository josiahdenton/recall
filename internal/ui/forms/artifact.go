package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"
)

const (
	artifactName = iota
	artifactTags
)

type editArtifact struct {
	artifact *domain.Artifact
}

func EditArtifact(artifact *domain.Artifact) tea.Cmd {
	return func() tea.Msg {
		return editArtifact{artifact: artifact}
	}
}

func NewArtifactForm() ArtifactFormModel {
	inputName := textinput.New()
	inputName.Focus()
	inputName.Width = 60
	inputName.CharLimit = 300
	inputName.Prompt = "Name: "
	inputName.PromptStyle = formLabelStyle
	inputName.Placeholder = "..."

	inputName.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("artifact name missing")
		}
		return nil
	}

	inputTags := textinput.New()
	inputTags.Width = 60
	inputTags.CharLimit = 300
	inputTags.Prompt = "Tags: "
	inputTags.PromptStyle = formLabelStyle
	inputTags.Placeholder = "(comma seperated list - tags are for improving searching only)"

	inputs := make([]textinput.Model, 2)
	inputs[artifactName] = inputName
	inputs[artifactTags] = inputTags

	return ArtifactFormModel{
		inputs:   inputs,
		artifact: &domain.Artifact{},
	}
}

type ArtifactFormModel struct {
	inputs   []textinput.Model
	active   int
	artifact *domain.Artifact
}

func (m ArtifactFormModel) Init() tea.Cmd {
	return nil
}

func (m ArtifactFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Artifact"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[artifactName].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[artifactTags].View())
	return b.String()
}

func (m ArtifactFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case editArtifact:
		m.artifact = msg.artifact
		m.inputs[artifactName].SetValue(m.artifact.Name)
		m.inputs[artifactTags].SetValue(m.artifact.Tags)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active < artifactTags {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			if err := m.inputs[artifactName].Err; err != nil {
				cmds = append(cmds, toast.ShowToast(fmt.Sprintf("%v", err)))
				return m, tea.Batch(cmds...)
			}
			m.artifact.Name = m.inputs[artifactName].Value()
			m.artifact.Tags = m.inputs[artifactTags].Value()
			cmds = append(cmds, addArtifact(m.artifact))
			// Reset the form
			m.inputs[artifactName].Reset()
			m.inputs[artifactName].Focus()
			m.inputs[artifactTags].Reset()
			m.inputs[artifactTags].Blur()
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

func addArtifact(artifact *domain.Artifact) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *artifact,
			Type:   state.ModifyArtifact,
		}
	}
}
