package main

import (
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
	"strings"
)

//export GOPRIVATE=github.com/markel1974/webautoma
//go mod tidy
//go mod vendor

const (
	defaultLogFile    = "log.json"
	defaultImagesFile = "images.json"
)

func createCapture(capture string) []string {
	var captureData []string
	if len(capture) > 0 {
		captureData = strings.Split(capture, ",")
	}
	return captureData
}

func createVariables(variables string) map[string]interface{} {
	var varData map[string]interface{}
	if len(variables) > 0 {
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
	}
	return varData
}

func main() {
	var showHelp bool
	var showVersion bool
	var execFile string
	var urlBase string
	var port int
	var urlPrefix string
	var wdPath string
	var wdOnly bool
	var logFile string
	var imgFile string
	var capture string
	var variables string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&execFile, "x", "", "exec file")
	flag.StringVar(&urlBase, "b", "http://127.0.0.1", "web driver base url")
	flag.StringVar(&urlPrefix, "u", "", "web driver prefix url")
	flag.IntVar(&port, "p", 9515, "web driver port")
	flag.StringVar(&wdPath, "s", "", "web driver service path")
	flag.BoolVar(&wdOnly, "w", false, "start webdriver only")
	flag.StringVar(&logFile, "l", defaultLogFile, "log file")
	flag.StringVar(&imgFile, "i", defaultImagesFile, "images file")
	flag.StringVar(&capture, "c", "", "capture event data (comma separated values)")
	flag.StringVar(&variables, "z", "", "variables (es a=10;b=20")
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
		log.Fatal(err.Error())
	}

	if len(wdPath) > 0 {
		svc, err := service.NewChromeService(wdOnly, wdPath, wdUrl, true)
		if err != nil {
			log.Fatal(err.Error())
		}
		if err := svc.Start(); err != nil {
			log.Fatal(err.Error())
		}
		if wdOnly {
			fmt.Println("service webdriver successfully started")
			return
		}
	}

	if len(execFile) == 0 {
		fmt.Println("empty exec file")
		flag.Usage()
		return
	}

	if len(logFile) == 0 {
		logFile = defaultLogFile
	}

	if len(imgFile) == 0 {
		imgFile = defaultImagesFile
	}

	exec, err := executor.New(execFile, logFile, imgFile, createCapture(capture), createVariables(variables))
	if err != nil {
		log.Fatal(err.Error())
	}
	logType, logLevel := exec.RequiredLogs()

	chromeCaps := chrome.Caps{}
	caps := base.Capabilities{}
	caps.SetLogLevel(logType, logLevel)
	caps.SetChrome(chromeCaps)
	driver := wd.NewWebDriver(caps, wdUrl, false)

	if err := driver.Start(); err != nil {
		log.Fatal(err.Error())
	}

	if err := exec.Start(driver); err != nil {
		log.Fatal(err.Error())
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
