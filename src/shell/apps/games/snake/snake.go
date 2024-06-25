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

package snake

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"math/rand"
	"strconv"
)

type direction int

const (
	Up direction = iota
	Down
	Left
	Right
)

type Point struct {
	X int
	Y int
}

type Attributes struct {
	C  rune
	Fg interfaces.ColorDef
	Bg interfaces.ColorDef
}

type Data struct {
	Point
	Attributes
}

type Snake struct {
	Position   Point
	Direction  direction
	Length     int
	Body       []Data
	Speed      int
	Rows       int
	Columns    int
	Food       Data
	Score      int
	GameOver   bool
	borders    []*Border
	SpriteHead Attributes
	SpriteBody Attributes
	SpriteTail Attributes
}

func New() *Snake {
	snake := &Snake{
		SpriteTail: Attributes{C: '░', Fg: interfaces.ColorYellowDef, Bg: interfaces.ColorNoneDef},
		SpriteBody: Attributes{C: '░', Fg: interfaces.ColorMagentaDef, Bg: interfaces.ColorNoneDef},
		SpriteHead: Attributes{C: '@', Fg: interfaces.ColorRedDef, Bg: interfaces.ColorNoneDef},
	}
	return snake
}

func (snake *Snake) SetSize(rows int, column int) {
	snake.Rows = rows
	snake.Columns = column
	snake.borders = nil

	snake.Start()
}

func (snake *Snake) Start() {
	snake.GameOver = false
	snake.Speed = 1
	snake.Direction = Right
	snake.Body = []Data{
		{Point: Point{X: 1, Y: 6}, Attributes: snake.SpriteTail}, // Tail
		{Point: Point{X: 2, Y: 6}, Attributes: snake.SpriteBody}, // Body
		{Point: Point{X: 3, Y: 6}, Attributes: snake.SpriteHead}, // Head
	}
	snake.moveFood()
	snake.borders = nil
	snake.borders = append(snake.borders, NewBorder(0, 1, snake.Columns, snake.Rows))
	//snake.borders = append(snake.borders, NewBorder(Coordinates{7, 10}, 3, 5))
	//snake.borders = append(snake.borders, NewBorder(Coordinates{32, 5}, 3, 20))
}

func (snake *Snake) moveFood() {
	minRows := 1
	maxRows := snake.Rows - 2
	minColumns := 1
	maxColumns := snake.Columns - 2

	mr := maxRows - minRows
	if mr <= 0 {
		mr = 1
	}

	mc := maxColumns - minColumns
	if mc <= 0 {
		mc = 1
	}

	snake.Food.Y = rand.Intn(mr) + minRows
	snake.Food.X = rand.Intn(mc) + minColumns
	snake.Food.C = '*'
	snake.Food.Fg = interfaces.ColorRedDef
	snake.Food.Bg = interfaces.ColorNoneDef
}

func (snake *Snake) createSurface() [][]string {
	var surface [][]string
	for i := 0; i < snake.Rows; i++ {
		var line []string
		for j := 0; j < snake.Columns; j++ {
			line = append(line, " ")
		}
		surface = append(surface, line)
	}
	return surface
}

func (snake *Snake) head() Point {
	if len(snake.Body) > 0 {
		var idx = len(snake.Body) - 1
		if idx >= 0 && idx < len(snake.Body) {
			return snake.Body[idx].Point
		}
	}
	return Point{}
}

func (snake *Snake) borderCollision() bool {
	k := snake.head()
	found := false
	for _, border := range snake.borders {
		if border.Contains(k) {
			found = true
			break
		}
	}
	return found
}

func (snake *Snake) foodCollision() bool {
	k := snake.head()
	if k.X == snake.Food.X && k.Y == snake.Food.Y {
		return true
	}
	return false
}

func (snake *Snake) snakeCollision() bool {
	return snake.Contains()
}

func (snake *Snake) Advance() {
	//for x := 0; x <= snake.Score; x ++ {
	snake.update()
	//}
}

func (snake *Snake) update() {
	if !snake.GameOver {
		nHead := snake.head()
		switch snake.Direction {
		case Up:
			nHead.Y--
		case Down:
			nHead.Y++
		case Left:
			nHead.X--
		case Right:
			nHead.X++
		}

		if snake.borderCollision() || snake.snakeCollision() {
			snake.GameOver = true
		} else {
			if snake.foodCollision() {
				snake.Score++
				snake.rebuildSnake(nHead)
				snake.moveFood()
			} else {
				snake.Body = snake.Body[1:]
				snake.rebuildSnake(nHead)
			}

			snake.Position.X = nHead.X
			snake.Position.Y = nHead.Y
		}
	}
}

func (snake *Snake) rebuildSnake(nHead Point) {
	snake.Body[len(snake.Body)-1].Attributes = snake.SpriteBody
	snake.Body = append(snake.Body, Data{Point: nHead, Attributes: snake.SpriteHead})
	snake.Body[0].Attributes = snake.SpriteTail
}

func (snake *Snake) Draw(surface interfaces.ISurface) {
	surface.DrawColor(snake.Food.Y, snake.Food.X, snake.Food.C, snake.Food.Fg, snake.Food.Bg, interfaces.ModeNormal)

	for _, c := range snake.Body {
		surface.DrawColor(c.Y, c.X, c.C, c.Fg, c.Bg, interfaces.ModeNormal)
	}

	for _, border := range snake.borders {
		border.Draw(surface)
	}

	score := "Score: " + strconv.Itoa(snake.Score)

	surface.DrawTextColor(0, 0, score, interfaces.ColorBlueDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	if snake.GameOver {
		rows, column := surface.GetSize()
		gameOver := "Game Over"
		surface.DrawTextColor(rows/2, (column/2)-(len(gameOver)/2), gameOver, interfaces.ColorBlueDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	}
}

func (snake *Snake) SetPosition(x int, _ int) {
	snake.Position.X = x
	snake.Position.Y = x
}

func (snake *Snake) Contains() bool {
	for i := 0; i < len(snake.Body)-1; i++ {
		if snake.head() == snake.Body[i].Point {
			return true
		}
	}
	return false
}

type Border struct {
	width  int
	height int
	coords map[Point]Data
}

func NewBorder(origX int, origY int, width, height int) *Border {
	b := new(Border)
	b.width, b.height = width-1, height-2
	b.coords = make(map[Point]Data)

	for x := origX; x < b.width; x++ {
		p1 := Point{X: x, Y: origY}
		p2 := Point{X: x, Y: b.height}

		b.coords[p1] = Data{Point: p1, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
		b.coords[p2] = Data{Point: p2, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
	}

	for y := origY; y < b.height+1; y++ {
		p1 := Point{X: origX, Y: y}
		p2 := Point{X: b.width, Y: y}

		b.coords[p1] = Data{Point: p1, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
		b.coords[p2] = Data{Point: p2, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
	}

	return b
}

func (b *Border) Contains(point Point) bool {
	_, exists := b.coords[point]
	return exists
}

func (b *Border) Draw(s interfaces.ISurface) {
	for _, c := range b.coords {
		s.DrawColor(c.Y, c.X, c.C, c.Fg, c.Bg, interfaces.ModeNormal)
	}
}
