package utils

import (
	"encoding/json"
	"log"
	"os"
)

func SaveJSON(data interface{}, fullpath string) error {
	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	if err := enc.Encode(data); err != nil {
		return err
	}

	log.Printf("Saved JSON to %s", fullpath)
	return nil
}
