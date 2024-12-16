package executor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Downloader struct {
	client *http.Client
}

func NewDownloader(client *http.Client) *Downloader {
	return &Downloader{client: client}
}

func (e *Downloader) Do(link string, cookies []*http.Cookie, fName string) error {
	if len(link) == 0 {
		return fmt.Errorf("invalid link")
	}
	if len(fName) == 0 {
		return fmt.Errorf("invalid value")
	}
	body := bytes.NewBufferString("")
	req, err := http.NewRequest("GET", link, body)
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	res, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("invalid status code %d", res.StatusCode)
	}
	out, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}
	return nil
}

/*
func (e *Downloader) prepareForDownload(target string, value string) (string, int, error) {
	timeout := 60
	if len(value) > 0 {
		t, err := strconv.Atoi(value)
		if err != nil {
			return "", 0, fmt.Errorf("invalid timeout %s", err.Error())
		}
		if t < 0 {
			return "", 0, fmt.Errorf("invalid timeout")
		}
		timeout = t
	}
	if strings.HasPrefix(target, "/") {
		target = e.baseUrl + target
	}
	pos := strings.LastIndex(target, "/")
	if pos < 0 {
		return "", 0, fmt.Errorf("invalid target")
	}
	filename := target[pos+1:]
	fName := e.downloadPath + filename
	_ = os.Remove(fName)
	return fName, timeout, nil
}

func (e *Downloader) waitForDownload(fp string, timeout int) error {
	counter := 0
	for {
		if _, err := os.Stat(fp); err == nil {
			return nil
		}
		time.Sleep(time.Second)
		if counter++; counter >= timeout {
			return fmt.Errorf("timeout")
		}
	}
}

func (e *Downloader) doDownload_old(target string, value string) error {
	if strings.HasPrefix(target, "/") {
		target = e.baseUrl + target
	}
	fName, timeout, err := e.prepareForDownload(target, value)
	if err != nil {
		return err
	}
	if err = e.adapter.Navigate(target); err != nil {
		return err
	}
	if err = e.waitForDownload(fName, timeout); err != nil {
		return err
	}
	return nil
}
*/
