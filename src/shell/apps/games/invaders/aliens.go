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
	"math/rand"
)

const (
	alienStartX        = 10
	alienStartY        = 7
	alienBulletSpeed   = 1
	alienShootValMax   = 100
	alienPadVertical   = 1
	alienPadHorizontal = 3
	alienMoveEvery     = 4
)

const (
	ufoMoveEvery = 3
	ufoReward    = 100
)

var (
	ufoSprite = []string{
		`
 YYYYYY
GMMMMMMM
 BB  BB`,
		`
 YYYYYY
MMGMMMMM
 BB  BB`,
		`
 YYYYYY
MMMMGMMM
 BB  BB`,
		`
 YYYYYY
MMMMMMGM
 BB  BB`,
		`
 YYYYYY
MMMMMMMM
 BB  BB`,
		`
 YYYYYY
MMMMMMMM
 BB  BB`,
		`
 YYYYYY
MMMMMMMM
 BB  BB`,
	}
)

var (
	smAlienSprite = []string{
		`
  YMMM
 GYYYYY
   CC
`,
		`
  MYMM
 YGYYYY
  C  C
`,
		`
  MMYM
 YYGYYY
 C    C
`,
		`
  MMMY
 YYYGYY
  C  C
`,
		`
  MMMM
 YYYYGY
   CC
`,
		`
  MMMM
 YYYYYG
  C  C
`,
		`
  MMMM
 YYYYYY
 C    C
`,
		`
  MMMM
 YYYYYY
  C  C
`,
	}

	mdAlienSprite = []string{
		`
G  BB  G
 RYYYYY
  MMMM
  B  B
`,
		`
   BB
GYRYYYYG
  MMMM 
  B  B
`,
		`
   BB
 YYRYYY
G MMMM G
  B  B
`,
		`
   BB
GYYYRYYG
  MMMM 
  B  B
`,
		`
G  BB  G
 YYYYRY
  MMMM 
  B  B
`,
		`
   BB  
GYYYYYRG
  MMMM 
  B  B
`,
		`
   BB  
 YYYYYY
G MMMM G
  B  B
`,
		`
   BB  
GYYYYYYG
  MMMM 
  B  B
`,
	}

	lgAlienSprite = []string{
		`
  CWWW
 GYYYYY
 GYYYYY
   WW
`,
		`
  WCWW
 YYGYYY
 YYYGYY
  W  W
`,
		`
  WWCW
 YYYYGY
 YYYYYG
 W    W
`,
		`
  WWWC
 YYYYYY
 YYYYYY
  W  W
`,
		`
  WWWW
 YYYYYY
 YYYYYY
   WW
`,
		`
  WWWW
 YYYYYY
 YYYYYY
  W  W
`,
		`
  WWWW
 YYYYYY
 YYYYYY
 W    W
`,
		`
  WWWW
 YYYYYY
 YYYYYY
  W  W
`,
	}
)

type Alien struct {
	*matrix.Entity
	idx     int
	reward  int
	counter int
}

func NewAlien(idx int, x int, y int, sprite []string, reward int) *Alien {
	return &Alien{
		Entity:  matrix.NewEntity(-1, x, y, 1, 1, 0, sprite, -1),
		idx:     idx,
		reward:  reward,
		counter: 3,
	}
}

func newUfo() *matrix.Entity {
	u := matrix.NewEntity(-1, 0, 0, 1, 0.5, 0.5, ufoSprite, -1)
	u.MoveTo(0-u.GetWidth(), u.GetHeight())
	return u
}

type Aliens struct {
	container []*Alien

	columns int
	rows    int
	rowsSm  int
	rowsMd  int
	rowsLg  int
	alienV  [2]int
	altered bool

	rightMove    [2]int
	leftMove     [2]int
	downMove     [2]int
	alienBullets *Bullets

	count int
	minY  int
	minX  int

	w    int
	h    int
	tree *matrix.AABBTree
}

func NewAliens() *Aliens {
	a := &Aliens{}

	a.rightMove = [2]int{1, 0}
	a.leftMove = [2]int{-1, 0}
	a.downMove = [2]int{0, 1}

	a.rowsSm = 2
	a.rowsMd = 2
	a.rowsLg = 1
	a.columns = 1
	a.rows = a.rowsSm + a.rowsMd + a.rowsLg
	a.altered = false
	a.alienV = a.rightMove
	a.alienBullets = NewBullets()

	a.count = 0

	return a
}

