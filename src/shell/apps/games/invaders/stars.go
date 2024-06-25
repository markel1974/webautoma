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
	starSymbol = '.'
	numStars   = 50
)

func NewStar(x int, y int, symbol rune, color string, speed float64) *matrix.Entity {
	s := matrix.NewEntity(symbol, x, y, 0.5, speed, speed, []string{color}, -1)
	return s
}

type Stars struct {
	data []*matrix.Entity
	tree *matrix.AABBTree
}

func NewStars() *Stars {
	s := &Stars{
		tree: matrix.NewAABBTree(numStars),
		data: make([]*matrix.Entity, 0, numStars),
	}
	return s
}

func (s *Stars) Update(w int, h int, fc uint8) {
	if len(s.data) != cap(s.data) && fc%3 == 0 {
		n := len(s.data)
		s.data = s.data[0 : n+1]
		var symbol, speed, color = s.randSymbol()
		star := NewStar(rand.Intn(w), 0, symbol, color, speed)
		s.data[n] = star
	}

	for i, obj1 := range s.data {
		obj1.AddTo(obj1.Vx, obj1.Vy)

		s.tree.RemoveObject(obj1)

		if obj1.GetX() < 0 || int(obj1.GetX()) > w || obj1.GetY() < 0 || int(obj1.GetY()) > h {
			var symbol, speed, color = s.randSymbol()
			star := NewStar(rand.Intn(w), rand.Intn(h), symbol, color, speed)
			s.data[i] = star
			s.tree.InsertObject(star)
		} else {
			s.tree.InsertObject(obj1)
		}

		k := s.tree.QueryOverlaps(obj1)
		if len(k) > 0 {
			obj1.DoPhysics(obj1)
			for _, z := range k {
				obj1.DoPhysics(z.(*matrix.Entity))
			}
		}
	}
}

func (s *Stars) Draw(surface interfaces.ISurface) {
	for _, s := range s.data {
		s.Draw(surface)
		//surface.DrawEntity(s)
	}
}

func (s *Stars) randColor() interfaces.ColorDef {
	z := rand.Intn(6)
	color := interfaces.ColorWhiteDef

	switch z {
	case 0:
		color = interfaces.ColorRedDef
	case 1:
		color = interfaces.ColorYellowDef
	case 2:
		color = interfaces.ColorBlueDef
	case 3:
		color = interfaces.ColorMagentaDef
	case 4:
		color = interfaces.ColorCyanDef
	case 5:
		color = interfaces.ColorWhiteDef
	}

	return color
}

func (s *Stars) randSymbol() (rune, float64, string) {
	var z = rand.Intn(9)
	var symbol = starSymbol
	var vx = -0.1
	var color = "W"

	switch z {
	case 1:
		symbol = '.'
		vx = -0.1
		color = "W"
	case 2:
		symbol = '.'
		vx = 0.1
		color = "W"
	case 3:
		symbol = '*'
		vx = -0.5
		color = "B"
	case 5:
		symbol = '*'
		vx = 0.5
		color = "B"
	case 7:
		symbol = '+'
		vx = -0.3
		color = "C"
	case 8:
		symbol = '+'
		vx = 0.3
		color = "C"
	}

	return symbol, vx, color
}
