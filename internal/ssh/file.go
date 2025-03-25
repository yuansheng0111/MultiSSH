package ssh

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Upload a file to the remote server
func UploadFile(client *ssh.Client, localFilePath string, remoteFilePath string) error {
	fmt.Printf("Uploading file from %s as %s\n", localFilePath, remoteFilePath)
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	localFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Copy the file contents
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return err
	}

	return nil
}

// Download a file from the remote server
func DownloadFile(client *ssh.Client, remoteFilePath string, localFilePath string) error {
	fmt.Printf("Downloading file from %s as %s\n", remoteFilePath, localFilePath)
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Copy the file contents
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return err
	}

	return nil
}
