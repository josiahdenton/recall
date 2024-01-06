package cycles

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	activeCycleStyle  = styles.SecondaryColor.Copy().PaddingLeft(2)
	defaultCycleStyle = styles.PrimaryGray.Copy().PaddingLeft(2)
	cursorStyle       = styles.PrimaryColor.Copy().PaddingLeft(2)
)

func renderCycleOption(cycle *domain.Cycle, selected bool) string {
	selectedMarker := " "
	style := defaultCycleStyle
	if selected {
		selectedMarker = ">"
		style = activeCycleStyle
	}

	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), style.Render(cycle.Title))
}
