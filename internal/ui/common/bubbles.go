package common

import (
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

func VerticalOptions(options []string, cursor int) string {
	var b strings.Builder
	for i, option := range options {
		if i == cursor {
			b.WriteString(styles.PrimaryColor.Render(option))
		} else {
			b.WriteString(styles.PrimaryGray.Render(option))
		}
	}
	return b.String()
}
