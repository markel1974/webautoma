package wd

import (
	"fmt"
	"github.com/markel1974/webautoma/src/wd/base"
	"github.com/markel1974/webautoma/src/wd/caps/chrome"
	"net/url"
	"os"
	"strings"
	"time"
)

const seq = `package main

import "fmt"

func main() {
	fmt.Println("Hello WebDriver!")
}`

func createDriver() *WebDriver {
	chromeCaps := chrome.Caps{}
	caps := base.Capabilities{}
	caps.SetChrome(chromeCaps)
	wdUrl, _ := url.Parse("http://127.0.0.1:9515")
	driver := NewWebDriver(nil, caps, wdUrl, false)
	if err := driver.Start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return driver
}

func Example() {
	driver := createDriver()

	// Navigate to the simple playground interface.
	if err := driver.Navigate("http://play.golang.org/?simple=1"); err != nil {
		panic(err)
	}

	// Get a reference to the text box containing code.
	elem, err := driver.FindElement(base.ByCSSSelector, "#code")
	if err != nil {
		panic(err)
	}
	// Remove the boilerplate code already in the text box.
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter some new code in text box.
	//err = elem.SendKeys(" ")
	//if err != nil {
	//	panic(err)
	//}

	/*
		for _, v := range []rune(seq) {
			err = elem.SendKeys(string(v))
			time.Sleep(time.Millisecond * 2)
			err = driver.KeyDown(base.DownArrowKey)
		}
	*/
	//offset1 := base.Point{X: 100, Y: 100}

	err = elem.SendKeys(" ")
	var actions []base.KeyAction
	for _, v := range []rune(seq) {
		actions = append(actions, base.KeyDownAction(string(v)))
		actions = append(actions, base.KeyPauseAction(100))
		//actions = append(actions, base.KeyDownAction(base.UpArrowKey))
	}
	driver.StoreKeyActions("keyboard1", actions...)
	if err := driver.PerformActions(); err != nil {
		panic(err)
	}
	if err := driver.ReleaseActions(); err != nil {
		panic(err)
	}

	// Click the run button.
	btn, err := driver.FindElement(base.ByID, "run")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Wait for the program to finish running and get the output.
	outputDiv, err := driver.FindElement(base.ByClassName, "Playground-output")
	if err != nil {
		panic(err)
	}

	var output string
	for {
		output, err = outputDiv.Text()
		if err != nil {
			panic(err)
		}
		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))
	// Example Output:
	// Hello WebDriver!
	//
	// Program exited.
	//os.Exit(0)
	os.Exit(1)
}

func Example1() {
	driver := createDriver()
	// The following shows an example of using the Actions API.
	// Please refer to the WC3 Actions spec for more detailed information.
	if err := driver.Navigate("http://play.golang.org/?simple=1"); err != nil {
		panic(err)
	}

	// Create a point which will be used as an offset to click on the
	// code editor text box element on the page.
	offset := base.Point{X: 100, Y: 200}

	// Call StorePointerActions to store a number of Pointer actions which
	// will be executed sequentially.
	// "mouse1" is used as a unique virtual device identifier for this
	// and future actions.
	// selenium.MousePointer is used to identify the type of the pointer.
	// The stored action chain will move the pointer and click on the code
	// editor text box on the page.

	driver.StorePointerActions("mouse1", base.MousePointer,
		// using base.FromViewport as the move origin
		// which calculates the offset from 0,0.
		// the other valid option is selenium.FromPointer.
		base.PointerMoveAction(0, offset, base.FromViewport),
		base.PointerPauseAction(100),
		base.PointerDownAction(base.LeftButton),
		//base.PointerPauseAction(100),
		//base.PointerUpAction(base.LeftButton),
	)

	// Call StoreKeyActions to store a number of Key actions which
	// will be executed sequentially.
	// "keyboard1" is used as a unique virtual device identifier
	// for this and future actions.
	// The stored action chain will send keyboard inputs to the browser.

	driver.StoreKeyActions("keyboard1",
		//base.KeyDownAction(base.AltKey),
		//base.KeyPauseAction(50),
		//base.KeyDownAction("a"),
		//base.KeyPauseAction(50),
		//base.KeyUpAction("a"),
		//base.KeyUpAction(base.AltKey),
		base.KeyDownAction(base.DownArrowKey),
		base.KeyPauseAction(100),
		base.KeyDownAction(base.DownArrowKey),
		base.KeyPauseAction(100),
		base.KeyDownAction(base.DownArrowKey),
		base.KeyPauseAction(100),
		base.KeyDownAction(base.DownArrowKey),
		base.KeyPauseAction(100),
		base.KeyDownAction(base.DownArrowKey),
	)

	if err := driver.PerformActions(); err != nil {
		panic(err)
	}
	if err := driver.ReleaseActions(); err != nil {
		panic(err)
	}
	os.Exit(0)
}

func Example3() {
	driver := createDriver()
	_ = driver.ResizeWindow("", 640, 480)
	if err := driver.Navigate("http://play.golang.org/?simple=1"); err != nil {
		panic(err)
	}
	startX := 0
	startY := 0
	if active, err := driver.ActiveElement(); err == nil {
		if activeLoc, err := active.Location(); err == nil {
			startX = int(activeLoc.X)
			startY = int(activeLoc.Y)
		}
	}
	elem, err := driver.FindElement(base.ByClassName, "Playground-output")
	if err != nil {
		panic(err)
	}
	loc, err := elem.Location()
	if err != nil {
		return
	}
	driver.StoreWheelActions("wheel1", base.CreateWheelAction(startX, startY, int(loc.X), int(loc.Y)))

	if err := driver.PerformActions(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if err := driver.ReleaseActions(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	os.Exit(1)
}