func (a *Aliens) Setup(w int, h int) {
	var i = 0
	var maxAlienWidth = 8
	var maxAlienHeight = 4

	a.columns = (int(float64(w) * 0.7)) / (maxAlienWidth + alienPadHorizontal)

	i = barricadeYPos(h) / (maxAlienHeight + alienPadVertical)

	switch i {
	case 0:
		a.rowsSm = 0
		a.rowsMd = 0
		a.rowsLg = 1
	case 1:
		a.rowsSm = 0
		a.rowsMd = 0
		a.rowsLg = 1
	case 2:
		a.rowsSm = 0
		a.rowsMd = 0
		a.rowsLg = 1
	case 3:
		a.rowsSm = 0
		a.rowsMd = 1
		a.rowsLg = 1
	case 4:
		a.rowsSm = 1
		a.rowsMd = 1
		a.rowsLg = 1
	case 5:
		a.rowsSm = 1
		a.rowsMd = 2
		a.rowsLg = 1
	default:
		a.rowsSm = 2
		a.rowsMd = 2
		a.rowsLg = 1
	}

	a.rows = a.rowsSm + a.rowsMd + a.rowsLg

	a.container = make([]*Alien, a.columns*a.rows)

	a.tree = matrix.NewAABBTree(uint(a.columns * a.rows))

	a.alienBullets.Setup(a.columns * a.rows / 10)

	a.alienV = a.rightMove
}

func (a *Aliens) Create(x int, y int, offset int) {
	y, offset = a.MakeAliens(a.columns, x, y, a.rowsLg, a.columns, rwdLg, offset, lgAlienSprite)
	y, offset = a.MakeAliens(a.columns, x, y, a.rowsMd, a.columns, rwdMd, offset, mdAlienSprite)
	a.MakeAliens(a.columns, x, y, a.rowsSm, a.columns, rwdSm, offset, smAlienSprite)
}

func (a *Aliens) Draw(surface interfaces.ISurface) {
	for _, k := range a.container {
		if k != nil {
			k.Draw(surface)
			k.Next()
		}
	}
	a.alienBullets.Draw(surface)
}

func (a *Aliens) MakeAliens(aliensHorizontal int, x int, y int, rows int, cols int, reward int, offset int, sprite []string) (int, int) {
	a.count = 0
	startX := x
	for i := 0; i < rows; i++ {
		var maxH = 0
		for j := 0; j < cols; j++ {
			idx := offset + i*aliensHorizontal + j
			if len(a.container) > 0 {
				if idx >= 0 && idx <= len(a.container) {
					alien := NewAlien(idx, x, y, sprite, reward)
					a.container[idx] = alien
					a.tree.InsertObject(alien)
					a.count++
					x += int(a.container[idx].GetWidth() + alienPadHorizontal)
					if int(a.container[idx].GetHeight()) > maxH {
						maxH = int(a.container[idx].GetHeight())
					}
				}
			}
		}
		x = startX
		y += maxH + alienPadVertical
	}
	return y, offset + rows*cols
}

func (a *Aliens) DoBulletCollision(bullet *matrix.Entity) (bool, int) {
	var reward = 0
	var found = false

	q := a.tree.QueryOverlaps(bullet)
	if len(q) > 0 {
		found = true
		for _, z := range q {
			var alien = z.(*Alien)
			alien.AddTo(bullet.Vx, bullet.Vy)
			bullet.DoPhysics(alien.Entity)
			alien.counter--
			if alien.counter <= 0 {
				reward = alien.reward
				a.tree.RemoveObject(alien)
				if alien.idx >= 0 && alien.idx < len(a.container) {
					a.container[alien.idx] = nil
				}
				a.count--
			}
		}
	}

	return found, reward
}

func (a *Aliens) DoMove(w int, _ int) {
	var downFlag = false
	var alienCount = 0

	a.minX = 999999999
	a.minY = 999999999

	for _, alien := range a.container {
		if alien != nil {
			alienCount++
			if a.altered {
				alien.Vx = float64(a.alienV[0])
				alien.Vy = float64(a.alienV[1])
			}

			alien.AddTo(alien.Vx, alien.Vy)

			a.tree.RemoveObject(alien)
			a.tree.InsertObject(alien)

			q := a.tree.QueryOverlaps(alien)
			if len(q) > 0 {
				for _, z := range q {
					var target = z.(*Alien)
					alien.DoPhysics(target.Entity)
				}
			}

			alienX := int(alien.GetX())
			alienY := int(alien.GetY())
			if alienX <= 0 || alienX+int(alien.GetWidth()) >= w {
				downFlag = true
				if alienX < a.minX {
					a.minX = alienX
				}
			}

			if alienY < a.minY {
				a.minY = alienY
			}

			if rand.Intn(alienShootValMax) == 6 {
				a.alienBullets.Update(int(alien.GetX())+int(alien.GetWidth())/2, int(alien.GetY()))
			}
		}
	}

	a.altered = false

	switch {
	case a.alienV == a.downMove:
		if a.minX == 0 {
			a.alienV = a.rightMove
		} else {
			a.alienV = a.leftMove
		}
		a.altered = true
	case downFlag:
		a.alienV = a.downMove
		a.altered = true
	}
}
