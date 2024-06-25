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

package matrix

type Rect struct {
	point  Point
	size   Size
	center Point
	z      float64
	aabb   *AABB
}

func NewRect(x float64, y float64, w float64, h float64, z float64) Rect {
	r := Rect{
		aabb:  &AABB{},
		point: NewPointFloat(x, y),
		size:  NewSize(w, h),
		z:     z,
	}
	r.rebuild()

	return r
}

func (r *Rect) rebuild() {
	r.center.x = r.point.x + (r.size.w / 2)
	r.center.y = r.point.y + (r.size.h / 2)

	r.aabb.minX = r.point.x
	r.aabb.maxX = r.point.x + r.size.w
	r.aabb.minY = r.point.y
	r.aabb.maxY = r.point.y + r.size.h
	r.aabb.minZ = 0
	r.aabb.maxZ = r.z
	r.aabb.surfaceArea = r.aabb.calculateSurfaceArea()
}

func (r *Rect) SetSize(w float64, h float64) {
	r.size.w += w
	r.size.h += h
	r.rebuild()
}

func (r *Rect) AddTo(x float64, y float64) {
	r.point.x += x
	r.point.y += y
	r.rebuild()
}

func (r *Rect) AddToX(x float64) {
	r.point.x += x
	r.rebuild()
}

func (r *Rect) AddToY(y float64) {
	r.point.y += y
	r.rebuild()
}

func (r *Rect) MoveTo(x float64, y float64) {
	r.point.x = x
	r.point.y = y
	r.rebuild()
}

func (r *Rect) MoveToX(x float64) {
	r.point.x = x
	r.rebuild()
}

func (r *Rect) MoveToY(y float64) {
	r.point.y = y
	r.rebuild()
}

func (r *Rect) GetX() float64 {
	return r.point.x
}

func (r *Rect) GetY() float64 {
	return r.point.y
}

func (r *Rect) GetWidth() float64 {
	return r.size.w
}

func (r *Rect) GetHeight() float64 {
	return r.size.h
}

func (r *Rect) GetAABB() *AABB {
	return r.aabb
}
