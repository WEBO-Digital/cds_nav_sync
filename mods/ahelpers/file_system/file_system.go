package filesystem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Save(path string, fileName string, data interface{}) error {
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

	//Specify the file path
	destinationPath := fmt.Sprintf("%s%s.json", *filePath, fileName)

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

func ReadFile(path string, fileName string) ([]byte, error) {
	// Get the current working directory
	currentDir, err := GetCurrentWorkingDirectory()
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s%s%s", currentDir, path, fileName)

	// Read your JSON data from a file
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading JSON data:", err)
	}
	return jsonData, nil
}

func GetAllFiles(path string) ([]string, error) {
	// Get the current working directory
	currentDir, err := GetCurrentWorkingDirectory()
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s%s", currentDir, path)

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
	currentDir, err := GetCurrentWorkingDirectory()
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

func Append(path string, fileName string, data string) error {
	//Create folder if it does not exists
	filePath, err := createDirectoryIfNotExists(path)
	if err != nil {
		return err
	}

	//Specify the file path
	destinationPath := fmt.Sprintf("%s%s", *filePath, fileName)

	//Create File if it does not exists
	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		//Create file
		_, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(destinationPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	// Append data to the file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func CleanAndSave(path string, fileName string, data string) error {
	//Create folder if it does not exists
	filePath, err := createDirectoryIfNotExists(path)
	if err != nil {
		return err
	}

	//Specify the file path
	destinationPath := fmt.Sprintf("%s%s", *filePath, fileName)

	//Create File if it does not exists
	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		//Create file
		_, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(destinationPath, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	// Append data to the file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func GetCurrentWorkingDirectory() (string, error) {
	currentDir, err := os.Getwd()
	return currentDir, err
}

func createDirectoryIfNotExists(path string) (*string, error) {
	// Get the current working directory
	currentDir, err := GetCurrentWorkingDirectory()
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
