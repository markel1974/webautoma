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
	"strings"
)

//export GOPRIVATE=github.com/markel1974/webautoma
//go mod tidy
//go mod vendor

func main() {
	var showHelp bool
	var showVersion bool
	var execFile string
	var baseUrl string
	var port int
	var webDriverPath string
	var startWebDriverOnly bool
	var logFile string
	var imgFile string
	var capture string
	var variables string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&execFile, "x", "", "exec file")
	flag.StringVar(&baseUrl, "b", "http://127.0.0.1", "web driver base url")
	flag.IntVar(&port, "p", 9515, "web driver port")
	flag.StringVar(&webDriverPath, "s", "", "web driver service path")
	flag.BoolVar(&startWebDriverOnly, "w", false, "start webdriver only")
	flag.StringVar(&logFile, "l", "log.json", "default log file")
	flag.StringVar(&imgFile, "i", "images.json", "default image file")
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

	if len(webDriverPath) > 0 {
		svc, err := service.NewChromeService(webDriverPath, port, baseUrl, true)
		if err != nil {
			log.Fatal(err.Error())
		}
		if err := svc.Start(); err != nil {
			log.Fatal(err.Error())
		}
		if startWebDriverOnly {
			return
		}
	}

	if len(execFile) == 0 {
		log.Fatal("empty exec file")
		return
	}

	if len(logFile) == 0 {
		logFile = "log.json"
	}

	if len(imgFile) == 0 {
		imgFile = "images.json"
	}

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

	var captureData []string = nil
	if len(capture) > 0 {
		captureData = strings.Split(capture, ",")
	}
	exec, err := executor.New(execFile, logFile, imgFile, captureData, varData)
	if err != nil {
		log.Fatal(err.Error())
	}
	logType, logLevel := exec.RequiredLogs()

	wdUrl := fmt.Sprintf("%s:%d", baseUrl, port)

	chromeCaps := chrome.Caps{}
	caps := base.Capabilities{}
	caps.SetLogLevel(logType, logLevel)
	caps.SetChrome(chromeCaps)
	driver := wd.NewWebDriver(caps, wdUrl)

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
