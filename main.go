package main

import (
	"github.com/markel1974/webautoma/src/executor"
	"github.com/markel1974/webautoma/src/wd"
	"github.com/markel1974/webautoma/src/wd/base"
	"github.com/markel1974/webautoma/src/wd/caps/chrome"
	"log"
)

//export GOPRIVATE=github.com/markel1974/webautoma
//go mod tidy
//go mod vendor

func main() {
	chromeCaps := chrome.Caps{}
	caps := base.Capabilities{}
	exec := executor.New()

	logType, logLevel := exec.RequiredLogs()
	caps.SetLogLevel(logType, logLevel)
	caps.AddChrome(chromeCaps)
	driver, err := wd.NewWebDriver(caps, "http://127.0.0.1:9515")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := exec.Setup(driver, "test"); err != nil {
		log.Fatal(err.Error())
	}

	if err := exec.Run(); err != nil {
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
