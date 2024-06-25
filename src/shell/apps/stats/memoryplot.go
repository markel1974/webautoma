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

package stats

import (
	"github.com/markel1974/webautoma/src/shell/apps/commandcreator"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"math"
	"runtime"
)

type rtPlotData struct {
	rtPlotData   []float64
	rtPlotType   int
	rtPlotMinVal float64
	rtPlotMaxVal float64
	rtPlotAuto   bool
}

func CreateMemoryPlot(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "rtplot"
	root.Short = "Runtime Plot"
	root.Long = "Runtime Plot"
	root.Activate = true
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()

		plt := &rtPlotData{
			rtPlotType:   0,
			rtPlotAuto:   true,
			rtPlotData:   nil,
			rtPlotMinVal: math.Inf(1),
			rtPlotMaxVal: math.Inf(-1),
		}

		if len(args) > 0 {
			switch args[0] {
			case "alloc":
				plt.rtPlotType = 0
			case "total":
				plt.rtPlotType = 1
			case "os":
				plt.rtPlotType = 2
			case "gc":
				plt.rtPlotType = 3
			}
		}
		r.SetContext(pid, plt)
		r.CreateTimer(pid, 0, 300, -1)
	}
	root.ReadEvent = func(cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
		plt := ctx.(*rtPlotData)

		interval := math.Abs(plt.rtPlotMaxVal - plt.rtPlotMinVal)
		scale := (interval * 10) / 100

		switch key {
		case 'a', '+':
			plt.rtPlotAuto = false
			plt.rtPlotMaxVal += scale
			plt.rtPlotMinVal -= scale
		case 'z', '-':
			plt.rtPlotAuto = false
			plt.rtPlotMaxVal -= scale
			plt.rtPlotMinVal += scale
		case 'r':
			plt.rtPlotAuto = !plt.rtPlotAuto
		}

	}
	root.TimerEvent = func(cmd *cli.Command, pid int, tid int, ctx interface{}, interval int) {
		var r = cmd.GetRootContext()
		var m runtime.MemStats
		plt := ctx.(*rtPlotData)

		runtime.ReadMemStats(&m)
		var val float64
		switch plt.rtPlotType {
		case 0:
			val = bToMb(m.Alloc)
		case 1:
			val = bToMb(m.TotalAlloc)
		case 2:
			val = bToMb(m.Sys)
		case 3:
			val = float64(m.NumGC)
		default:
			val = bToMb(m.Alloc)
		}

		if val < plt.rtPlotMinVal {
			plt.rtPlotMinVal = val
		}
		if val > plt.rtPlotMaxVal {
			plt.rtPlotMaxVal = val
		}

		plt.rtPlotData = append(plt.rtPlotData, val)
		if len(plt.rtPlotData) > 10 {
			plt.rtPlotData = plt.rtPlotData[1:]
		}

		r.PaintRequest(pid)
	}
	root.PaintEvent = func(cmd *cli.Command, pid int, ctx interface{}, surface interfaces.ISurface) {
		var min float64 = 0
		var max float64 = 0
		plt := ctx.(*rtPlotData)
		if !plt.rtPlotAuto {
			min = plt.rtPlotMinVal
			max = plt.rtPlotMaxVal
		}

		surface.DrawSeries(plt.rtPlotData, -1, -1, min, max)
	}

	return root
}
