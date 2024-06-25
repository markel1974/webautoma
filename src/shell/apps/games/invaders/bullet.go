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
)

var (
	bulletSprite = []string{"R"}
)

type Bullet struct {
	*matrix.Entity
}

func NewBullet(x int, y int, vy int) *Bullet {
	b := &Bullet{Entity: matrix.NewEntity(-1, x, y, 2, 0, float64(vy), bulletSprite, 0)}
	return b
}

type Bullets struct {
	data []*Bullet
}

func NewBullets() *Bullets {
	return &Bullets{}
}

func (ab *Bullets) Setup(num int) {
	ab.data = make([]*Bullet, num)
}

func (ab *Bullets) Draw(surface interfaces.ISurface) {
	for _, k := range ab.data {
		if k != nil {
			k.Draw(surface)
			k.Next()
		}
	}
}

func (ab *Bullets) Update(x int, y int) {
	for j := range ab.data {
		if ab.data[j] == nil {
			ab.data[j] = NewBullet(x, y, alienBulletSpeed)
			break
		}
	}
}
