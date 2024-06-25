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

type Cell struct {
	x, y  int
	color interfaces.ColorDef
}

func NewCell(x, y int, ch rune) *Cell {
	return &Cell{x: x, y: y, color: colorMapping[ch]}
}

func (c *Cell) conflicts(board *Board) bool {
	return c.isOnWall(board) || c.isOverlapped(board)
}

func (c *Cell) isOverlapped(board *Board) bool {
	if !board.isOnBoard(c.x, c.y) {
		return false
	}
	return board.colors[c.x][c.y] != blankColor
}

func (c *Cell) isOnWall(board *Board) bool {
	return c.x < 0 || board.w <= c.x || board.h <= c.y
}
