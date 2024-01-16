package shared

import tea "github.com/charmbracelet/bubbletea"

// RequestState grabs state from the repository and
// sends back a LoadedStateMsg. If "ID" is 0, RequestState
// will give back an array of those items
func RequestState(loadType LoadType, id uint) tea.Cmd {
	return func() tea.Msg {
		return RequestStateMsg{
			Type: loadType,
			ID:   id,
		}
	}
}
