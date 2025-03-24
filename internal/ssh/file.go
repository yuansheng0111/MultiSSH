package ssh

import (
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// UploadFile uploads a file to the remote server
func UploadFile(client *ssh.Client, localFilePath string, remoteFilePath string) error {
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
