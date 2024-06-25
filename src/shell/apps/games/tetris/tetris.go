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
	"fmt"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"math/rand"
	"time"
)

const (
	levelMax = 20
	scoreMax = 999999
)

type Tetris struct {
	board           *Board
	currentMino     *Tetromino
	nextMino        *Tetromino
	score           int
	level           int
	initLevel       int
	deleteLines     int
	gameOver        bool
	backgroundLines []string

	animationGameOverCount int
	animationDeleteCount   int
	animationDelete        []int
}

func New(w int, h int) *Tetris {
	t := &Tetris{}
	t.Init(w, h)
	return t
}

func (t *Tetris) Init(w int, h int) {
	rand.Seed(time.Now().UnixNano())
	t.board = NewBoard(w, h)
	t.level = t.initLevel
	t.score = 0
	t.deleteLines = 0
	t.currentMino = nil
	t.nextMino = nil
	t.gameOver = false
	t.animationGameOverCount = 0
	t.animationDeleteCount = 0
	t.animationDelete = nil
	t.backgroundLines = background

	t.initMino()
}

func (t *Tetris) GetSize() (int, int) {
	return t.board.w, t.board.h
}

func (t *Tetris) initMino() {
	t.currentMino = nil
	t.nextMino = nil
	t.pushMino()
	t.pushMino()
}

func (t *Tetris) deleteCheck() {
	if !t.board.hasFullLine() {
		return
	}

	lines := t.board.fullLines()

	t.animationDeleteCount = 4
	t.animationDelete = lines

	for _, line := range lines {
		t.board.deleteLine(line)
	}
	t.deleteLines += len(lines)
	switch len(lines) {
	case 1:
		t.addScore(40 * (t.level + 1))
	case 2:
		t.addScore(100 * (t.level + 1))
	case 3:
		t.addScore(300 * (t.level + 1))
	case 4:
		t.addScore(1200 * (t.level + 1))
	}
	t.levelUpdate()
}

func (t *Tetris) levelUpdate() {
	if t.level == levelMax {
		return
	}

	targetLevel := t.deleteLines / 10
	if t.level < targetLevel {
		t.level = targetLevel
	}
}

func (t *Tetris) addScore(add int) {
	t.score += add
	if t.score > scoreMax {
		t.score = scoreMax
	}
}

func (t *Tetris) pushMino() {
	t.deleteCheck()

	t.currentMino = t.nextMino
	if t.currentMino != nil {
		t.currentMino.x, t.currentMino.y = defaultTetrominoX, defaultTetrominoY
		if t.currentMino.conflicts(t.board) {
			ranking := NewRanking()
			ranking.insertScore(t.score)
			ranking.save()
			t.gameOver = true
			return
		}
	}
	t.nextMino = NewMino()
}

func (t *Tetris) ApplyGravity() {
	if t.gameOver {
		return
	}
	t.MoveDown()
}

func (t *Tetris) Drop() {
	if t.gameOver {
		return
	}
	t.addScore(t.currentMino.putBottom(t.board))
	t.board.setCells(t.currentMino.cells())
	t.pushMino()
}

func (t *Tetris) MoveDown() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.forceMoveDown()

	if dstMino.conflicts(t.board) {
		t.board.setCells(t.currentMino.cells())
		t.pushMino()
	} else {
		t.currentMino.forceMoveDown()
	}
}

func (t *Tetris) MoveLeft() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.x--
	if !dstMino.conflicts(t.board) {
		t.currentMino.x--
	}
}

func (t *Tetris) MoveRight() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.x++
	if !dstMino.conflicts(t.board) {
		t.currentMino.x++
	}
}

func (t *Tetris) RotateRight() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.forceRotateRight()
	if !dstMino.conflicts(t.board) {
		t.currentMino.forceRotateRight()
	}
}

func (t *Tetris) RotateLeft() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.forceRotateLeft()
	if !dstMino.conflicts(t.board) {
		t.currentMino.forceRotateLeft()
	}
}

func (t *Tetris) Draw(surface interfaces.ISurface) {
	t.drawBackGround(surface, 0, 0)
	t.drawBoard(surface, boardXOffset, boardYOffset)
	t.drawMino(surface, t.nextMino, nextMinoXOffset-t.nextMino.x, nextMinoYOffset-t.nextMino.y)
	t.drawTexts(surface)
	t.drawDropMarker(surface)
	t.drawMino(surface, t.currentMino, boardXOffset, boardYOffset)
	t.drawAnimationDelete(surface)
	t.drawGameOver(surface)
}

func (t *Tetris) drawTexts(surface interfaces.ISurface) {
	surface.DrawTextColor(9, 32, "SCORE", interfaces.ColorWhiteDef, interfaces.ColorBlueDef, interfaces.ModeNormal)
	surface.DrawTextColor(10, 32, fmt.Sprintf("%7d", t.score), interfaces.ColorBlackDef, interfaces.ColorWhiteDef, interfaces.ModeNormal)
	surface.DrawTextColor(13, 32, "LEVEL", interfaces.ColorWhiteDef, interfaces.ColorBlueDef, interfaces.ModeNormal)
	surface.DrawTextColor(14, 32, fmt.Sprintf("%5d", t.level), interfaces.ColorBlackDef, interfaces.ColorWhiteDef, interfaces.ModeNormal)
	surface.DrawTextColor(16, 32, "LINES", interfaces.ColorWhiteDef, interfaces.ColorBlueDef, interfaces.ModeNormal)
	surface.DrawTextColor(17, 32, fmt.Sprintf("%5d", t.deleteLines), interfaces.ColorBlackDef, interfaces.ColorWhiteDef, interfaces.ModeNormal)
}

