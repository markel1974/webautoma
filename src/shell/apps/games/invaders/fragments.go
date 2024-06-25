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
	fragmentLifetime       = 15
	fragmentNum            = 10
	fragmentMaxX           = 16
	fragmentMaxY           = 8
	fragmentOffsetX        = -(fragmentMaxX / 2)
	fragmentOffsetY        = -(fragmentMaxY / 2)
	fragmentLifetimeUpdate = 2
)

type Fragment struct {
	matrix.Point
	life     int
	sprites  []*matrix.Entity
	fragment string
}

func NewFragment(x int, y int, life int, fragments int) Fragment {
	f := Fragment{
		Point:   matrix.NewPointFloat(float64(x), float64(y)),
		life:    life,
		sprites: make([]*matrix.Entity, fragments),
	}
	return f
}

type Fragments struct {
	data []*Fragment
}

func NewFragments() *Fragments {
	return &Fragments{}
}

func (f *Fragments) Setup() {
	f.data = make([]*Fragment, 0)
}

func (f *Fragments) AddFragment(x int, y int) {
	fg := NewFragment(x, y, 0, fragmentNum)
	f.data = append(f.data, &fg)
}

func (f *Fragments) Update() {
	var fragmentsCopy = make([]*Fragment, 0)
	for i := range f.data {
		if f.data[i].life > fragmentLifetime {
			continue
		}
		if f.data[i].life%fragmentLifetimeUpdate == 0 {
			for j := 0; j < fragmentNum; j++ {
				x := rand.Intn(fragmentMaxX) + int(f.data[i].GetX()) + fragmentOffsetX
				y := rand.Intn(fragmentMaxY) + int(f.data[i].GetY()) + fragmentOffsetY
				f.data[i].sprites[j] = matrix.NewEntity('+', x, y, 0.2, 0.5, 0.5, []string{f.randSprite()}, 0)
			}
		}

		f.data[i].life++
		fragmentsCopy = append(fragmentsCopy, f.data[i])
	}
	f.data = fragmentsCopy
	if len(f.data) > 0 {
		f.data = f.data[:]
	}
}

func (f *Fragments) Draw(surface interfaces.ISurface) {
	for i := range f.data {
		for _, sprite := range f.data[i].sprites {
			sprite.Draw(surface)
			sprite.Next()
		}
	}
}

func (f *Fragments) randSprite() string {
	var z = rand.Intn(10)
	var color string
	switch z {
	case 0:
		color = "G"
	case 1:
		color = "W"
	case 2:
		color = "YB"
	case 3:
		color = "B"
	case 4:
		color = "M"
	default:
		color = "Y"
	}

	return color
}
