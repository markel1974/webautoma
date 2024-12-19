package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/markel1974/webautoma/src/email"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/markel1974/webautoma/src/executor"
	"github.com/markel1974/webautoma/src/server"
	"github.com/markel1974/webautoma/src/server/asset"
	"github.com/markel1974/webautoma/src/server/config"
	"github.com/markel1974/webautoma/src/server/handlers/httphandler"
	"github.com/markel1974/webautoma/src/service"
	"github.com/markel1974/webautoma/src/version"
	"github.com/markel1974/webautoma/src/wd"
	"github.com/markel1974/webautoma/src/wd/base"
	"github.com/markel1974/webautoma/src/wd/caps/chrome"
)

//export GOPRIVATE=github.com/markel1974/webautoma
//go mod tidy
//go mod vendor

func createServer(listen string) error {
	cfg := &config.Config{
		Listen: listen,
	}
	s := server.New()
	h := httphandler.New()
	if err := s.Setup(h, cfg); err != nil {
		fmt.Println(err.Error())
		return err
	}
	if err := s.Start(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func createAsset(createAsset string) {
	c := strings.Split(createAsset, ":")
	if len(c) < 3 {
		fmt.Println("asset error: ", "wrong format (package:name:dir)")
		return
	}
	packageAsset := c[0]
	packageName := c[1]
	packageDir := c[2]
	fmt.Println("generating asset: ", packageName, packageAsset, packageDir, ".....")
	w := asset.NewWriter()
	if err := w.Setup(packageDir); err != nil {
		fmt.Println("asset error: ", err.Error())
		return
	}
	if err := w.Marshal(packageAsset, "AssetContent", packageName); err != nil {
		fmt.Println("asset error: ", err.Error())
		return
	}
	fmt.Println("asset generated successfully")
	return
}

func createCapture(capture string) []string {
	if len(capture) == 0 {
		return nil
	}
	return strings.Split(capture, ",")
}

func createVariables(variables string) (map[string]interface{}, error) {
	if len(variables) == 0 {
		return nil, nil
	}
	var varData map[string]interface{}
	if err := json.Unmarshal([]byte(variables), &varData); err != nil {
		return nil, err
	}
	return varData, nil
}

func launch(console bool, wdUrl *url.URL, args []string, sideFile string, resultFile string, imgFile string, imgDump bool, capture string, variables string) error {
	chromeCaps := chrome.Caps{}
	for _, arg := range args {
		chromeCaps.Args = append(chromeCaps.Args, arg /*"headless"*/)
	}
	caps := base.Capabilities{}
	logType, logLevel := executor.RequiredLogs()
	caps.SetLogLevel(logType, logLevel)
	caps.SetChrome(chromeCaps)
	driver := wd.NewWebDriver(nil, caps, wdUrl, false)
	if err := driver.Start(); err != nil {
		return err
	}
	variablesData, err := createVariables(variables)
	if err != nil {
		return err
	}
	captureData := createCapture(capture)
	exec := executor.New(driver)

	if err = exec.Setup(sideFile, resultFile, imgFile, imgDump, captureData, variablesData); err != nil {
		return err
	}
	if err = exec.Start(console); err != nil {
		return err
	}
	return nil
}

/*
func test() {
	download := executor.NewDownloader(&http.Client{})
	html := `<html>
<a href="https://www.w3schools.com">Visit W3Schools</a>
<a href="https://www.w3schools.com">Visit W3Schools</a>
<a href="https://www.w3schools.com">Visit W3Schools</a>
	</html>`

	download.DoRetrieveHref(nil, html, "")
	os.Exit(0)
}
*/

func emailTester() {
	//Email: dimon.probes@telecomitalia.it
	//const userId = "CD038910"
	//const pwd = "Ggfrt56_87&"
	//const host = "mail.telecomitalia.it:993" // "10.14.252.100:993" //
	//mode := email.ModeTLS
	target := "startTLS://markel@tin.it:Cristiana1976@box.tin.it:143"
	value := ".+One-Time Password" + "|||" + "([0-9]+) is your One-Time Password to login|||60|||1440"
	options, err := email.NewOptions(value)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(0)
	}
	z, err := email.NewClientFromTarget(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	verifyInterval := time.Now().Unix() + options.VerifyInterval()
	var k string
	for {
		if k, err = z.Retrieve(options); err == nil {
			break
		}
		time.Sleep(time.Second * 5)
		if time.Now().Unix() > verifyInterval {
			break
		}
	}
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(0)
	}
	fmt.Println(k)

	os.Exit(0)
}

func main() {
	//emailTester()
	//launcher.Start()
	//os.Exit(-1)
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
	var driverArgs string
	var buildAsset string
	var imgDump bool
	var listen string
	var sam bool

	//wd.Example3()

	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&driverArgs, "d", "", "webdriver args (semicolon separated)")
	flag.StringVar(&sideFile, "x", "", "side file")
	flag.StringVar(&urlBase, "b", "http://127.0.0.1", "webdriver base url")
	flag.StringVar(&urlPrefix, "u", "", "webdriver prefix url")
	flag.IntVar(&port, "p", 9515, "webdriver port")
	flag.StringVar(&wdPath, "s", "", "webdriver service path")
	flag.BoolVar(&wdOnly, "w", false, "start webdriver only")
	flag.StringVar(&resultFile, "l", executor.DefaultResultFile, "result file")
	flag.StringVar(&imgFile, "i", executor.DefaultImagesFile, "images file")
	flag.StringVar(&capture, "c", "", "capture event data (comma separated values)")
	flag.StringVar(&buildAsset, "t", "", "create file asset (format package:name:dir)")
	flag.StringVar(&listen, "y", "", "server listen")
	flag.StringVar(&variables, "z", "", "side variables (es a=10;b=20), if you start the data with the letter @, the rest should be a filename (in ndjson format)")
	flag.BoolVar(&imgDump, "a", false, "image dump")
	flag.BoolVar(&sam, "e", false, "launch SAM console")

	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName, version.AppVersion)
		return
	}

	if len(buildAsset) > 0 {
		createAsset(buildAsset)
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

	if len(listen) > 0 {
		if err := createServer(listen); err != nil {
			fmt.Println(err.Error())
			return
		}
		return
	}

	if len(sideFile) == 0 {
		fmt.Println("empty side file")
		flag.Usage()
		return
	}

	var args []string
	if len(driverArgs) > 0 {
		for _, arg := range strings.Split(driverArgs, ";") {
			args = append(args, arg)
		}
	}

	if len(variables) > 0 && variables[0] == '@' {
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
			if err := launch(sam, wdUrl, args, sideFile, resultFile, imgFile, imgDump, capture, scanner.Text()); err != nil {
				log.Printf("line %d: %s", counter, err.Error())
			}
		}
		if scanner.Err() != nil {
			log.Printf("line %d: %s", counter, err.Error())
		}
	} else {
		if err = launch(sam, wdUrl, args, sideFile, resultFile, imgFile, imgDump, capture, variables); err != nil {
			log.Fatal(err.Error())
		}
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
