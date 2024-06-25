package executor

import (
	"fmt"
	"github.com/markel1974/webautoma/src/shell"
	"github.com/markel1974/webautoma/src/shell/cli"
)

func commandHandler() *cli.Command {
	root := cli.NewCommand()
	root.Run = func(cmd *cli.Command, pid int, args []string) {}
	root.Use = "root"
	root.Long = "Root Command"
	root.Short = root.Long

	next := cli.NewCommand()
	next.Activate = false
	next.Use = "next"
	next.Long = "Move cursor to next command"
	next.Short = next.Long
	next.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\nBar is opened, pid %d", pid))
		r.Deactivate(pid)
	}
	next.ReadEvent = func(cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
	}

	prev := cli.NewCommand()
	prev.Use = "prev"
	prev.Long = "Move cursor to previous command"
	prev.Short = prev.Long
	prev.Activate = false
	prev.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "prev 1"))
		r.Deactivate(pid)
	}

	step := cli.NewCommand()
	step.Use = "step"
	step.Long = "Run current command only"
	step.Short = step.Long
	step.Activate = false
	step.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "stop 2"))
		r.Deactivate(pid)
	}

	redo := cli.NewCommand()
	redo.Use = "redo"
	redo.Long = "Repeat current command only"
	redo.Short = redo.Long
	redo.Activate = false
	redo.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "redo 3"))
		r.Deactivate(pid)
	}

	stop := cli.NewCommand()
	stop.Use = "stop"
	stop.Long = "Stop execution"
	stop.Short = stop.Long
	stop.Activate = false
	stop.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "stop 4"))
		r.Deactivate(pid)
	}

	run := cli.NewCommand()
	run.Use = "run"
	run.Long = "Start execution"
	run.Short = run.Long
	run.Activate = false
	run.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "run 5"))
		r.Deactivate(pid)
	}

	edit := cli.NewCommand()
	edit.Use = "edit"
	edit.Long = "Edit current command"
	edit.Short = edit.Long
	edit.Activate = false
	edit.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "edit 6"))
		r.Deactivate(pid)
	}

	goTo := cli.NewCommand()
	goTo.Use = "goto"
	goTo.Long = "Set current command"
	goTo.Short = goTo.Long
	goTo.Activate = false
	goTo.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "goto 7"))
		r.Deactivate(pid)
	}

	list := cli.NewCommand()
	list.Use = "list"
	list.Long = "list entrie list"
	list.Short = list.Long
	list.Activate = false
	list.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "Test"))
		r.Deactivate(pid)
	}

	cursor := cli.NewCommand()
	cursor.Use = "cursor"
	cursor.Long = "Current cursor position"
	cursor.Short = cursor.Long
	cursor.Activate = true
	cursor.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "Test"))
		r.Deactivate(pid)
	}

	printS := cli.NewCommand()
	printS.Use = "print"
	printS.Long = "Print available data"
	printS.Short = printS.Long
	printS.Activate = true
	printS.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "Test"))
		r.Deactivate(pid)
	}

	printHTML := cli.NewCommand()
	printHTML.Use = "html"
	printHTML.Long = "Print current html page"
	printHTML.Short = printHTML.Long
	printHTML.Activate = true
	printHTML.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "Test"))
		r.Deactivate(pid)
	}

	printNetwork := cli.NewCommand()
	printNetwork.Use = "net"
	printNetwork.Long = "Print current network state"
	printNetwork.Short = printNetwork.Long
	printNetwork.Activate = true
	printNetwork.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn(fmt.Sprintf("\r\n" + "Test"))
		r.Deactivate(pid)
	}

	_ = printS.AddCommand(printHTML)
	_ = printS.AddCommand(printNetwork)

	_ = root.AddCommand(next)
	_ = root.AddCommand(prev)
	_ = root.AddCommand(step)
	_ = root.AddCommand(redo)
	_ = root.AddCommand(stop)
	_ = root.AddCommand(run)
	_ = root.AddCommand(goTo)
	_ = root.AddCommand(list)
	_ = root.AddCommand(cursor)
	_ = root.AddCommand(printS)
	return root
}

func CreateConsole() error {
	if err := shell.Create(false, commandHandler()); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
