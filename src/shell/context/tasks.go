/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/markel1974/webautoma/src/shell/adaptiveticker"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type taskState int

const (
	taskStateSetup   taskState = iota
	taskStateRunning taskState = iota
)

const tasksFileExtension = ".task"

const (
	commandActivate = "activate"
	commandTask     = "task"
)

type Task struct {
	cmd     *cli.Command
	context interface{}
	timers  []int
	pid     int
	state   taskState
	caption string
	Line    string
	OffsetX int
	OffsetY int
	Scale   float64
}

func NewTask(cmd *cli.Command, line string) *Task {
	return &Task{
		cmd:     cmd,
		context: nil,
		state:   taskStateSetup,
		caption: "",
		Line:    line,
		OffsetX: 0,
		OffsetY: 0,
		Scale:   1.0,
	}
}

func (t *Task) SetId(id int) {
	t.pid = id
}

func (t *Task) Unset() {
	t.pid = adaptiveticker.UnknownId
}

func (t *Task) Paint(surface *Surface) {
	if t.cmd.PaintEvent == nil {
		return
	}
	caption := strconv.Itoa(t.pid)
	if len(t.caption) > 0 {
		caption += " - " + t.caption
	}
	surface.SetOffsetX(t.OffsetX)
	surface.SetOffsetY(t.OffsetY)
	surface.SetScale(t.Scale)
	surface.SetCaption(caption)
	surface.Begin()
	t.cmd.PaintEvent(t.cmd, t.pid, t.context, surface)
	surface.End()
}

type TaskSelector struct {
	pid       int
	available []int
	idx       int
}

func NewTaskSelector() *TaskSelector {
	return &TaskSelector{
		pid:       adaptiveticker.UnknownId,
		available: nil,
		idx:       0,
	}
}

type TaskManager struct {
	ticker     *adaptiveticker.AdaptiveTicker
	foreground *Task
	selector   *TaskSelector
	root       *cli.Command
	path       string
	dirty      bool
	width      int
	height     int
	fullPaint  bool
	timersChan chan *adaptiveticker.TimerHandler
	ids        *adaptiveticker.Ids
}

func NewTaskManager(ticker *adaptiveticker.AdaptiveTicker, timersChannel chan *adaptiveticker.TimerHandler, root *cli.Command) *TaskManager {
	t := &TaskManager{
		ticker:     ticker,
		foreground: nil,
		selector:   NewTaskSelector(),
		timersChan: timersChannel,
		root:       root,
		dirty:      false,
		fullPaint:  true,
		width:      132,
		height:     40,
		ids:        adaptiveticker.NewIds(1024),
	}

	return t
}

func (c *TaskManager) Execute(line string, template *Task) bool {
	args := strings.Split(line, " ")
	if len(args) <= 0 {
		return false
	}

	c.root.SetArgs(args)

	pCmd, flags, err := c.root.Prepare()
	if err != nil {
		return false
	}

	if pCmd == nil {
		return false
	}

	task, err := c.create(pCmd, line)
	if err != nil {
		return false
	}

	if template != nil {
		task.OffsetY = template.OffsetY
		task.OffsetX = template.OffsetX
		task.Scale = template.Scale
	}

	if err = c.root.Execute(task.cmd, flags, task.pid); err == nil {
		if task.cmd.Activate {
			if !task.cmd.Background {
				c.foreground = task
			}
			task.state = taskStateRunning
		}
	}

	if task.state == taskStateSetup {
		c.Kill(task.pid)
	}

	return true
}

func (c *TaskManager) create(cmd *cli.Command, line string) (*Task, error) {
	task := NewTask(cmd, line)

	c.ids.Set(task)
	if task.pid == adaptiveticker.UnknownId {
		return nil, errors.New("max slot reached")
	}

	//c.tasks[task.pid] = task

	return task, nil
}

func (c *TaskManager) SetScreenSize(width int, height int) {
	c.width = width
	c.height = height
	c.fullPaint = true

	//if fgPid := c.GetForegroundPid(); fgPid > unknownId {
	//	c.PaintRequest()
	//}
}

func (c *TaskManager) GetScreenSize() (int, int) {
	return c.width, c.height
}

func (c *TaskManager) SetSelectionMode(requestedPid int) {
	var idx = 0
	var firstPid = adaptiveticker.UnknownId
	var firstIdx = 0

	c.selector.idx = 0
	c.selector.pid = adaptiveticker.UnknownId
	c.selector.available = nil

	for e := c.ids.All().Front(); e != nil; e = e.Next() {
		task := e.Value.(*Task)
		if task.cmd.PaintEvent != nil {
			c.selector.available = append(c.selector.available, task.pid)
			if firstPid == adaptiveticker.UnknownId {
				firstPid = task.pid
				firstIdx = idx
			}
			if task.pid == requestedPid {
				c.selector.pid = requestedPid
				c.selector.idx = idx
			}
			idx++
		}
	}

	if c.selector.pid == adaptiveticker.UnknownId {
		if firstPid == adaptiveticker.UnknownId {
			return
		}
		c.selector.pid = firstPid
		c.selector.idx = firstIdx
	}

	c.PaintRequest()
}

