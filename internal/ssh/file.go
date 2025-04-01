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
	// Get remote address information for debugging
	// remoteAddr := client.RemoteAddr().String()
	// localAddr := client.LocalAddr().String()
	// fmt.Printf("[DEBUG] UploadFile called for connection from %s to %s\n", localAddr, remoteAddr)
	// fmt.Printf("[DEBUG] Uploading file from %s to %s as %s\n", localFilePath, remoteAddr, remoteFilePath)

	// Create a new SFTP client for this upload
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	localFile, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer remoteFile.Close()

	// Copy the file contents from local to remote
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

// Download a file from the remote server
func DownloadFile(client *ssh.Client, remoteFilePath string, localFilePath string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// Copy the file contents from remote to local
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
