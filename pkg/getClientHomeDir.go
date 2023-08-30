package pkg

import (
	"os"
	"path/filepath"
	"runtime"
)

//get client's home directory
//get file path to save the file into the local machine
//logFilepath is created from the client's home directory(Linux or Mac) and logFileName, for Windows, the Documents folder is included.
func (app *Application) getClientHomeDir(logFileName string) (string, error) {
	//get client's home directory
	var homeDir string
	var err error
	if runtime.GOOS == "windows" {
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		homeDir = filepath.Join(homeDrive, homePath)

	} else {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			app.App.ErrorLog.Printf("Failed to get user home directory: %v", err)
			return "", err
		}
	}

	//get file path to save the file into the local machine
	//logFilepath is created from the client's home directory(Linux or Mac) and logFileName, for Windows, the Documents folder is included.
	var localFilePath string
	if runtime.GOOS == "windows" {
		downloadFolder := filepath.Join(homeDir, "Documents")
		localFilePath = filepath.Join(downloadFolder, logFileName)
	} else {
		localFilePath = filepath.Join(homeDir, logFileName)
	}
	app.App.InfoLog.Println("Location of file on local directory", localFilePath)
	return localFilePath, err
}
