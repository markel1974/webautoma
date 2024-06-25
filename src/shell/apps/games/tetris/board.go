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

package tetris

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
)

type Board struct {
	colors [][]interfaces.ColorDef
	w      int
	h      int
}

func NewBoard(w int, h int) *Board {
	b := &Board{
		w: w,
		h: h,
	}
	b.colors = make([][]interfaces.ColorDef, w)
	for r := range b.colors {
		b.colors[r] = make([]interfaces.ColorDef, h)
		for c := range b.colors[r] {
			b.colors[r][c] = blankColor
		}
	}
	return b
}

func (b *Board) deleteLine(y int) {
	for i := 0; i < b.w; i++ {
		b.colors[i][y] = blankColor
	}
	for j := y; j > 0; j-- {
		for i := 0; i < b.w; i++ {
			b.colors[i][j] = b.colors[i][j-1]
		}
	}
	for i := 0; i < b.w; i++ {
		b.colors[i][0] = blankColor
	}
}

func (b *Board) fullLines() []int {
	var fullLines []int
	for j := 0; j < b.h; j++ {
		if b.isFullLine(j) {
			fullLines = append(fullLines, j)
		}
	}
	return fullLines
}

func (b *Board) isFullLine(y int) bool {
	hasBlank := false
	for i := 0; i < b.w; i++ {
		if b.colors[i][y] == blankColor {
			hasBlank = true
			break
		}
	}
	return !hasBlank
}

func (b *Board) hasFullLine() bool {
	for j := 0; j < b.h; j++ {
		if b.isFullLine(j) {
			return true
		}
	}
	return false
}

/*
func (b *Board) text() string {
	text := ""
	for j := 0; j < b.h; j++ {
		for i := 0; i < b.w; i++ {
			text = fmt.Sprintf("%s%c", text, charByColor(b.colors[i][j]))
		}
		text = fmt.Sprintf("%s\n", text)
	}
	return text
}
*/

func (b *Board) setCell(cell *Cell) {
	b.colors[cell.x][cell.y] = cell.color
}

func (b *Board) setCells(cells []*Cell) {
	for _, cell := range cells {
		b.setCell(cell)
	}
}

func (b *Board) isOnBoard(x, y int) bool {
	return (0 <= x && x < b.w) && (0 <= y && y < b.h)
}