func (c *TaskManager) SetSelectionModeNext() {
	if len(c.selector.available) == 0 {
		return
	}
	next := c.selector.idx + 1
	if next >= len(c.selector.available) {
		next = 0
	}
	c.selector.idx = next
	c.selector.pid = c.selector.available[next]

	c.PaintRequest()
}

func (c *TaskManager) SetSelectionModePrevious() {
	if len(c.selector.available) == 0 {
		return
	}
	prev := c.selector.idx - 1
	if prev < 0 {
		prev = len(c.selector.available) - 1
	}
	c.selector.idx = prev
	c.selector.pid = c.selector.available[prev]

	c.PaintRequest()
}

func (c *TaskManager) SetSelectionDisabled() {
	c.selector.idx = 0
	c.selector.pid = adaptiveticker.UnknownId
	c.selector.available = nil

	c.PaintRequest()
}

func (c *TaskManager) SetBasePath(arg string) bool {
	c.path = arg

	/*
		var cmd * cli.Command

		arg = strings.Trim(arg, " ")

		if arg == ".." {
			if c.cmd.Parent() != nil {
				cmd = c.cmd.Parent()
			}
		} else {
			path := strings.Split(arg, "/")
			if len(path) > 0 {
				var e error
				fmt.Println("TRAVERSING", path)
				if cmd, _, e = c.cmd.Traverse(path); e != nil {
					cmd = nil
				} else {
					fmt.Println("NEW CMD", cmd)
				}
			}
		}

		if cmd != nil {
			c.cmd = cmd
		}
	*/

	return true
}

func (c *TaskManager) SetSelectionOptions(option rune, value float64) bool {
	t, ok := c.ids.Get(c.selector.pid)
	if !ok {
		return false
	}

	task := t.(*Task)

	switch option {
	case 'y':
		task.OffsetY += int(value)
	case 'x':
		task.OffsetX += int(value)
	case 'z':
		if scale := task.Scale + value; scale >= 0.2 && scale <= 1 {
			task.Scale = scale
		}
	}

	c.fullPaint = true

	c.PaintRequest()

	return true
}

func (c *TaskManager) SetCaption(pid int, caption string) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	task.caption = caption
	return true
}

func (c *TaskManager) SetContext(pid int, context interface{}) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	task.context = context
	return true
}

func (c *TaskManager) SetFg(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	c.foreground = task
	return true
}

func (c *TaskManager) GetSuggestion(in string, _ int) (string, []string, bool) {
	var data string
	var cmd *cli.Command = nil

	args := strings.Split(in, " ")
	if len(args) == 1 {
		data = args[0]
		cmd = c.root
	} else if len(args) > 1 {
		var e error
		data = args[len(args)-1]
		args = args[:len(args)-1]

		if cmd, _, e = c.root.Traverse(args); e != nil {
			cmd = nil
		}
		if cmd == c.root {
			cmd = nil
		}
	}

	if cmd == nil {
		return data, nil, false
	}

	return data, cmd.SuggestionsFor(data), true
}

func (c *TaskManager) PaintRequest() bool {
	ret := false
	if !c.dirty {
		c.dirty = true
		ret = true
		c.ticker.Create(c.timersChan, newMessagePaint(), -1, -1, 1)
	}
	return ret
}

func (c *TaskManager) CreateTimer(pid int, first int, interval int, count int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	if task.cmd.TimerEvent == nil {
		return false
	}

	m := newMessageTimer(pid, interval)

	m.tid = c.ticker.Create(c.timersChan, m, int64(first), int64(interval), int64(count))
	if m.tid > -1 {
		task.timers = append(task.timers, m.tid)
	}

	return true
}

func (c *TaskManager) StopTimer(pid int, tid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	return c.closeTimer(task, tid)
}

func (c *TaskManager) IsActive(pid int) bool {
	_, ret := c.ids.Get(pid)
	return ret
}

func (c *TaskManager) GetForegroundPid() int {
	pid := adaptiveticker.UnknownId
	if c.foreground != nil {
		pid = c.foreground.pid
	}
	return pid
}

func (c *TaskManager) GetForegroundName() (int, string) {
	var name string
	pid := adaptiveticker.UnknownId

	if c.foreground != nil {
		pid = c.foreground.pid
		name = c.foreground.cmd.Name()
	}
	return pid, name
}

func (c *TaskManager) SetBackground() bool {
	ret := false
	if c.foreground != nil {
		c.foreground = nil
		ret = true
	}
	return ret
}

func (c *TaskManager) KillForeground() {
	if c.foreground != nil {
		c.Kill(c.foreground.pid)
	}
}

func (c *TaskManager) Kill(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}

	task := t.(*Task)

	if len(task.timers) > 0 {
		c.ticker.Remove(task.timers)
	}

	if c.foreground != nil {
		if c.foreground.pid == pid {
			c.foreground = nil
		}
	}

	c.ids.Unset(pid)

	return true
}

