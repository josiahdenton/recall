package zettel

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/adapters/editors"
	"github.com/josiahdenton/recall/internal/adapters/editors/nvim"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/zettel/forms"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	viewportStyle   = lipgloss.NewStyle().
			Padding(2).
			Width(80).
			Height(20).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3a3b5b"))
	alignCenterStyle    = lipgloss.NewStyle().Align(lipgloss.Center)
	activeViewportStyle = viewportStyle.Copy().BorderForeground(lipgloss.Color("#D120AF"))
	titleStyle          = styles.PrimaryColor.Copy().Align(lipgloss.Center)
	defaultLinksTitle   = styles.SecondaryGray.Copy()
	activeLinksTitle    = styles.PrimaryColor.Copy()
)

const (
	content = iota
	links
)

type section = int

func New() Model {
	vp := viewport.New(80, 20)
	vp.Style = activeViewportStyle

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(50),
	)
	if err != nil {
		log.Printf("failed to create glamour renderer: %v", err)
		os.Exit(1)
	}

	vp.SetContent("no content")

	return Model{
		form:     forms.NewZettelForm(),
		editor:   nvim.New(),
		vp:       vp,
		renderer: renderer,
	}
}

type Model struct {
	zettel     *domain.Zettel
	form       tea.Model
	links      list.Model
	showForm   bool
	ready      bool
	active     section
	editActive bool
	editor     editors.Editor // not needed anymore
	vp         viewport.Model
	renderer   *glamour.TermRenderer
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	// TODO - tie in glamour for displaying the content
	var b strings.Builder
	if m.showForm {
		b.WriteString(m.form.View())
	} else {
		b.WriteString(titleStyle.Render(m.zettel.Name))
		b.WriteString("\n")
		b.WriteString(alignCenterStyle.Render(m.vp.View()))
		b.WriteString("\n")
		b.WriteString(m.links.View())
	}
	return styles.WindowStyle.Render(b.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		zettel := msg.State.(*domain.Zettel)
		m.zettel = zettel
		m.links = list.New(toItemList(m.zettel.Links), zettelDelegate{}, 80, 30)
		m.links.Title = "Links"
		m.links.Styles.PaginationStyle = paginationStyle
		m.links.Styles.Title = defaultLinksTitle
		m.links.SetShowHelp(false)
		m.links.KeyMap.Quit.Unbind()
		cmds = append(cmds, loadContent(m.zettel))
	case zettelContentMsg:
		content, err := m.renderer.Render(msg.content)
		if err != nil {
			log.Printf("failed render content: %v", err)
		} else {
			m.vp.SetContent(content)
			m.ready = true
		}
	case forms.ZettelFormMsg:
		// save this zettel with new link...
		m.zettel.Links = append(m.zettel.Links, &msg.Zettel)
		m.links.InsertItem(len(m.zettel.Links), m.zettel.Links[len(m.zettel.Links)-1])
		cmds = append(cmds, modifyLinks(*m.zettel))
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && m.showForm {
			m.showForm = false
		} else if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		}
	}

	if !m.ready {
		return m, tea.Batch(cmds...)
	}

	if m.showForm {
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	if m.active == content {
		m.vp, cmd = m.vp.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.active == links {
		m.links, cmd = m.links.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.active = nextSection(m.active)
			if m.active == content {
				m.vp.Style = activeViewportStyle
				m.links.Styles.Title = defaultLinksTitle
			} else if m.active == links {
				m.vp.Style = viewportStyle
				m.links.Styles.Title = activeLinksTitle
			}
		case tea.KeyShiftTab:
			m.active = nextSection(m.active)
			if m.active == content {
				m.vp.Style = activeViewportStyle
				m.links.Styles.Title = defaultLinksTitle
			} else if m.active == links {
				m.vp.Style = viewportStyle
				m.links.Styles.Title = activeLinksTitle
			}
		case tea.KeyEnter:
			if m.active == content {
				cmds = append(cmds, editZettelContent(m.zettel))
			} else if m.active == links {
				selected := m.links.SelectedItem().(*domain.Zettel)
				cmds = append(cmds, router.GotoPage(domain.ZettelPage, selected.ID))
			}
		}

		switch msg.String() {
		case "a": // add zettel
			if m.active == links {
				m.showForm = true
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func editZettelContent(zettel *domain.Zettel) tea.Cmd {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("editZettelContent: failed to get user home dir: %v", err)
	}
	// need to make this a setting...
	path := fmt.Sprintf("%s/%s/%s", home, "recall-notes", zettel.ContentLocation)
	// TODO - add err handling
	cmd := exec.Command("nvim")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Args = append(cmd.Args, path)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return loadContent(zettel)
	})
}

type zettelContentMsg struct {
	content string
}

func loadContent(zettel *domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("editZettelContent: failed to get user home dir: %v", err)
			return zettelContentMsg{}
		}
		// need to make this a setting...
		path := fmt.Sprintf("%s/%s/%s", home, "recall-notes", zettel.ContentLocation)
		// what if the path does not exist? then the file needs to be created...
		var bytes []byte
		for {
			bytes, err = os.ReadFile(path)
			if os.IsNotExist(err) {
				f, err := os.Create(path)
				if err != nil {
					log.Printf("failed creating file (%s) for reason: %v", path, err)
				}
				f.Close()
			} else if err != nil {
				log.Printf("failed openinng file (%s) for reason: %v", path, err)
				return zettelContentMsg{}
			} else {
				break
			}
		}
		log.Printf("content: %v", string(bytes))
		return zettelContentMsg{content: string(bytes)}
	}
}

func modifyLinks(zettel domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return shared.SaveStateMsg{
			Update: zettel,
			Type:   shared.ModifyZettel,
		}
	}
}

func nextSection(section section) section {
	if section == content {
		return links
	}
	return content
}

func toItemList(links []*domain.Zettel) []list.Item {
	items := make([]list.Item, len(links))
	for i := range links {
		items[i] = links[i]
	}
	return items
}
