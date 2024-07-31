package internal

import (
	"encoding/json"
	"os"
)

func SaveObject(obj interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
