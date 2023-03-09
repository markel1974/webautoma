package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/webautoma/src/executor"
	"github.com/markel1974/webautoma/src/version"
	"github.com/markel1974/webautoma/src/wd"
	"github.com/markel1974/webautoma/src/wd/base"
	"github.com/markel1974/webautoma/src/wd/caps/chrome"
	"log"
)

//export GOPRIVATE=github.com/markel1974/webautoma
//go mod tidy
//go mod vendor

func main() {
	var showHelp bool
	var showVersion bool
	var execFile string
	var url string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&execFile, "x", "", "exec file")
	flag.StringVar(&url, "u", "http://127.0.0.1:9515", "url")
	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName, version.AppVersion)
		return
	}

	if len(execFile) == 0 {
		log.Fatal("empty exec file")
		return
	}

	exec, err := executor.New(execFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	logType, logLevel := exec.RequiredLogs()

	chromeCaps := chrome.Caps{}
	caps := base.Capabilities{}
	caps.SetLogLevel(logType, logLevel)
	caps.SetChrome(chromeCaps)
	driver := wd.NewWebDriver(caps, url)

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
