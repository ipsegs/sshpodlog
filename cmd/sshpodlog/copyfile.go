package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
)

func (app *Application) sftpClientCopy(conn *ssh.Client, logFileName string) error {

	//sftp client connections
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		app.ErrorLog.Println("Failed to create SFTP client:", err)
		return err
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(logFileName)
	if err != nil {
		app.ErrorLog.Println("Failed to open remote file:", err)
		return err
	}

	remoteFileInfo, err := remoteFile.Stat()
	if err != nil {
		app.ErrorLog.Println("Unable to get file size", err)
		return err
	}
	fileSize := remoteFileInfo.Size()

	localFilePath, err := app.fileDir(logFileName)
	if err != nil {
		app.ErrorLog.Println("Failed to create SFTP client:", err)
		return err
	}

	//create the file name in the local machine
	localFile, err := os.Create(localFilePath)
	if err != nil {
		app.ErrorLog.Println("Failed to create the local file:", err)
		return err
	}
	defer localFile.Close()

	bar := progressbar.DefaultBytes(fileSize, "copying to local")

	//copy the file from remote to local
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		app.ErrorLog.Println("Error copying file:", err)
		return err
	}

	bar.Finish()
	fmt.Printf("Copied %d kilobytes content.\n", fileSize/1024)

	return nil
}
