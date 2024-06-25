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
	"math/rand"
	"strings"
)

const (
	defaultTetrominoX, defaultTetrominoY = 3, -1
	TetrominoWidth, TetrominoHeight      = 4, 4
)

var (
	srcBlocks = []string{
		`
			....
			.GG.
			GG..
			....
		`, `
			....
			.RR.
			..RR
			....
		`, `
			....
			.YY.
			.YY.
			....
		`, `
			....
			....
			CCCC
			....
		`, `
			....
			.M..
			MMM.
			....
		`, `
			....
			.b..
			.bbb
			....
		`, `
			....
			..m.
			mmm.
			....
		`,
	}
)

var blocks = initializeBlocks(srcBlocks)

func initializeBlocks(blocks []string) []string {
	var out []string
	for _, block := range blocks {
		block = strings.Replace(block, "\t", "", -1)
		block = strings.Replace(block, " ", "", -1)
		block = strings.Trim(block, "\n")
		out = append(out, block)
	}
	return out
}

type Tetromino struct {
	block string
	x     int
	y     int
}

func NewMino() *Tetromino {
	return &Tetromino{
		block: blocks[rand.Intn(len(blocks))],
	}
}

func (m *Tetromino) cell(x, y int) rune {
	return rune(m.block[x+(TetrominoWidth+1)*y])
}

func (m *Tetromino) setCell(x, y int, cell rune) {
	buf := []rune(m.block)
	buf[x+(TetrominoWidth+1)*y] = cell
	m.block = string(buf)
}

func (m *Tetromino) putBottom(board *Board) int {
	distance := -1
	dstMino := *m
	for !dstMino.conflicts(board) {
		*m = dstMino
		dstMino.forceMoveDown()
		distance++
	}
	if distance < 0 {
		distance = 0
	}
	return distance
}

func (m *Tetromino) forceMoveDown() {
	m.y++
}

func (m *Tetromino) forceRotateRight() {
	oldMino := *m
	for j := 0; j < TetrominoHeight; j++ {
		for i := 0; i < TetrominoWidth; i++ {
			m.setCell(TetrominoHeight-j-1, i, oldMino.cell(i, j))
		}
	}
}

func (m *Tetromino) forceRotateLeft() {
	oldMino := *m
	for j := 0; j < TetrominoHeight; j++ {
		for i := 0; i < TetrominoWidth; i++ {
			m.setCell(j, TetrominoWidth-i-1, oldMino.cell(i, j))
		}
	}
}

func (m *Tetromino) conflicts(board *Board) bool {
	for _, cell := range m.cells() {
		if cell.conflicts(board) {
			return true
		}
	}
	return false
}

func (m *Tetromino) cells() []*Cell {
	var cells []*Cell
	for i := 0; i < TetrominoWidth; i++ {
		for j := 0; j < TetrominoHeight; j++ {
			if m.cell(i, j) != '.' {
				cells = append(cells, NewCell(m.x+i, m.y+j, m.cell(i, j)))
			}
		}
	}
	return cells
}

func (m *Tetromino) lines() []string {
	return strings.Split(m.block, "\n")
}
