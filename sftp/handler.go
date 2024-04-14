package sftp

import (
	"io"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type FileTransferParams struct {
	User              string `json:"user"`
	Password          string `json:"password"`
	ServerWithSshPort string `json:"server_with_ssh_port"`
	TargetPath        string `json:"target_path"`
}

func getClients(params *FileTransferParams) (*ssh.Client, *sftp.Client, error) {
	auth := ssh.Password(params.Password)
	config := &ssh.ClientConfig{
		User: params.User,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", params.ServerWithSshPort, config)
	if err != nil {
		return nil, nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, nil, err
	}

	return sshClient, sftpClient, nil
}

func uploadFile(params *FileTransferParams, file io.Reader) error {
	sshClient, sftpClient, err := getClients(params)
	if err != nil {
		return err
	}

	defer func() {
		sftpClient.Close()
		sshClient.Close()
	}()

	target, err := sftpClient.Create(params.TargetPath)
	if err != nil {
		return err
	}

	if _, err = target.ReadFrom(file); err != nil {
		return err
	}

	return nil
}

func downloadFile(params *FileTransferParams) (io.Reader, error) {
	sshClient, sftpClient, err := getClients(params)
	if err != nil {
		return nil, err
	}

	defer func() {
		sftpClient.Close()
		sshClient.Close()
	}()

	target, err := sftpClient.Open(params.TargetPath)
	if err != nil {
		return nil, err
	}

	return target, nil
}
