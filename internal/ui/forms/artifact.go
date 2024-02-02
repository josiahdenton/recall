package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/state"
	"strings"
)

const (
	artifactName = iota
	artifactTags
	artifactEditor
	artifactPath
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

	inputEditor := textinput.New()
	inputEditor.Width = 60
	inputEditor.CharLimit = 60
	inputEditor.Prompt = "Editor: "
	inputEditor.PromptStyle = formLabelStyle
	inputEditor.Placeholder = "(goland)" // TODO - to support nvim, you would need to pause the app...

	// TODO - make these optional
	inputEditor.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("artifact editor missing")
		}
		return nil
	}

	inputPath := textinput.New()
	inputPath.Width = 60
	inputPath.CharLimit = 300
	inputPath.Prompt = "File Location: "
	inputPath.PromptStyle = formLabelStyle
	inputPath.Placeholder = "~/projects/foo/bar/baz"

	inputPath.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("artifact file path missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 4)
	inputs[artifactName] = inputName
	inputs[artifactTags] = inputTags
	inputs[artifactEditor] = inputEditor
	inputs[artifactPath] = inputPath

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
	b.WriteString("\n\n")
	b.WriteString(m.inputs[artifactEditor].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[artifactPath].View())
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
		m.inputs[artifactEditor].SetValue(m.artifact.Editor)
		m.inputs[artifactPath].SetValue(m.artifact.Path)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active < artifactPath {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			if cmd := validateForm(
				m.inputs[artifactName].Err,
				m.inputs[artifactEditor].Err,
				m.inputs[artifactPath].Err,
			); cmd != nil {
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			m.artifact.Name = m.inputs[artifactName].Value()
			m.artifact.Tags = m.inputs[artifactTags].Value()
			m.artifact.Editor = m.inputs[artifactEditor].Value()
			m.artifact.Path = m.inputs[artifactPath].Value()
			cmds = append(cmds, saveArtifact(m.artifact))
			// Reset the form
			m.inputs[artifactName].Reset()
			m.inputs[artifactName].Focus()
			m.inputs[artifactTags].Reset()
			m.inputs[artifactTags].Blur()
			m.inputs[artifactEditor].Reset()
			m.inputs[artifactEditor].Blur()
			m.inputs[artifactPath].Reset()
			m.inputs[artifactPath].Blur()
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

func saveArtifact(artifact *domain.Artifact) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *artifact,
			Type:   state.ModifyArtifact,
		}
	}
}
