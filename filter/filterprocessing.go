package filter

import (
	"bufio"
	"fmt"
	"regexp"
	"sync"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Match(conn *ssh.Client, logFileName, filter string) error {
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("Failed to create SFTP client connection: %v\n", err)
		return err
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(logFileName)
	if err != nil {
		fmt.Printf("Failed to open remote file using sftp: %v\n", err)
		return err
	}
	defer remoteFile.Close()

	scanner := bufio.NewScanner(remoteFile)
	pattern := regexp.MustCompile(filter)

	var wg sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			wg.Add(1)

			go func(line string) {
				defer wg.Done()
			// Print the matching line to the terminal
			fmt.Println(line)
			} (line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning the remote file: %v\n", err)
		return err
	}

	return err
}

// output filtered logs to a file
// package filter

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"regexp"

// 	"github.com/ipsegs/sshpodlog/pkg"
// 	"github.com/pkg/sftp"
// 	"golang.org/x/crypto/ssh"
// )

// func Match(conn *ssh.Client, logFileName, filter string) (string, error) {
// 	inst := &pkg.Application{}
// 	sftpClient, err := sftp.NewClient(conn)
// 	if err != nil {
// 		fmt.Printf("Failed to create SFTP client connection: %v\n", err)
// 	}
// 	defer sftpClient.Close()

// 	remoteFile, err := sftpClient.Open(logFileName)
// 	if err != nil {
// 		fmt.Printf("Failed to open remote file: %v\n", err)
// 		return "", err
// 	}

// 	defer remoteFile.Close()
// 	outputFilePath, _ := inst.GetClientHomeDir(logFileName)

// 	outFile, err := os.Create(outputFilePath)
// 	//outFile, err := os.Create("filtered_" + logFileName)
// 	//outFile, err := inst.GetClientHomeDir(logFileName)
// 	if err != nil {
// 		return "", err
// 	}
// 	//defer outFile.Close()

// 	scanner := bufio.NewScanner(remoteFile)
// 	pattern := regexp.MustCompile(filter)

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if pattern.MatchString(line) {
// 			if _, err := fmt.Fprintln(outFile, line); err != nil {
// 				return "", err
// 			}
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return "", err
// 	}

// 	return outFile.Name(), err
// }
