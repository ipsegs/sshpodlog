package pkg

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
)

func (app *Application) CopyFileFromRemoteToLocal(conn *ssh.Client, outFile string) error {

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		app.App.ErrorLog.Printf("Failed to create SFTP client connection: %v\n", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(outFile)
	if err != nil {
		app.App.ErrorLog.Printf("Failed to open remote file: %v\n", err)
		return err
	}

	defer remoteFile.Close()

	localFilePath, err := app.GetClientHomeDir(outFile)
	if err != nil {
		app.App.ErrorLog.Println("Failed to get client home directory:", err)
		return err
	}

	//create the file name in the local machine
	localFile, err := os.Create(localFilePath)
	if err != nil {
		app.App.ErrorLog.Println("Failed to create the local file:", err)
		return err
	}
	defer localFile.Close()

	fmt.Println("Copying logs to", localFilePath)
	remoteFileInfo, err := remoteFile.Stat()
	if err != nil {
		app.App.ErrorLog.Println("Unable to get file status:", err)
		return err
	}
	fileSize := remoteFileInfo.Size()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		defer wg.Done()
		// track copy progress
		bar := progressbar.DefaultBytes(fileSize, "Downloading")

		//startTime := time.Now()

		//copy the file from remote to local
		_, err = io.Copy(io.MultiWriter(localFile, bar), remoteFile)
		if err != nil {
			app.App.ErrorLog.Println("Error copying file:", err)
			return
		}
		//elapsedTime := time.Since(startTime)
		bar.Finish()
	}()
	wg.Wait()

	fmt.Printf("Copied %d kilobytes content.\n", fileSize/1024)

	return nil
}
