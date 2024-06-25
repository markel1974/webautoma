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
	"fmt"
	"github.com/markel1974/webautoma/src/shell/context/matrix"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"math/rand"
	"strconv"
	"time"
)

const (
	fgPlayText               = interfaces.ColorCyanDef
	bgPlayText               = interfaces.ColorRedDef
	playerSpriteBottomOffset = 1

	scoreX = 10
	scoreY = 1

	livesRightOffset = 0

	rwdSm = 10
	rwdMd = 20
	rwdLg = 30

	nonIndex = -50
)

type GameState uint8

const (
	MenuState GameState = iota
	PlayState
	HighscoresState
)

type Invaders struct {
	highscores []*HighScore
	state      GameState
	fc         uint8
	hmi        int
	w          int
	h          int

	startX int
	startY int
	player *Player

	ufo            *matrix.Entity
	ufoTimer       *time.Time
	aliens         *Aliens
	alienMoveEvery uint8

	fragments  *Fragments
	barricade  *Barricade
	lvl        int
	freeze     int
	freezeText string

	stars *Stars
}

func NewGame(w int, h int) *Invaders {
	g := &Invaders{
		h:     h,
		w:     w,
		stars: NewStars(),
	}
	g.init()
	return g
}

func (g *Invaders) init() {
	g.highscores = make([]*HighScore, 0)
	g.fc = 1

	g.barricade = NewBarricade()
	g.fragments = NewFragments()
	g.aliens = NewAliens()
}

func (g *Invaders) SetSize(w int, h int) {
	g.w = w
	g.h = h
	g.init()
}

func (g *Invaders) GetSize() (int, int) {
	return g.w, g.h
}

func (g *Invaders) drawMenu(surface interfaces.ISurface) {
	x := g.w/2 - logoLineLength/2
	y := logoY

	for _, line := range logoLines {
		drawColor(surface, x, y, fgMenu, bgMenu, line)
		y++
	}

	length := 0
	i := 0
	for _, v := range menuItems {
		length += len(v)
		if i+1 != len(menuItems) {
			length += menuPad
		}
		i++
	}

	x = g.w/2 - length/2
	y += logoHeight + 5
	for i := FirstMenuItem; i < NumMenuItems; i++ {
		v := menuItems[i]
		if i == g.hmi {
			drawColor(surface, x, y, fgMenuHighlight, bgMenuHighlight, v)
		} else {
			drawColor(surface, x, y, fgMenu, bgMenu, v)
		}
		x += len(v) + menuPad
	}

	g.stars.Draw(surface)
}

func (g *Invaders) drawFlash(surface interfaces.ISurface) {
	if g.freeze > 0 {
		drawColor(surface, g.w/2-len(g.freezeText)/2, g.h/2, fgPlayText, bgPlayText, g.freezeText)
		g.freeze--
	}
}

func (g *Invaders) drawPlay(surface interfaces.ISurface) {
	g.stars.Draw(surface)
	g.barricade.Draw(surface)
	if g.ufo != nil {
		g.ufo.Draw(surface)
		g.ufo.Next()
	}
	g.aliens.Draw(surface)
	g.player.Draw(surface)
	drawColor(surface, scoreX, scoreY, fgPlayText, bgPlayText, strconv.Itoa(g.player.score))
	drawColor(surface, g.w-livesRightOffset-len(g.player.GetLivesStr()), scoreY, fgPlayText, bgPlayText, g.player.GetLivesStr())
	g.fragments.Draw(surface)
	g.drawFlash(surface)
}

func (g *Invaders) Draw(surface interfaces.ISurface) {
	switch g.state {
	case MenuState:
		g.drawMenu(surface)
	case PlayState:
		g.drawPlay(surface)
	case HighscoresState:
		//g.DrawHighscores()
	}
}

func (g *Invaders) Update() {
	g.fc++
	switch g.state {
	case MenuState:
		g.updateMenu()
	case PlayState:
		g.updatePlay()
	case HighscoresState:
		//g.UpdateHighscores()
	}
}

func (g *Invaders) updateMenu() {
	g.stars.Update(g.w, g.h, g.fc)
}

