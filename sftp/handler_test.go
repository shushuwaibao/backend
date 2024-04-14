package sftp

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestUploadFile(t *testing.T) {
	params := &FileTransferParams{
		User:              "default",
		Password:          "1234",
		ServerWithSshPort: "172.16.13.73:2224",
		TargetPath:        "/home/default/ok.txt",
	}

	file, err := os.Open("./data/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = uploadFile(params, file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadFile(t *testing.T) {
	params := &FileTransferParams{
		User:              "default",
		Password:          "1234",
		ServerWithSshPort: "172.16.13.73:2224",
		TargetPath:        "/home/default/ok.txt",
	}

	file, err := downloadFile(params)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(bytes))
}
