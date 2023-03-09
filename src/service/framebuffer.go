package service

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FrameBuffer struct {
	display    string
	authPath   string
	cmd        *exec.Cmd
	screenSize string
}

// "{width}x{height}[x{depth}]", example: "1024x768x24"
func NewFrameBuffer(screenSize string) *FrameBuffer {
	return &FrameBuffer{
		display:    "",
		authPath:   "",
		cmd:        nil,
		screenSize: screenSize,
	}
}

func (fb *FrameBuffer) Display() (string, string) {
	return fb.display, fb.authPath
}

func (fb *FrameBuffer) Start() error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()

	auth, err := ioutil.TempFile("", "webautom-xvfb")
	if err != nil {
		return err
	}
	fb.authPath = auth.Name()
	if err := auth.Close(); err != nil {
		return err
	}

	args := []string{"-displayfd", "3", "-nolisten", "tcp"}
	if len(fb.screenSize) > 0 {
		var screenSizeExpression = regexp.MustCompile(`^\d+x\d+(?:x\d+)$`)
		if !screenSizeExpression.MatchString(fb.screenSize) {
			return fmt.Errorf("invalid screen size: expected 'WxH[xD]', got %q", fb.screenSize)
		}
		args = append(args, "-screen", "0", fb.screenSize)
	}

	fb.cmd = exec.Command("Xvfb", args...)
	fb.cmd.ExtraFiles = []*os.File{w}
	fb.cmd.Env = append(fb.cmd.Env, "XAUTHORITY="+fb.authPath)
	if err := fb.cmd.Start(); err != nil {
		return err
	}
	w.Close()

	type resp struct {
		display string
		err     error
	}

	ch := make(chan resp)
	go func() {
		buf := bufio.NewReader(r)
		s, err := buf.ReadString('\n')
		ch <- resp{s, err}
	}()

	select {
	case resp := <-ch:
		if resp.err != nil {
			return resp.err
		}
		fb.display = strings.TrimSpace(resp.display)
		if _, err := strconv.Atoi(fb.display); err != nil {
			return fmt.Errorf("xvfb did not print the display number")
		}
	case <-time.After(3 * time.Second):
		return fmt.Errorf("timeout waiting for Xvfb")
	}

	xAuthCmd := exec.Command("xauth", "generate", ":"+fb.display, ".", "trusted")
	xAuthCmd.Stderr = os.Stderr
	xAuthCmd.Stdout = os.Stdout
	xAuthCmd.Env = append(xAuthCmd.Env, "XAUTHORITY="+fb.authPath)

	if err := xAuthCmd.Run(); err != nil {
		return err
	}
	return nil
}

func (fb *FrameBuffer) Stop() error {
	if err := fb.cmd.Process.Kill(); err != nil {
		return err
	}
	_ = os.Remove(fb.authPath)
	if err := fb.cmd.Wait(); err != nil && err.Error() != "signal: killed" {
		return err
	}
	return nil
}
