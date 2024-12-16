package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/markel1974/webautoma/src/shell"
	"github.com/markel1974/webautoma/src/shell/cli"
)

const (
	MessageNext      = 0
	MessagePrev      = 1
	MessageStep      = 2
	MessageList      = 3
	MessageRedo      = 4
	MessageStop      = 5
	MessageRun       = 6
	MessageEdit      = 7
	MessageJump      = 8
	MessageCurr      = 9
	MessagePrintHtml = 11
	MessagePrintNet  = 12
	MessageWindows   = 13
	MessageQuit      = 255
)

type ISamMessage interface {
	GetType() int
	SetResponse(id string, x int, result string, err error)
	GetResponse() []string
}

type SamMessage struct {
	kind     int
	response chan []string
}

func NewSamMessage(kind int) *SamMessage {
	return &SamMessage{
		kind:     kind,
		response: make(chan []string),
	}
}
func (s *SamMessage) GetType() int {
	return s.kind
}
func (s *SamMessage) SetResponse(id string, x int, result string, err error) {
	var out []string
	var status string
	if err != nil {
		v := strings.Replace(err.Error(), "\r", "", -1)
		out = append(out, strings.Split(v, "\n")...)
		status = "err"
	} else {
		if len(result) > 0 {
			p := strings.Replace(result, "\r", "", -1)
			out = append(out, strings.Split(p, "\n")...)
		}
		status = "ok"
	}
	out = append(out, fmt.Sprintf("%s *[%d] -> %s", status, x, id))
	s.response <- out
}
func (s *SamMessage) GetResponse() []string {
	v := <-s.response
	return v
}

type SamMessageEdit struct {
	*SamMessage
	key string
	val string
}

func NewSamMessageEdit(key string, val string) *SamMessageEdit {
	return &SamMessageEdit{
		SamMessage: NewSamMessage(MessageEdit),
		key:        key,
		val:        val,
	}
}
func (sme *SamMessageEdit) KeyVal() (string, string) {
	return sme.key, sme.val
}

type SamMessageJump struct {
	*SamMessage
	jump int
}

func NewSamMessageJump(jump int) *SamMessageJump {
	return &SamMessageJump{
		SamMessage: NewSamMessage(MessageJump),
		jump:       jump,
	}
}
func (smj *SamMessageJump) Jump() int {
	return smj.jump
}

type SamMessageStep struct {
	*SamMessage
	advance bool
	count   int
}

func NewSamMessageStep(advance bool, count int) *SamMessageStep {
	return &SamMessageStep{
		SamMessage: NewSamMessage(MessageStep),
		advance:    advance,
		count:      count,
	}
}
func (smj *SamMessageStep) Advance() bool {
	return smj.advance
}

type SamMessageHTML struct {
	*SamMessage
	fileId string
}

func NewSamMessageHTML(fileId string) *SamMessageHTML {
	return &SamMessageHTML{
		SamMessage: NewSamMessage(MessagePrintHtml),
		fileId:     fileId,
	}
}
func (smj *SamMessageHTML) FileId() string {
	return smj.fileId
}

type SamMessageNet struct {
	*SamMessage
	fileId string
}

func NewSamMessageNet(fileId string) *SamMessageNet {
	return &SamMessageNet{
		SamMessage: NewSamMessage(MessagePrintNet),
		fileId:     fileId,
	}
}
func (smj *SamMessageNet) FileId() string {
	return smj.fileId
}

type SamMessages chan ISamMessage

type SamConsole struct {
	e        *Executor
	messages SamMessages
}

func NewSam(e *Executor) *SamConsole {
	return &SamConsole{
		e:        e,
		messages: make(SamMessages, 64),
	}
}

func (sm *SamConsole) Start() error {
	quit := make(chan bool)
	if err := shell.Create(false, sm.commandHandler(), quit); err != nil {
		return err
	}
	sm.eventLoop(quit)
	return nil
}

func (sm *SamConsole) Print(s string) {
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\n", "\r\n", -1)
	fmt.Printf("\r\n%s", s)
}

func (sm *SamConsole) runCommand(cmd *cli.Command, pid int, s ISamMessage) {
	fmt.Printf("\r\n%s", "waiting....")
	r := cmd.GetRootContext()
	sm.e.SamSendMessage(s)
	resp := s.GetResponse()
	fmt.Printf("\r%s", "           ")
	fmt.Printf("\r%s", strings.Join(resp, "\r\n"))
	r.Deactivate(pid)
}

func (sm *SamConsole) SendMessage(m ISamMessage) {
	sm.messages <- m
}

