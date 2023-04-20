package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Start() {
	const varLauncherPath = "CHROME_LAUNCHER_PATH"
	const varLauncherArgs = "CHROME_LAUNCHER_ARGS"
	fmt.Println("using path variable", varLauncherPath)
	fmt.Println("using args variable", varLauncherArgs)
	path := os.Getenv(varLauncherPath)
	if len(path) == 0 {
		path = "." + string(os.PathSeparator) + "chrome.app.exe"
	} else {
		if strings.HasPrefix(path, "\"") {
			path = strings.TrimPrefix(path, "\"")
		}
		if strings.HasSuffix(path, "\"") {
			path = strings.TrimSuffix(path, "\"")
		}
	}
	env := os.Getenv(varLauncherArgs)
	env = strings.Replace(env, "\"", "", -1)
	args := strings.Split(env, ";")
	fmt.Println("current path:", path)
	fmt.Println("current args:", args)
	cmd := exec.Command(path, args...)
	err := cmd.Start()
	if err != nil {
		fmt.Printf(err.Error())
	}
}
