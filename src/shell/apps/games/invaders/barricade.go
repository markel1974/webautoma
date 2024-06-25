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

package invaders

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"strings"
)

const (
	barricadeSymbol       = '▓'
	numBarricades         = 4
	fgBarricade           = interfaces.ColorWhiteDef
	bgBarricade           = interfaces.ColorGreenDef
	barricadeSpriteWidth  = 11
	barricadeSpriteHeight = 5
	barricadeSprite       = `    xxx
  xxxxxxx
xxxxxxxxxxx
xxx     xxx
xxx     xxx`
)

func barricadeYPos(h int) int {
	//TODO pass playerSpriteHeight
	var playerSpriteHeight = 3
	return (h - playerSpriteBottomOffset - playerSpriteHeight) - barricadeSpriteHeight - 2
}

type Barricade struct {
	w    int
	h    int
	data [][]int
}

func NewBarricade() *Barricade {
	b := &Barricade{}
	return b
}

func (b *Barricade) Setup(w int, h int) *Barricade {
	b.w = w
	b.h = h

	b.data = make([][]int, b.w)
	for i := range b.data {
		b.data[i] = make([]int, b.h)
		for j := range b.data[i] {
			b.data[i][j] = nonIndex
		}
	}

	gap := (w - 4*barricadeSpriteWidth) / 5
	x := gap
	y := barricadeYPos(h)
	initY := y
	lines := strings.Split(barricadeSprite, "\n")
	for i := 0; i < numBarricades; i++ {
		initX := gap*(i+1) + barricadeSpriteWidth*i
		x = initX
		for _, l := range lines {
			for _, c := range l {
				if c != ' ' {
					b.Set(x, y, i)
				}
				x++
			}
			y++
			x = initX
		}
		y = initY
	}
	return b
}

func (b *Barricade) Unset(x int, y int) bool {
	var ret = false
	if b.Get(x, y) != nonIndex {
		b.Set(x, y, nonIndex)
		ret = true
	}

	return ret
}

func (b *Barricade) Draw(surface interfaces.ISurface) {
	for i := range b.data {
		for j := range b.data[i] {
			if b.data[i][j] != nonIndex {
				surface.DrawColor(j, i, barricadeSymbol, fgBarricade, bgBarricade, interfaces.ModeNormal)
			}
		}
	}
}

func (b *Barricade) Set(x int, y int, data int) {
	if len(b.data) > 0 {
		if x >= 0 && y >= 0 && x < b.w && y < b.h {
			b.data[x][y] = data
		}
	}
}

func (b *Barricade) Get(x int, y int) int {
	data := nonIndex
	if len(b.data) > 0 {
		if x >= 0 && y >= 0 && x < b.w && y < b.h {
			data = b.data[x][y]
		}
	}
	return data
}
