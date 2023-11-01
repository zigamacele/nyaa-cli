package torrent

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

func OpenFile(filename string) error {
	cmd := exec.Command("open", filename) // windows:  use "start" instead of "open"
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func DownloadFile(filepath string, url string) (err error) {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	OpenFile(filepath)
	return nil
}

func DownloadTorrent(download_link string, filename string) {
	var URL string = "https://nyaa.si" + download_link
	var FileName string = filename + ".torrent"
	DownloadFile(FileName, URL)
	OpenFile(FileName)
}
