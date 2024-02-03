package jsonified

import "github.com/josiahdenton/recall/internal/domain"

func LoadKeybindings(path string) (domain.Keybindings, error) {
	return domain.Keybindings{}, nil
}

func SaveKeybindings(keys domain.Keybindings) error {
	return nil
}

func KeybindingsExist() {}

func SetDefaultKeybindings() {}
