package filesystem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Save(path string, data interface{}) error {
	//Marshal data to JSON
	jsonData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	//Create folder if it does not exists
	filePath, err := createDirectoryIfNotExists(path)
	if err != nil {
		return err
	}

	//get current timestamp
	timestamp := time.Now().Format("2006-01-02T15-04-05")

	//Specify the file path
	destinationPath := fmt.Sprintf("%s%s.json", *filePath, timestamp)

	//Create File if it does not exists
	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		//Create file
		_, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(destinationPath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetAllFiles(path string) ([]string, error) {
	// Get the current working directory
	currentDir, err := getCurrentWorkingDirectory()
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s", currentDir, path)

	// Walk the directory and add file names to the slice
	var fileNames []string
	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileNames = append(fileNames, info.Name())
		}
		return nil
	})
	return fileNames, nil
}

func MoveFile(sourceFileName string, sourceDirectory string, destinationDirectory string) error {
	//Create folder if it does not exists
	_, err := createDirectoryIfNotExists(destinationDirectory)
	if err != nil {
		return err
	}

	// Get the current working directory
	currentDir, err := getCurrentWorkingDirectory()
	if err != nil {
		return err
	}

	//Construct the Source path
	sourcePath := filepath.Join(currentDir, sourceDirectory)
	sourcePath = filepath.Join(sourcePath, sourceFileName)

	//Construct the Destination path
	destinationPath := filepath.Join(currentDir, destinationDirectory)
	destinationPath = filepath.Join(destinationPath, sourceFileName)

	//Move the file
	err = os.Rename(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	return nil
}

func getCurrentWorkingDirectory() (string, error) {
	currentDir, err := os.Getwd()
	return currentDir, err
}

func createDirectoryIfNotExists(path string) (*string, error) {
	// Get the current working directory
	currentDir, err := getCurrentWorkingDirectory()
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s%s", currentDir, path)

	//Give all permissions
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return nil, err
	}
	return &filePath, nil
}
