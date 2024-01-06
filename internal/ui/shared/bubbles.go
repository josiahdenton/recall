package shared

import (
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

func VerticalOptions(options []string, cursor int) string {
	var b strings.Builder
	for i, option := range options {
		if i == cursor {
			b.WriteString(styles.PrimaryColor.PaddingRight(1).Render(option))
		} else {
			b.WriteString(styles.PrimaryGray.PaddingRight(1).Render(option))
		}
	}
	return b.String()
}
