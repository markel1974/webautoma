package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/markel1974/webautoma/src/executor"
	"github.com/markel1974/webautoma/src/service"
	"github.com/markel1974/webautoma/src/version"
	"github.com/markel1974/webautoma/src/wd"
	"github.com/markel1974/webautoma/src/wd/base"
	"github.com/markel1974/webautoma/src/wd/caps/chrome"
	"log"
	"net/url"
	"os"
	"strings"
)

//export GOPRIVATE=github.com/markel1974/webautoma
//go mod tidy
//go mod vendor

const (
	defaultResultFile = "log.json"
	defaultImagesFile = "images.json"
)

func createCapture(capture string) []string {
	if len(capture) == 0 {
		return nil
	}
	return strings.Split(capture, ",")
}

func createVariables(variables string) map[string]interface{} {
	if len(variables) == 0 {
		return nil
	}
	var varData map[string]interface{}
	for _, v := range strings.Split(variables, ",") {
		kv := strings.Split(v, "=")
		if len(kv) < 2 {
			continue
		}
		if varData == nil {
			varData = make(map[string]interface{})
		}
		varData[strings.TrimSpace(kv[0])] = kv[1]
	}
	return varData
}

func launch(wdUrl *url.URL, sideFile string, resultFile string, imgFile string, capture string, variablesData map[string]interface{}) error {
	if len(resultFile) == 0 {
		resultFile = defaultResultFile
	}

	if len(imgFile) == 0 {
		imgFile = defaultImagesFile
	}
	captureData := createCapture(capture)
	exec, err := executor.New(sideFile, resultFile, imgFile, captureData, variablesData)
	if err != nil {
		return err
	}
	logType, logLevel := exec.RequiredLogs()
	chromeCaps := chrome.Caps{}
	caps := base.Capabilities{}
	caps.SetLogLevel(logType, logLevel)
	caps.SetChrome(chromeCaps)
	driver := wd.NewWebDriver(caps, wdUrl, false)
	if err := driver.Start(); err != nil {
		return err
	}
	if err := exec.Start(driver); err != nil {
		return err
	}
	return nil
}

func main() {
	var showHelp bool
	var showVersion bool
	var sideFile string
	var urlBase string
	var port int
	var urlPrefix string
	var wdPath string
	var wdOnly bool
	var resultFile string
	var imgFile string
	var capture string
	var variables string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&sideFile, "x", "", "side file")
	flag.StringVar(&urlBase, "b", "http://127.0.0.1", "web driver base url")
	flag.StringVar(&urlPrefix, "u", "", "web driver prefix url")
	flag.IntVar(&port, "p", 9515, "web driver port")
	flag.StringVar(&wdPath, "s", "", "web driver service path")
	flag.BoolVar(&wdOnly, "w", false, "start webdriver only")
	flag.StringVar(&resultFile, "l", defaultResultFile, "result file")
	flag.StringVar(&imgFile, "i", defaultImagesFile, "images file")
	flag.StringVar(&capture, "c", "", "capture event data (comma separated values)")
	flag.StringVar(&variables, "z", "", "side variables (es a=10;b=20), if you start the data with the letter @, the rest should be a filename (in ndjson format)")
	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName, version.AppVersion)
		return
	}

	sUrl := fmt.Sprintf("%s:%d", urlBase, port)
	if len(urlPrefix) > 0 {
		if !strings.HasPrefix(urlPrefix, "/") {
			urlPrefix = "/" + urlPrefix
		}
		sUrl += urlPrefix
	}

	wdUrl, err := url.Parse(sUrl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(wdPath) > 0 {
		svc, err := service.NewChromeService(wdOnly, wdPath, wdUrl, true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err := svc.Start(); err != nil {
			fmt.Println(err.Error())
			return
		}
		if wdOnly {
			fmt.Println("service webdriver successfully started")
			return
		}
	}

	if len(sideFile) == 0 {
		fmt.Println("empty side file")
		flag.Usage()
		return
	}

	if len(variables) == 0 {
		if err := launch(wdUrl, sideFile, resultFile, imgFile, capture, nil); err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	if variables[0] != '@' {
		variablesData := createVariables(variables)
		if err := launch(wdUrl, sideFile, resultFile, imgFile, capture, variablesData); err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	variables = variables[1:]
	file, err := os.OpenFile(variables, os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(file)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*100)
	counter := -1
	for scanner.Scan() {
		counter++
		var variablesData map[string]interface{}
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &variablesData); err != nil {
			log.Printf("line %d: %s", counter, err.Error())
			continue
		}
		if err := launch(wdUrl, sideFile, resultFile, imgFile, capture, variablesData); err != nil {
			log.Printf("line %d: %s", counter, err.Error())
		}
	}
	if scanner.Err() != nil {
		log.Printf("line %d: %s", counter, err.Error())
	}

	/*
		if err := driver.Navigate("https://www.google.com"); err != nil {
			log.Fatal(err.Error())
		}
		body, err := driver.PageSource()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("---- BODY -----")
		fmt.Println(body)

		fmt.Println("---- COOKIES -----")
		cookies, _ := driver.GetCookies()
		for _, cookie := range cookies {
			fmt.Println(cookie)
		}
		_ = driver.Quit()
	*/
}
