package matrix

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"math"
	"math/rand"
)

type Entity struct {
	Rect
	Mass    float64
	Vx      float64
	Vy      float64
	sprites []*Sprite
	index   int
	Bounce  bool
	Gravity bool
}

func NewEntity(base rune, x int, y int, mass float64, vx float64, vy float64, raw []string, startIdx int) *Entity {
	var sprites []*Sprite
	var width float64 = 0
	var height float64 = 0
	var index = 0

	for _, s := range raw {
		sd := NewSprite(s, base)
		if sd.size.w > width {
			width = sd.size.w
		}
		if sd.size.h > height {
			height = sd.size.h
		}
		sprites = append(sprites, sd)
	}

	if startIdx > len(sprites) {
		startIdx = -1
	}

	if startIdx < 0 {
		index = rand.Intn(len(sprites))
	} else {
		index = startIdx
	}

	a := &Entity{
		Rect:    NewRect(float64(x), float64(y), width, height, 1),
		sprites: sprites,
		Mass:    mass,
		Vx:      vx,
		Vy:      vy,
		Bounce:  true,
		Gravity: true,
		index:   index,
	}

	return a
}

func (e *Entity) Draw(surface interfaces.ISurface) {
	if len(e.sprites) == 0 {
		return
	}
	if e.index >= len(e.sprites) {
		e.index = 0
	}

	for y, v := range e.sprites[e.index].attr {
		for x, h := range v {
			surface.DrawColor(int(e.point.y)+y, int(e.point.x)+x, h.cur, h.fg, h.bg, interfaces.ModeNormal)
		}
	}
}

func (e *Entity) Next() {
	e.index++
}

func (e *Entity) HasCollision(obj2 *Entity) bool {
	return e.rectIntersect(e.point.x, e.point.y, e.size.w, e.size.h, obj2.point.x, obj2.point.y, obj2.size.w, obj2.size.h)
}

func (e *Entity) rectIntersect(x1 float64, y1 float64, w1 float64, h1 float64, x2 float64, y2 float64, w2 float64, h2 float64) bool {
	if x2 > w1+x1 || x1 > w2+x2 || y2 > h1+y1 || y1 > h2+y2 {
		return false
	}
	return true
}

func (e *Entity) DoPhysics(obj2 *Entity) bool {
	away := false
	if e.Bounce {
		//var distance = e.getDistance(e.point.x, e.point.y, obj2.point.x, obj2.point.y)
		//var vecCollision = Point{ x: obj2.point.x - e.point.x, y: obj2.point.y - e.point.y }
		var distance = e.getDistance(e.center.x, e.center.y, obj2.center.x, obj2.center.y)
		var vecCollision = Point{x: obj2.center.x - e.center.x, y: obj2.center.y - e.center.y}

		var vecCollisionNorm = Point{x: vecCollision.x / distance, y: vecCollision.y / distance}
		var vRelativeVelocity = Point{x: e.Vx - obj2.Vx, y: e.Vy - obj2.Vy}
		var speed = vRelativeVelocity.x*vecCollisionNorm.x + vRelativeVelocity.y*vecCollisionNorm.y

		if speed < 0 {
			away = true
		} else {
			if e.Gravity {
				var impulse = 2 * speed / (e.Mass + obj2.Mass)
				e.Vx -= impulse * obj2.Mass * vecCollisionNorm.x
				e.Vy -= impulse * obj2.Mass * vecCollisionNorm.y
				obj2.Vx += impulse * e.Mass * vecCollisionNorm.x
				obj2.Vy += impulse * e.Mass * vecCollisionNorm.y
			} else {
				e.Vx -= speed * vecCollisionNorm.x
				e.Vy -= speed * vecCollisionNorm.y
				obj2.Vx += speed * vecCollisionNorm.x
				obj2.Vy += speed * vecCollisionNorm.y
			}
		}
	}
	return away
}

func (e *Entity) getDistance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	d := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	if d == 0 {
		d = 0.001
	}
	return d
}