func (t *Tetris) drawDropMarker(surface interfaces.ISurface) {
	marker := *t.currentMino
	marker.putBottom(t.board)

	for y, line := range marker.lines() {
		for x, char := range line {
			if t.board.isOnBoard(x+marker.x, y+marker.y) && colorByChar(char) != blankColor &&
				colorByChar(char) != interfaces.ColorNoneDef {
				t.drawCell(surface, x+marker.x+boardXOffset, y+marker.y+boardYOffset, colorByChar('K'))
			}
		}
	}
}

func (t *Tetris) drawMino(surface interfaces.ISurface, mino *Tetromino, xOffset, yOffset int) {
	for y, line := range mino.lines() {
		for x, char := range line {
			if t.board.isOnBoard(x+mino.x, y+mino.y) {
				color := colorByChar(char)
				t.drawCell(surface, x+mino.x+xOffset, y+mino.y+yOffset, color)
			}
		}
	}
}

func (t *Tetris) drawBoard(surface interfaces.ISurface, left int, top int) {
	for j := 0; j < t.board.h; j++ {
		for i := 0; i < t.board.w; i++ {
			t.drawCell(surface, left+i, top+j, t.board.colors[i][j])
		}
	}
}

func (t *Tetris) drawCell(surface interfaces.ISurface, x, y int, color interfaces.ColorDef) {
	if color != interfaces.ColorNoneDef && color != blankColor {
		if color == colorByChar('K') {
			surface.DrawColor(y, 2*x-1, '▓', color, interfaces.ColorWhiteDef, interfaces.ModeNormal)
			surface.DrawColor(y, 2*x, ' ', color, interfaces.ColorWhiteDef, interfaces.ModeNormal)
		} else {
			var bg interfaces.ColorDef
			switch color {
			case interfaces.ColorRedDef:
				bg = interfaces.ColorBrightRedDef
			case interfaces.ColorGreenDef:
				bg = interfaces.ColorBrightGreenDef
			case interfaces.ColorYellowDef:
				bg = interfaces.ColorBrightYellowDef
			case interfaces.ColorBlueDef:
				bg = interfaces.ColorBrightBlueDef
			case interfaces.ColorMagentaDef:
				bg = interfaces.ColorBrightMagentaDef
			case interfaces.ColorCyanDef:
				bg = interfaces.ColorBrightCyanDef
			case interfaces.ColorWhiteDef:
				bg = interfaces.ColorBrightWhiteDef
			default:
				bg = color
			}
			surface.DrawColor(y, 2*x-1, '▓', color, bg, interfaces.ModeNormal)
			surface.DrawColor(y, 2*x, ' ', color, bg, interfaces.ModeNormal)
		}
	}
}

func (t *Tetris) drawBackGround(surface interfaces.ISurface, left int, top int) {
	for y, line := range t.backgroundLines {
		for x, char := range line {
			t.drawBack(surface, left+x, top+y, colorByChar(char))
		}
	}
}

func (t *Tetris) drawBack(surface interfaces.ISurface, x, y int, color interfaces.ColorDef) {
	surface.DrawColor(y, 2*x-1, ' ', interfaces.ColorNoneDef, color, interfaces.ModeNormal)
	surface.DrawColor(y, 2*x, ' ', interfaces.ColorNoneDef, color, interfaces.ModeNormal)
}

func (t *Tetris) drawAnimationDelete(surface interfaces.ISurface) {
	if t.animationDeleteCount > 0 {
		for _, line := range t.animationDelete {
			color := interfaces.ColorCyanDef
			if t.animationDeleteCount%2 == 0 {
				color = interfaces.ColorMagentaDef
			}
			t.colorizeLine(surface, line, color)
		}
		t.animationDeleteCount--
	}
}

func (t *Tetris) drawGameOver(surface interfaces.ISurface) {
	if t.gameOver {
		if t.animationGameOverCount < t.board.h {
			for y := t.board.h - 1; y > t.board.h-1-t.animationGameOverCount; y -= 1 {
				t.colorizeLine(surface, y, interfaces.ColorBlackDef)
			}
			t.animationGameOverCount++
			return
		}

		for j := 0; j < t.board.h; j++ {
			t.colorizeLine(surface, j, interfaces.ColorBlackDef)
		}
		surface.DrawTextColor(4, 10, "GAME OVER", interfaces.ColorWhiteDef, interfaces.ColorBlackDef, interfaces.ModeNormal)

		ranking := NewRanking()
		for idx, sc := range ranking.scores {
			fg := availableColors[rand.Intn(len(availableColors))]
			surface.DrawTextColor(8+idx, 9, fmt.Sprintf("%2d: %6d", idx+1, sc), fg, interfaces.ColorBlackDef, interfaces.ModeNormal)
		}
	}
}

func (t *Tetris) colorizeLine(surface interfaces.ISurface, line int, color interfaces.ColorDef) {
	for i := 0; i < t.board.w; i++ {
		t.drawBack(surface, i+boardXOffset, line+boardYOffset, color)
	}
}
