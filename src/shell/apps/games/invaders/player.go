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
	"github.com/markel1974/webautoma/src/shell/context/matrix"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"strings"
)

const (
	initLives         = 5
	livesSprite       = `| `
	playerMoveSpeed   = 2
	playerBulletSpeed = -1
)

var playerSprite = []string{
	`
  YY
GcGGGG
GGGGGG
`,
}

type Player struct {
	sprite   *matrix.Entity
	bullet   *Bullet
	score    int
	lives    int
	livesStr string
}

func NewPlayer(x int, y int) *Player {
	player := &Player{
		sprite: matrix.NewEntity(-1, x, y, 1, 1.0, 1.0, playerSprite, 0),
		score:  0,
		lives:  initLives,
		bullet: nil,
	}
	player.syncLives()
	return player
}

func (player *Player) HasCollision(e *matrix.Entity) bool {
	return player.sprite.HasCollision(e)
}

func (player *Player) Draw(surface interfaces.ISurface) {
	player.sprite.Draw(surface)
	player.sprite.Next()
	if player.bullet != nil {
		player.bullet.Draw(surface)
		player.bullet.Next()
	}
}

func (player *Player) NewBullet() {
	if player.bullet == nil {
		player.bullet = NewBullet(int(player.sprite.GetX()+player.sprite.GetWidth()/2), int(player.sprite.GetY()), playerBulletSpeed)
	}
}

func (player *Player) WipeBullet() {
	player.bullet = nil
}

func (player *Player) MoveLeft(minW int) {
	player.sprite.AddToX(-playerMoveSpeed)
	if int(player.sprite.GetX()) < minW {
		player.sprite.MoveToX(float64(minW))
	}
}

func (player *Player) MoveRight(maxW int) {
	player.sprite.AddToX(playerMoveSpeed)
	if int(player.sprite.GetX()+player.sprite.GetWidth()) > maxW {
		player.sprite.MoveToX(float64(maxW) - player.sprite.GetWidth())
	}
}

func (player *Player) DecLives() bool {
	player.lives -= 1
	player.syncLives()
	if player.lives <= 0 {
		return false
	}
	return true
}

func (player *Player) GetLivesStr() string {
	return player.livesStr
}

func (player *Player) syncLives() {
	player.livesStr = strings.Repeat(livesSprite, player.lives)
}