func (sm *SamConsole) eventLoop(quit chan bool) {
	test := sm.e.cfg.Tests[0]
	pc := 0

	for {
		select {
		case q := <-quit:
			if q {
				return
			}
		case msg := <-sm.messages:
			switch msg.GetType() {
			case MessageNext:
				pc++
				if pc >= len(test.Commands) {
					pc = 0
				}
				msg.SetResponse(test.Commands[pc].Id, pc, "", nil)
			case MessagePrev:
				pc--
				if pc < 0 {
					pc = 0
				}
				msg.SetResponse(test.Commands[pc].Id, pc, "", nil)
			case MessageStep:
				sms := msg.(*SamMessageStep)
				var pre string
				var err error
				for x := 0; x < sms.count; x++ {
					var jump int
					_, jump, err = sm.e.doCommand(test.Commands, pc)
					if err != nil {
						pre = ""
						break
					}
					if sms.Advance() {
						if jump >= 0 {
							pc = jump
							pre = "jump found"
						} else {
							pc++
						}
						if pc < 0 || pc >= len(test.Commands) {
							pc = 0
						}
					}
				}
				msg.SetResponse(test.Commands[pc].Id, pc, pre, err)
			case MessageList:
				m, _ := json.MarshalIndent(test.Commands, "", "  ")
				msg.SetResponse(test.Commands[pc].Id, pc, string(m), nil)
			case MessageCurr:
				m, _ := json.MarshalIndent(test.Commands[pc], "", "  ")
				msg.SetResponse(test.Commands[pc].Id, pc, string(m), nil)
			case MessageJump:
				msgJump := msg.(*SamMessageJump)
				pc = msgJump.Jump()
				if pc < 0 || pc >= len(test.Commands) {
					pc = 0
				}
				msg.SetResponse(test.Commands[pc].Id, pc, "", nil)
			case MessageEdit:
				msgEdit := msg.(*SamMessageEdit)
				k, v := msgEdit.KeyVal()
				err := AlterObject(k, v, &test.Commands[pc])
				msg.SetResponse(test.Commands[pc].Id, pc, "", err)
			case MessagePrintHtml:
				smh := msg.(*SamMessageHTML)
				h, err := sm.e.doRetrievePageSource()
				if len(smh.FileId()) > 0 {
					_ = os.WriteFile(smh.FileId(), []byte(h), 0644)
				}
				msg.SetResponse(test.Commands[pc].Id, pc, h, err)
			case MessagePrintNet:
				smh := msg.(*SamMessageNet)
				n, err := sm.e.doRetrieveNetworkHeaders()
				if len(smh.FileId()) > 0 {
					_ = os.WriteFile(smh.FileId(), []byte(n), 0644)
				}
				msg.SetResponse(test.Commands[pc].Id, pc, n, err)
			case MessageWindows:
				v, _ := sm.e.doListWindows()
				m, _ := json.MarshalIndent(v, "", "  ")
				msg.SetResponse(test.Commands[pc].Id, pc, string(m), nil)
			case MessageQuit:
				return
			}
		}
	}
}

