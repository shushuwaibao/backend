package sftp

import "testing"

func TestUploadFile(t *testing.T) {
	req := &FileTransferReq{
		User: "default",
		Password: "1234",
		ServerWithSshPort: "172.16.13.73:2224",
		WorkType: "upload",
		InputPath: "./test.txt",
		OutputPath: "/home/default/ok.txt",
	}
	err := HandleFileTransferReq(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadFile(t *testing.T) {
	req := &FileTransferReq{
		User: "default",
		Password: "1234",
		ServerWithSshPort: "172.16.13.73:2224",
		WorkType: "download",
		InputPath: "/home/default/ok.txt",
		OutputPath: "./backcat.txt",
	}
	err := HandleFileTransferReq(req)
	if err != nil {
		t.Fatal(err)
	}
}