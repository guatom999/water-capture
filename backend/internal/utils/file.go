package utils

import (
	"fmt"
	"log"
	"os"
)

func DeleteFile(filePath string) error {

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Error File IsNotExist")
			return fmt.Errorf("error file is not exist")
		}
		return fmt.Errorf("error failed to delete file %v", err.Error())
	}

	return nil
}
