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

package plotter

import (
	"fmt"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"math"
)

type Plotter struct {
	maximum   float64
	minimum   float64
	precision int
	ratio     float64
	interval  float64
	maxWidth  int
	min2      float64
	max2      float64
	series    []float64

	offset int

	Width  int
	Height int
	Rows   int
	Column int
}

func NewPlotter(w int, h int) *Plotter {
	return &Plotter{
		ratio:     1,
		precision: 2,
		Width:     w,
		Height:    h,
		offset:    3,
	}
}

func (p *Plotter) minMaxFloat64Slice(v []float64) (min, max float64) {
	min = math.Inf(1)
	max = math.Inf(-1)

	if len(v) > 0 {
		for _, e := range v {
			if e < min {
				min = e
			}
			if e > max {
				max = e
			}
		}
	}
	return
}

func (p *Plotter) round(input float64) float64 {
	if math.IsNaN(input) {
		return math.NaN()
	}
	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}
	_, decimal := math.Modf(input)
	var rounded float64
	if decimal >= 0.5 {
		rounded = math.Ceil(input)
	} else {
		rounded = math.Floor(input)
	}
	return rounded * sign
}

func (p *Plotter) linearInterpolate(before, after, atPoint float64) float64 {
	return before + (after-before)*atPoint
}

func (p *Plotter) interpolateArray(data []float64, width int) []float64 {
	var interpolatedData []float64
	springFactor := float64(len(data)-1) / float64(width-1)
	interpolation := data[0]

	interpolatedData = append(interpolatedData, interpolation)

	for i := 1; i < width-1; i++ {
		spring := float64(i) * springFactor
		before := math.Floor(spring)
		after := math.Ceil(spring)
		atPoint := spring - before
		interpolation = p.linearInterpolate(data[int(before)], data[int(after)], atPoint)
		interpolatedData = append(interpolatedData, interpolation)
	}

	interpolation = data[len(data)-1]
	interpolatedData = append(interpolatedData, interpolation)

	return interpolatedData
}

func (p *Plotter) Setup(series []float64, inputMinimum float64, inputMaximum float64) {
	if len(series) == 0 {
		return
	}
	p.series = series
	p.minimum = inputMinimum
	p.maximum = inputMaximum
	if p.minimum >= p.maximum {
		p.minimum, p.maximum = p.minMaxFloat64Slice(p.series)
	}

	p.interval = math.Abs(p.maximum - p.minimum)

	logMaximum := math.Log10(math.Max(math.Abs(p.maximum), math.Abs(p.minimum))) //to find number of zeros after decimal
	if p.minimum == float64(0) && p.maximum == float64(0) {
		logMaximum = float64(-1)
	}

	if logMaximum < 0 {
		if math.Mod(logMaximum, 1) != 0 {
			p.precision = p.precision + int(math.Abs(logMaximum))
		} else {
			p.precision = p.precision + int(math.Abs(logMaximum)-1.0)
		}
	} else if logMaximum > 2 {
		p.precision = 0
	}

	maxNumLength := len(fmt.Sprintf("%0.*f", p.precision, p.maximum))
	minNumLength := len(fmt.Sprintf("%0.*f", p.precision, p.minimum))
	p.maxWidth = int(math.Max(float64(maxNumLength), float64(minNumLength)))

	if p.Width > 0 {
		w := p.maxWidth + 3
		if p.Width > w {
			p.Width = p.Width - w
		}
		p.series = p.interpolateArray(p.series, p.Width)
	}

	if p.Height <= 0 {
		if int(p.interval) <= 0 {
			p.Height = int(p.interval * math.Pow10(int(math.Ceil(-math.Log10(p.interval)))))
		} else {
			p.Height = int(p.interval)
		}
	}

	if p.interval != 0 {
		p.ratio = float64(p.Height) / p.interval
	}

	p.min2 = p.round(p.minimum * p.ratio)
	p.max2 = p.round(p.maximum * p.ratio)

	intMin2 := int(math.Floor(p.min2))
	intMax2 := int(math.Ceil(p.max2))

	p.Rows = intMax2 - intMin2
	//p.Rows = int(math.Abs(p.max2 - p.min2))
	p.Column = len(p.series) + p.offset
}