func (c *TaskManager) KillAll(name string) int {
	count := 0
	var tasks []*Task
	for e := c.ids.All().Front(); e != nil; e = e.Next() {
		tasks = append(tasks, e.Value.(*Task))
	}

	for _, task := range tasks {
		deactivate := false
		if len(name) == 0 {
			deactivate = true
		} else {
			if task.cmd.Name() == name {
				deactivate = true
			}
		}

		if deactivate {
			if ok := c.Kill(task.pid); ok {
				count++
			}
		}
	}
	return count
}

func (c *TaskManager) List() string {
	out := "\r\nPid: Task"
	for e := c.ids.All().Front(); e != nil; e = e.Next() {
		task := e.Value.(*Task)
		if task != nil {
			out += fmt.Sprintf("\r\n%d: %s", task.pid, task.cmd.Name())
		}
	}
	return out
}

func (c *TaskManager) ExecTimer(pid int, tid int, interval int) bool {
	ret := false
	if t, ok := c.ids.Get(pid); ok {
		task := t.(*Task)
		if task.cmd.TimerEvent != nil {
			task.cmd.TimerEvent(task.cmd, task.pid, tid, task.context, interval)
			ret = true
		}
	}

	return ret
}

func (c *TaskManager) ExecRead(pid int, code int, buffer rune) bool {
	ret := false
	if t, ok := c.ids.Get(pid); ok {
		task := t.(*Task)
		if task.cmd.ReadEvent != nil {
			task.cmd.ReadEvent(task.cmd, task.pid, task.context, code, buffer)
			ret = true
		}
	}
	return ret
}

func (c *TaskManager) ExecPaint(terminal interfaces.ITerminal) bool {
	if !c.dirty {
		return false
	}

	w, h := c.GetScreenSize()
	surface := newSurface(terminal, h, w)

	if c.fullPaint {
		surface.SetCompletePaint()
		c.fullPaint = false
	}

	var selectedTask *Task = nil

	//TODO Z-INDEX!!!!!

	/*
		for _, pid := range c.ids.All() {
			if t, ok := c.ids.Get(pid); ok {
				task := t.(*Task)
				if task.pid == c.selector.pid {
					selectedTask = task
				} else {
					surface.SetSelectionMode(false)
					task.Paint(surface)
				}
			}
		}
	*/

	for e := c.ids.All().Front(); e != nil; e = e.Next() {
		task := e.Value.(*Task)
		if task.pid == c.selector.pid {
			selectedTask = task
		} else {
			surface.SetSelectionMode(false)
			task.Paint(surface)
		}
	}

	if selectedTask != nil {
		surface.SetSelectionMode(true)
		selectedTask.Paint(surface)
	}

	surface.Render()

	c.dirty = false

	return true
}

func (c *TaskManager) ListTasks() []string {
	var out []string
	dir := "./"
	if files, err := ioutil.ReadDir(dir); err == nil {
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			file := f.Name()
			pos := strings.LastIndex(file, tasksFileExtension)
			if pos < 0 {
				continue
			}

			out = append(out, file[:pos])
		}
	}
	return out
}

func (c *TaskManager) SaveTasks(name string) bool {
	var tasks map[int]*Task
	tasks = make(map[int]*Task)

	for e := c.ids.All().Front(); e != nil; e = e.Next() {
		task := e.Value.(*Task)
		if strings.HasPrefix(task.Line, commandTask) {
			continue
		}
		tasks[task.pid] = task
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		log.Println("Error marshalling task file ", name, ": ", err.Error())
		return false
	}

	if pos := strings.LastIndex(name, string(os.PathSeparator)); pos > -1 {
		name = name[pos+1:]
	}

	name += tasksFileExtension

	if err = ioutil.WriteFile(name, data, 0644); err != nil {
		log.Println("Error writing task file ", name, ": ", err.Error())
		return false
	}

	return true
}

func (c *TaskManager) RestoreTasks(name string) bool {
	var tasks map[int]*Task

	if pos := strings.LastIndex(name, string(os.PathSeparator)); pos > -1 {
		name = name[pos+1:]
	}
	name += tasksFileExtension

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return false
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return false
	}

	for _, task := range tasks {
		if strings.HasPrefix(task.Line, commandTask) {
			continue
		}
		c.Execute(task.Line, task)
	}

	c.Execute(commandActivate, nil)

	return true
}

func (c *TaskManager) ExecActivate() bool {
	pid, name := c.GetForegroundName()

	if pid == adaptiveticker.UnknownId {
		return false
	}
	if name == commandActivate {
		return false
	}

	c.SetBackground()
	c.Execute(fmt.Sprint(commandActivate, " ", pid), nil)

	return false
}

func (c *TaskManager) closeTimer(task *Task, tid int) bool {
	ret := false
	if task != nil {
		for _, timer := range task.timers {
			if timer == tid {
				ret = c.ticker.Remove([]int{timer})
				break
			}
		}
	}
	return ret
}
