package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"strings"
)

type LinkFormMsg struct {
	Zettel domain.Zettel
}

const (
	none = iota
	newZettel
	existingZettel
)

type linkType = int

type linkZettelOption struct {
	DisplayName string
	TypeOption  linkType
}

func (r *linkZettelOption) FilterValue() string {
	return ""
}

func NewLinkForm() LinkFormModel {
	inputName := textinput.New()
	inputName.Focus()
	inputName.Width = 60
	inputName.CharLimit = 60
	inputName.Prompt = "Name: "
	inputName.PromptStyle = formLabelStyle
	inputName.Placeholder = "..."

	inputName.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("zettel name missing")
		}
		return nil
	}

	options := []linkZettelOption{
		{
			DisplayName: "New Zettel",
			TypeOption:  newZettel,
		},
		{
			DisplayName: "Existing Zettel",
			TypeOption:  existingZettel,
		},
	}

	items := make([]list.Item, len(options))
	for i := range options {
		item := &options[i]
		items[i] = item
	}

	createOptions := list.New(items, createZettelOptionDelegate{}, 50, 6)
	createOptions.Title = "attach one of the following types"
	createOptions.SetShowStatusBar(false)
	createOptions.SetFilteringEnabled(false)
	createOptions.Styles.PaginationStyle = paginationStyle
	createOptions.Styles.Title = fadedTitleStyle
	createOptions.SetShowHelp(false)
	createOptions.KeyMap.Quit.Unbind()

	return LinkFormModel{
		nameInput:     inputName,
		choice:        linkZettelOption{},
		createOptions: createOptions,
	}

}

type LinkFormModel struct {
	nameInput     textinput.Model
	choice        linkZettelOption
	createOptions list.Model
	existing      list.Model // populate this depending on choice
	status        string
	existingReady bool
}

func (m LinkFormModel) Init() tea.Cmd {
	return nil
}

func (m LinkFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Link a Zettel"))
	b.WriteString("\n")
	b.WriteString(m.createOptions.View())
	b.WriteString("\n")
	b.WriteString(titleStyle.Render(fmt.Sprintf("linking a zettel of type: %s", m.choice.DisplayName)))
	if m.choice.TypeOption == existingZettel && m.existingReady {
		b.WriteString(m.existing.View())
	} else if m.choice.TypeOption == newZettel {
		b.WriteString(m.nameInput.View())
	}
	b.WriteString("\n")
	b.WriteString(errorStyle.Render(m.status))

	return b.String()
}

func (m LinkFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.choice.TypeOption == none {
		m.createOptions, cmd = m.createOptions.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case shared.LoadedStateMsg:
		zettels := msg.State.([]domain.Zettel)
		m.existing = list.New(toItemList(zettels), createZettelOptionDelegate{}, 50, 10)
		m.existing.Title = "attach one of the following types"
		m.existing.SetShowStatusBar(false)
		m.existing.SetFilteringEnabled(false)
		m.existing.Styles.PaginationStyle = paginationStyle
		m.existing.Styles.Title = fadedTitleStyle
		m.existing.SetShowHelp(false)
		m.existing.KeyMap.Quit.Unbind()
		m.existingReady = true
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			// reset
			m.choice = linkZettelOption{}
			m.nameInput.Reset()
		case tea.KeyEnter:
			// depends on which model is active...
			switch m.choice.TypeOption {
			case none:
				selected := m.createOptions.SelectedItem().(*linkZettelOption)
				m.choice = *selected
				if m.choice.TypeOption == existingZettel {
					cmds = append(cmds, shared.RequestState(shared.LoadZettel, 0))
				}
			case newZettel:
				if m.nameInput.Err != nil {
					m.status = "missing name for new zettel"
					// TODO - clear status
					break
				}
				cmds = append(cmds, linkZettel(domain.Zettel{Name: m.nameInput.Value()}))
			case existingZettel:
				selected := m.existing.SelectedItem().(*domain.Zettel)
				cmds = append(cmds, linkZettel(*selected))
			}
		}
	}

	if m.choice.TypeOption == newZettel {
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.choice.TypeOption == existingZettel {
		m.existing, cmd = m.existing.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func toItemList(zettels []domain.Zettel) []list.Item {
	items := make([]list.Item, len(zettels))
	for i := range zettels {
		item := &zettels[i]
		items[i] = item
	}
	return items
}

func linkZettel(zettel domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return LinkFormMsg{Zettel: zettel}
	}
}