func (p *Plotter) Draw(surface interfaces.ISurface) {
	if len(p.series) == 0 {
		return
	}
	//intMin2 := int(p.min2)
	//intMax2 := int(p.max2)

	intMin2 := int(math.Floor(p.min2))
	intMax2 := int(math.Ceil(p.max2))

	lwidth := p.maxWidth + 1

	for y := intMin2; y < intMax2+1; y++ {
		var magnitude float64
		if p.Rows > 0 {
			magnitude = p.maximum - (float64(y-intMin2) * p.interval / float64(p.Rows))
		} else {
			magnitude = float64(y)
		}

		label := fmt.Sprintf("%*.*f", p.maxWidth+1, p.precision, magnitude)
		ly := y - intMin2
		lx := int(math.Max(float64(p.offset)-float64(lwidth), 0))

		surface.DrawTextColor(ly, lx, label, interfaces.ColorYellowDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
		if y == 0 {
			surface.DrawColor(ly, p.offset-1+lwidth, '┼', interfaces.ColorMagentaDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
		} else {
			surface.DrawColor(ly, p.offset-1+lwidth, '┤', interfaces.ColorMagentaDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
		}
	}
	p.offset += lwidth

	var y0 = int(p.round(p.series[0]*p.ratio) - p.min2)
	var y1 int

	surface.DrawColor(p.Rows-y0, p.offset-1, '┼', interfaces.ColorMagentaDef, interfaces.ColorNoneDef, interfaces.ModeNormal)

	for x := 0; x < len(p.series)-1; x++ {
		y0 = int(p.round(p.series[x+0]*p.ratio) - float64(intMin2))
		y1 = int(p.round(p.series[x+1]*p.ratio) - float64(intMin2))
		if y0 == y1 {
			surface.DrawColor(p.Rows-y0, x+p.offset, '─', interfaces.ColorWhiteDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
		} else {
			idy1 := p.Rows - y1
			idy0 := p.Rows - y0
			idx := x + p.offset
			if y0 > y1 {
				surface.DrawColor(idy1, idx, '╰', interfaces.ColorWhiteDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
				surface.DrawColor(idy0, idx, '╮', interfaces.ColorWhiteDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			} else {
				surface.DrawColor(idy1, idx, '╭', interfaces.ColorWhiteDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
				surface.DrawColor(idy0, idx, '╯', interfaces.ColorWhiteDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			}

			start := int(math.Min(float64(y0), float64(y1))) + 1
			end := int(math.Max(float64(y0), float64(y1)))
			for y := start; y < end; y++ {
				surface.DrawColor(p.Rows-y, x+p.offset, '│', interfaces.ColorWhiteDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			}
		}
	}
}

// Plot returns ascii graph for a series.
/*
func PlotLineChart(values []float64, newline string, minimum float64, maximum float64, options ...Option) string {
	var logMaximum float64

	config := configure(config{
		Offset: 3,
	}, options)

	defaultSpacing := 1

	interval := math.Abs(maximum - minimum)

	if config.Height <= 0 {
		if int(interval) <= 0 {
			config.Height = int(interval * math.Pow10(int(math.Ceil(-math.Log10(interval)))))
		} else {
			config.Height = int(interval)
		}
	}

	if config.Offset <= 0 {
		config.Offset = 3
	}

	var ratio float64
	if interval != 0 {
		ratio = float64(config.Height) / interval
	} else {
		ratio = 1
	}
	min2 := p.round(minimum * ratio)
	max2 := p.round(maximum * ratio)

	intmin2 := int(min2)
	intMax2 := int(max2)

	rows := int(math.Abs(float64(intMax2 - intmin2)))
	width := config.Offset + len(values) + (len(values) * defaultSpacing)

	//fmt.Println("intmin2=", intmin2, " intMax2=", intMax2, " rows=", rows, " width=", width, " ratio=", ratio)

	var plot [][]string

	// initialise empty 2D grid
	for i := 0; i < rows+1; i++ {
		var line []string
		for j := 0; j < width; j++ {
			line = append(line, " ")
		}
		plot = append(plot, line)
	}

	precision := 2
	logMaximum = math.Log10(math.Max(math.Abs(maximum), math.Abs(minimum))) //to find number of zeros after decimal
	if minimum == float64(0) && maximum == float64(0) {
		logMaximum = float64(-1)
	}
	//fmt.Println("logMaximum=", logMaximum)

	if logMaximum < 0 {
		// negative log
		if math.Mod(logMaximum, 1) != 0 {
			// non-zero digits after decimal
			precision = precision + int(math.Abs(logMaximum))
		} else {
			precision = precision + int(math.Abs(logMaximum)-1.0)
		}
	} else if logMaximum > 2 {
		precision = 0
	}

	maxNumLength := len(fmt.Sprintf("%0.*f", precision, maximum))
	minNumLength := len(fmt.Sprintf("%0.*f", precision, minimum))
	maxWidth := int(math.Max(float64(maxNumLength), float64(minNumLength)))

	// axis and labels
	for y := intmin2; y < intMax2+1; y++ {
		var magnitude float64
		if rows > 0 {
			magnitude = maximum - (float64(y-intmin2) * interval / float64(rows))
		} else {
			magnitude = float64(y)
		}

		label := fmt.Sprintf("%*.*f", maxWidth+1, precision, magnitude)
		w := y - intmin2
		h := int(math.Max(float64(config.Offset)-float64(len(label)), 0))

		plot[w][h] = label
		plot[w][config.Offset-2] = "┤"
	}

	//Plot chart
	for x := 0; x < len(values); x++ {
		idx := x + config.Offset + (x * defaultSpacing)
		//fmt.Println("PLOT CHART idx=", idx)
		for y := intmin2; y < intMax2+1; y++ {
			fCurr := float64(y) / ratio
			if fCurr <= values[x] {
				//fmt.Println("fCurr=", fCurr, " value=", values[x])
				w := y - intmin2
				plot[rows-w][idx] = "█"
			}
		}
	}

	// join columns
	var lines bytes.Buffer
	for h, horizontal := range plot {
		if h != 0 {
			lines.WriteString(newline)
		}
		for _, v := range horizontal {
			lines.WriteString(v)
		}
	}

	// add caption if not empty
	if len(config.Caption) > 0 {
		captions := strings.Split(config.Caption, " ")
		lines.WriteString(newline)

		for x := 0; x < len(captions); x++ {
			idx := config.Offset
			if x > 0 {
				idx += x * (defaultSpacing / 6)
			}
			//fmt.Println("idx=", idx)
			lines.WriteString(strings.Repeat(" ", idx))
			lines.WriteString(captions[x])
		}
	}

	return lines.String()
}
*/
