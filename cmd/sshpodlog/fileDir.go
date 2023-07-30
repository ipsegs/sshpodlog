package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func (app *Application) fileDir(logFileName string) (string, error) {
	//get home directory of the local server
	var homeDir string
	var err error
	if runtime.GOOS == "windows" {
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		homeDir = filepath.Join(homeDrive, homePath)

	} else {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			app.ErrorLog.Printf("Failed to get user home directory: %v", err)
			return "", err
		}
	}

	//get file path to save the file into the local machine
	var localFilePath string

	if runtime.GOOS == "windows" {
		downloadFolder := filepath.Join(homeDir, "Downloads")
		localFilePath = filepath.Join(downloadFolder, logFileName)
	} else {
		localFilePath = filepath.Join(homeDir, logFileName)
	}
	app.InfoLog.Println("Location of file on local directory", localFilePath)
	return localFilePath, err
}
