package keyset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/josiahdenton/recall/internal/domain"
)

const keybindFileName = "keybindings.json"

func New(recallFolder string) Handler {
	filePath := fmt.Sprintf("%s/%s", recallFolder, keybindFileName)
	return Handler{keybindingsFileLocation: filePath}
}

type Handler struct {
	keybindingsFileLocation string
}

func (h Handler) Load() (domain.Keybindings, error) {
	f, err := os.Open(h.keybindingsFileLocation)
	if os.IsNotExist(err) {
		log.Printf("failed to ")
		f, err = os.Create(h.keybindingsFileLocation)
		if err != nil {
			log.Printf("%v", err)
			return domain.Keybindings{}, err
		}
	} else if err != nil {
		log.Printf("%v", err)
		return domain.Keybindings{}, err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, f)

	var keybinds domain.Keybindings
	if n == 0 {
		keybinds = domain.DefaultKeybindings()
	} else {
		err = json.Unmarshal(buf.Bytes(), &keybinds)
		if err != nil {
			log.Printf("%v", err)
			return domain.Keybindings{}, err
		}
	}

	return keybinds, err
}

func (h *Handler) Save(keys domain.Keybindings) error {
	bytes, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	err = os.WriteFile(h.keybindingsFileLocation, bytes, 0666)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}