func (sm *SamConsole) commandHandler() *cli.Command {
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
		msg := NewSamMessage(MessageNext)
		sm.runCommand(cmd, pid, msg)
	}
	next.ReadEvent = func(cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
	}

	prev := cli.NewCommand()
	prev.Use = "prev"
	prev.Long = "Move cursor to previous command"
	prev.Short = prev.Long
	prev.Activate = false
	prev.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessagePrev)
		sm.runCommand(cmd, pid, msg)
	}

	step := cli.NewCommand()
	step.Use = "step"
	step.Long = "Run current command only"
	step.Short = step.Long
	step.Activate = false
	step.Run = func(cmd *cli.Command, pid int, args []string) {
		count := 1
		if len(args) > 0 {
			if c, err := strconv.Atoi(args[0]); err == nil && c > 0 {
				count = c
			}
		}
		msg := NewSamMessageStep(true, count)
		sm.runCommand(cmd, pid, msg)
	}

	redo := cli.NewCommand()
	redo.Use = "redo"
	redo.Long = "Repeat current command only"
	redo.Short = redo.Long
	redo.Activate = false
	redo.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessageStep(false, 1)
		sm.runCommand(cmd, pid, msg)
	}

	stop := cli.NewCommand()
	stop.Use = "stop"
	stop.Long = "Stop execution"
	stop.Short = stop.Long
	stop.Activate = false
	stop.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessageStop)
		sm.runCommand(cmd, pid, msg)
	}

	run := cli.NewCommand()
	run.Use = "run"
	run.Long = "Start execution"
	run.Short = run.Long
	run.Activate = false
	run.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessageRun)
		sm.runCommand(cmd, pid, msg)
	}

	edit := cli.NewCommand()
	edit.Use = "edit"
	edit.Long = "Edit current command"
	edit.Short = edit.Long
	edit.Activate = false
	edit.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		if len(args) < 2 {
			fmt.Printf("\r\nmissing arguments")
			r.Deactivate(pid)
			return
		}
		key := args[0]
		val := args[1]
		msg := NewSamMessageEdit(key, val)
		sm.runCommand(cmd, pid, msg)
	}

	jump := cli.NewCommand()
	jump.Use = "jump"
	jump.Long = "Set current command"
	jump.Short = jump.Long
	jump.Activate = false
	jump.Run = func(cmd *cli.Command, pid int, args []string) {
		j := -1
		if len(args) > 0 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				j = v
			}
		}
		msg := NewSamMessageJump(j)
		sm.runCommand(cmd, pid, msg)
	}

	list := cli.NewCommand()
	list.Use = "list"
	list.Long = "list entrie list"
	list.Short = list.Long
	list.Activate = false
	list.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessageList)
		sm.runCommand(cmd, pid, msg)
	}

	curr := cli.NewCommand()
	curr.Use = "curr"
	curr.Long = "Current cursor position"
	curr.Short = curr.Long
	curr.Activate = true
	curr.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessageCurr)
		sm.runCommand(cmd, pid, msg)
	}

	printS := cli.NewCommand()
	printS.Use = "print"
	printS.Long = "Print available data"
	printS.Short = printS.Long
	printS.Activate = true
	printS.Run = func(cmd *cli.Command, pid int, args []string) {
		fmt.Printf("\r\nmissing arguments")
		cmd.GetRootContext().Deactivate(pid)
	}

	printHTML := cli.NewCommand()
	printHTML.Use = "html"
	printHTML.Long = "Print current html page [fileId]"
	printHTML.Short = printHTML.Long
	printHTML.Activate = true
	printHTML.Run = func(cmd *cli.Command, pid int, args []string) {
		fileId := ""
		if len(args) > 0 {
			fileId = args[0]
		}
		msg := NewSamMessageHTML(fileId)
		sm.runCommand(cmd, pid, msg)
	}

	printNetwork := cli.NewCommand()
	printNetwork.Use = "net"
	printNetwork.Long = "Print current network state"
	printNetwork.Short = printNetwork.Long
	printNetwork.Activate = true
	printNetwork.Run = func(cmd *cli.Command, pid int, args []string) {
		fileId := ""
		if len(args) > 0 {
			fileId = args[0]
		}
		msg := NewSamMessageNet(fileId)
		sm.runCommand(cmd, pid, msg)
	}

	windows := cli.NewCommand()
	windows.Use = "windows"
	windows.Long = "Print available windows"
	windows.Short = windows.Long
	windows.Activate = true
	windows.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessageWindows)
		sm.runCommand(cmd, pid, msg)
	}

	quit := cli.NewCommand()
	quit.Use = "exit"
	quit.Long = "exit"
	quit.Short = quit.Long
	quit.Activate = true
	quit.Run = func(cmd *cli.Command, pid int, args []string) {
		msg := NewSamMessage(MessageQuit)
		sm.runCommand(cmd, pid, msg)
	}

	_ = printS.AddCommand(printHTML)
	_ = printS.AddCommand(printNetwork)

	_ = root.AddCommand(next)
	_ = root.AddCommand(prev)
	_ = root.AddCommand(step)
	_ = root.AddCommand(redo)
	_ = root.AddCommand(stop)
	_ = root.AddCommand(run)
	_ = root.AddCommand(jump)
	_ = root.AddCommand(edit)
	_ = root.AddCommand(list)
	_ = root.AddCommand(curr)
	_ = root.AddCommand(printS)
	_ = root.AddCommand(windows)
	_ = root.AddCommand(quit)
	return root
}

func AlterObject(k string, v string, cmd interface{}) error {
	rv := reflect.ValueOf(cmd)
	trv := reflect.TypeOf(cmd)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("object isn't a pointer")
	}
	rv = reflect.ValueOf(cmd).Elem()
	trv = reflect.TypeOf(cmd).Elem()
	if rv.IsZero() {
		return fmt.Errorf("invalid object")
	}
	//fmt.Println(rv.CanSet())
	//value := reflect.New(rv.Type())
	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		fieldType := trv.Field(i)
		name := fieldType.Name
		tagName, _ := fieldType.Tag.Lookup("json")
		if tagName != k && name != k {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(v)
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint64, reflect.Uint32:
			d, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(d)
			return nil
		case reflect.Float64, reflect.Float32:
			f, err := strconv.ParseFloat(v, 64)
			if err == nil {
				return err
			}
			fieldValue.SetFloat(f)
			return nil
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			fieldValue.SetBool(b)
			return nil
		}
	}
	return fmt.Errorf("unupported")
}
