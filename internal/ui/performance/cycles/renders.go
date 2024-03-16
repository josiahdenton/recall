package cycles

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	activeCycleStyle  = styles.SecondaryColorStyle.Copy().PaddingLeft(2)
	defaultCycleStyle = styles.PrimaryGrayStyle.Copy().PaddingLeft(2)
	cursorStyle       = styles.PrimaryColorStyle.Copy().PaddingLeft(2)
)

func renderCycleOption(cycle *domain.Cycle, selected bool) string {
	selectedMarker := " "
	style := defaultCycleStyle
	if selected {
		selectedMarker = ">"
	}

	if cycle.Active {
		style = activeCycleStyle
	}

	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), style.Render(cycle.Title))
}
