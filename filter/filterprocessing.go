package filter

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Match(conn *ssh.Client, logFileName, filter string) error {
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client connection: %v", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(logFileName)
	if err != nil {
		return fmt.Errorf("failed to open remote file using SFTP: %v", err)
	}
	defer remoteFile.Close()

	pattern, err := regexp.Compile(filter)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	var outputMutex sync.Mutex
	const maxTokenSize = 256 * 1024 // 1 MB
	scanner := bufio.NewScanner(remoteFile)
	scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)

	//scanner := bufio.NewScanner(remoteFile)

	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			wg.Add(1)
			go func(line string) {
				defer wg.Done()

				outputMutex.Lock()
				defer outputMutex.Unlock()

				fmt.Println(line)
			}(line)
		}
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning the remote file: %v", err)
	}

	return nil
}

// package filter

// import (
// 	"bufio"
// 	"fmt"
// 	"regexp"
// 	"sync"
// 	"strings"

// 	"github.com/pkg/sftp"
// 	"golang.org/x/crypto/ssh"
// )

// func Match(conn *ssh.Client, logFileName, filter string) error {
// 	const bufferSize = 1024
// 	sftpClient, err := sftp.NewClient(conn)
// 	if err != nil {
// 		fmt.Printf("Failed to create SFTP client connection: %v\n", err)
// 		return err
// 	}
// 	defer sftpClient.Close()

// 	remoteFile, err := sftpClient.Open(logFileName)
// 	if err != nil {
// 		fmt.Printf("Failed to open remote file using sftp: %v\n", err)
// 		return err
// 	}
// 	defer remoteFile.Close()

// 	scanner := bufio.NewScanner(remoteFile)
// 	pattern := regexp.MustCompile(filter)

// 	var wg sync.WaitGroup
// 	var buffer strings.Builder
// 	var bufferMutex sync.Mutex

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if pattern.MatchString(line) {
// 			wg.Add(1)
// 			go func(line string) {
// 				defer wg.Done()

// 				bufferMutex.Lock()
// 				defer bufferMutex.Unlock()

// 				buffer.WriteString(line + "\n")

// 				if buffer.Len() >= bufferSize {
// 					// Print the buffered lines and reset the buffer
// 					fmt.Print(buffer.String())
// 					buffer.Reset()
// 				}
// 			}(line)
// 		}
// 	}

// 	wg.Wait()

// 	// Print any remaining buffered lines
// 	bufferMutex.Lock()
// 	defer bufferMutex.Unlock()
// 	fmt.Print(buffer.String())

// 	if err := scanner.Err(); err != nil {
// 		fmt.Printf("Error scanning the remote file: %v\n", err)
// 		return err
// 	}

// 	return nil
// }

// package filter

// import (
// 	"bufio"
// 	"fmt"
// 	"regexp"
// 	"sync"

// 	"github.com/pkg/sftp"
// 	"golang.org/x/crypto/ssh"
// )

// func Match(conn *ssh.Client, logFileName, filter string) error {
// 	sftpClient, err := sftp.NewClient(conn)
// 	if err != nil {
// 		fmt.Printf("Failed to create SFTP client connection: %v\n", err)
// 		return err
// 	}
// 	defer sftpClient.Close()

// 	remoteFile, err := sftpClient.Open(logFileName)
// 	if err != nil {
// 		fmt.Printf("Failed to open remote file using sftp: %v\n", err)
// 		return err
// 	}
// 	defer remoteFile.Close()

// 	scanner := bufio.NewScanner(remoteFile)
// 	pattern := regexp.MustCompile(filter)

// 	var wg sync.WaitGroup
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if pattern.MatchString(line) {
// 			wg.Add(1)

// 			go func(line string) {
// 				defer wg.Done()
// 			// Print the matching line to the terminal
// 			fmt.Println(line)
// 			} (line)
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		fmt.Printf("Error scanning the remote file: %v\n", err)
// 		return err
// 	}

// 	return err
// }

// // output filtered logs to a file
// // package filter

// // import (
// // 	"bufio"
// // 	"fmt"
// // 	"os"
// // 	"regexp"

// // 	"github.com/ipsegs/sshpodlog/pkg"
// // 	"github.com/pkg/sftp"
// // 	"golang.org/x/crypto/ssh"
// // )

// // func Match(conn *ssh.Client, logFileName, filter string) (string, error) {
// // 	inst := &pkg.Application{}
// // 	sftpClient, err := sftp.NewClient(conn)
// // 	if err != nil {
// // 		fmt.Printf("Failed to create SFTP client connection: %v\n", err)
// // 	}
// // 	defer sftpClient.Close()

// // 	remoteFile, err := sftpClient.Open(logFileName)
// // 	if err != nil {
// // 		fmt.Printf("Failed to open remote file: %v\n", err)
// // 		return "", err
// // 	}

// // 	defer remoteFile.Close()
// // 	outputFilePath, _ := inst.GetClientHomeDir(logFileName)

// // 	outFile, err := os.Create(outputFilePath)
// // 	//outFile, err := os.Create("filtered_" + logFileName)
// // 	//outFile, err := inst.GetClientHomeDir(logFileName)
// // 	if err != nil {
// // 		return "", err
// // 	}
// // 	//defer outFile.Close()

// // 	scanner := bufio.NewScanner(remoteFile)
// // 	pattern := regexp.MustCompile(filter)

// // 	for scanner.Scan() {
// // 		line := scanner.Text()
// // 		if pattern.MatchString(line) {
// // 			if _, err := fmt.Fprintln(outFile, line); err != nil {
// // 				return "", err
// // 			}
// // 		}
// // 	}

// // 	if err := scanner.Err(); err != nil {
// // 		return "", err
// // 	}

// // 	return outFile.Name(), err
// // }
