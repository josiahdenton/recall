package config

import (
	"fmt"
	"log"
	"os"
)

func SetupDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get home dir: %v", err)
		os.Exit(1)
	}

	parentPath := fmt.Sprintf("%s/%s", home, ".recall")
	// check if Dir exists, if not, create it
	fi, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(parentPath, 0775)
		if err != nil {
			log.Printf("failed to mkdir for recall: %v", err)
			os.Exit(1)
		}
	} else if err != nil {
		log.Printf("failed to stat recall dir: %v", err)
		os.Exit(1)
	}

	if !fi.IsDir() {
		log.Printf(".recall exists and is not a dir")
		os.Exit(1)
	}

	return parentPath
}