func (g *Invaders) updatePlay() {
	g.stars.Update(g.w, g.h, g.fc)

	for b := range g.aliens.alienBullets.data {
		if g.aliens.alienBullets.data[b] != nil {
			g.aliens.alienBullets.data[b].AddToY(alienBulletSpeed)

			if int(g.aliens.alienBullets.data[b].GetY()) >= g.h {
				g.aliens.alienBullets.data[b] = nil
			} else {
				x, y := g.aliens.alienBullets.data[b].GetX(), g.aliens.alienBullets.data[b].GetY()

				if g.player.HasCollision(g.aliens.alienBullets.data[b].Entity) {
					if !g.player.DecLives() {
						g.gameOver()
						return
					}

					g.wipeBullets()

					g.freeze = 20
					g.freezeText = fmt.Sprintf("Level %d", g.lvl)
					return
				} else if g.barricade.Unset(int(x), int(y)) {
					g.aliens.alienBullets.data[b] = nil
				}
			}
		}
	}

	if g.player.bullet != nil {
		g.player.bullet.AddToY(g.player.bullet.Vy)
		if g.player.bullet.GetY() < 0 {
			g.player.bullet = nil
		} else {
			found, reward := g.aliens.DoBulletCollision(g.player.bullet.Entity)
			if !found {
				if g.ufo != nil {
					if g.ufo.HasCollision(g.player.bullet.Entity) {
						reward = ufoReward
						g.ufo = nil
						found = true
					}
				}
			}

			if found {
				g.explode(int(g.player.bullet.GetX()), int(g.player.bullet.GetY()))
				g.player.score += reward
				g.player.WipeBullet()
			} else {
				if g.barricade.Unset(int(g.player.bullet.GetX()), int(g.player.bullet.GetY())) {
					g.player.WipeBullet()
				}
			}
		}
	}

	if g.ufo != nil && g.fc%ufoMoveEvery == 0 {
		g.ufo.AddToX(1)
		if int(g.ufo.GetX()) > g.w {
			g.ufo = nil
			t := time.Now()
			t = t.Add(time.Duration(rand.Intn(20)+15) * time.Second)
			g.ufoTimer = &t
		}
	}

	if g.fc%g.alienMoveEvery == 0 {
		g.aliens.DoMove(g.w, g.h)

		if g.aliens.minY >= int(g.player.sprite.GetY()-g.player.sprite.GetHeight()) {
			g.gameOver()
			return
		}

		if g.aliens.count <= 0 && g.ufo == nil {
			g.lvl += 1
			if g.alienMoveEvery != 2 {
				g.alienMoveEvery--
			}
			g.beginNextLevel()
			g.wipeBullets()
		}
	}

	g.fragments.Update()

	if g.ufoTimer != nil {
		if (*g.ufoTimer).Before(time.Now()) {
			g.ufo = newUfo()
			g.ufoTimer = nil
		}
	}
}

func (g *Invaders) HandleKey(k rune) {
	switch g.state {
	case MenuState:
		g.handleKeyMenu(k)
	case PlayState:
		g.handleKeyPlay(k)
	case HighscoresState:
		//g.HandleKeyHighscores(k)
	}
}

func (g *Invaders) handleKeyMenu(k rune) {
	switch k {
	case 'a':
		g.hmi = (g.hmi - 1 + NumMenuItems) % NumMenuItems
	case 'd':
		g.hmi = (g.hmi + 1) % NumMenuItems
	case ' ':
		switch g.hmi {
		case Highscores:
			//g.GoHighscores()
		case Howto:
			//g.GoHowto()
		case Play:
			g.SetPlayState()
		}
	}
}

func (g *Invaders) handleKeyPlay(k rune) {
	switch k {
	case 'd':
		g.player.MoveRight(g.w)
	case 'a':
		g.player.MoveLeft(0)
	case ' ':
		g.player.NewBullet()
	}
}

func (g *Invaders) SetMenuState() {
	g.state = MenuState
	g.hmi = FirstMenuItem
}

func (g *Invaders) SetPlayState() {
	g.state = PlayState
	g.initPlay()
	g.beginNextLevel()
}

func (g *Invaders) beginNextLevel() {
	g.aliens.Create(alienStartX, alienStartY, 0)
}

func (g *Invaders) explode(x, y int) {
	g.fragments.AddFragment(x, y)
}

func (g *Invaders) wipeBullets() {
	g.player.WipeBullet()
	g.aliens.alienBullets.Setup(g.aliens.columns * g.aliens.rows / 10)
}

func (g *Invaders) initPlay() {
	g.fc = 1
	g.lvl = 1
	g.freeze = 0
	g.freezeText = ""

	g.player = NewPlayer(0, 0)

	g.startX = g.w/2 - int(g.player.sprite.GetWidth()/2)
	g.startY = g.h - playerSpriteBottomOffset - int(g.player.sprite.GetHeight())

	g.player.sprite.MoveTo(float64(g.startX), float64(g.startY))

	g.alienMoveEvery = uint8(alienMoveEvery)

	g.aliens.Setup(g.w, g.h)

	g.fragments.Setup()

	g.barricade.Setup(g.w, g.h)

	t := time.Now()
	t = t.Add(time.Duration(rand.Intn(20)+15) * time.Second)
	g.ufoTimer = &t
}

func (g *Invaders) gameOver() {
	g.freeze = 200
	g.freezeText = "GAME OVER"
	g.SetMenuState()
}
