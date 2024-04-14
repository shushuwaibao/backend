package sftp

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type FileTransferReq struct {
	User              string
	Password          string
	ServerWithSshPort string
	WorkType          string
	InputPath         string
	OutputPath        string
}

func HandleFileTransferReq(req *FileTransferReq) error {
	auth := ssh.Password(req.Password)
	config := &ssh.ClientConfig{
		User: req.User,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", req.ServerWithSshPort, config)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	if req.WorkType == "upload" {
		err = uploadFile(sftpClient, req.InputPath, req.OutputPath)
	} else if req.WorkType == "download" {
		err = downloadFile(sftpClient, req.InputPath, req.OutputPath)
	} else {
		err = fmt.Errorf("unknown work type")
	}

	if err != nil {
		return err
	}

	return nil
}

func uploadFile(client *sftp.Client, inputPath string, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	output, err := client.Create(outputPath)
	if err != nil {
		return err
	}

	if _, err = output.ReadFrom(input); err != nil {
		return err
	}

	return nil
}

func downloadFile(client *sftp.Client, inputPath string, outputPath string) error {
	input, err := client.Open(inputPath)
	if err != nil {
		return err
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(output, input); err != nil {
		return err
	}

	return nil
}
