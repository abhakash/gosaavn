package util

import (
	"github.com/abhakash/gosaavn/internal/logging"
	"os"
	"path/filepath"
)

func CreateDirectory(locationDirectory string) error {
	if !doesDirectoryExist(locationDirectory) {
		err := os.MkdirAll(locationDirectory, 0755)
		if err != nil {
			logging.Log.Error("Error creating directory:", err)
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func CreateFile(locationDirectory string, fileName string) (*os.File, error) {
	absoluteDirPath, err := filepath.Abs(locationDirectory)
	if err != nil {
		return nil, err
	}
	absoluteFilePath := filepath.Join(absoluteDirPath, fileName)
	file, err := os.Create(absoluteFilePath)
	if err != nil {
		logging.Log.Errorf("Error creating file: %s", err)
		return nil, err
	} else {
		return file, err
	}
}

func doesDirectoryExist(locationDirectory string) bool {
	_, err := os.Stat(locationDirectory)
	return os.IsExist(err)
}
