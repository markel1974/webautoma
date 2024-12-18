package executor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

// Downloader is a structure that facilitates file downloads using an HTTP client.
// It manages HTTP requests and handles response processing for file saving.
type Downloader struct {
	client *http.Client
}

// NewDownloader creates and returns a new Downloader instance configured with the provided HTTP client.
func NewDownloader(client *http.Client) *Downloader {
	return &Downloader{client: client}
}

// Do perform an HTTP GET request to the specified link, uses provided cookies, and writes the response to a file.
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

const _r = `<\s*a\s+href\s*=\s*"([^"]+)"`

var _rgx = regexp.MustCompile(_r)

func (e *Downloader) DoRetrieveHref(cookies []*http.Cookie, html string, pattern string) error {
	values := _rgx.FindAllStringSubmatch(html, -1)
	for _, v := range values {
		if len(v) > 1 {
			fmt.Println(v[1])
		}
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
