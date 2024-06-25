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

package interfaces

type ISurface interface {
	GetSize() (int, int)
	Draw(rows int, column int, c rune)
	DrawColor(rows int, column int, c rune, fg ColorDef, bg ColorDef, mode ColorMode)
	DrawText(rows int, column int, c string)
	DrawTextColor(rows int, column int, c string, fg ColorDef, bg ColorDef, mode ColorMode)
	//DrawEntity(e * matrix.Entity)
	DrawSeries(data []float64, w int, h int, min float64, max float64)
}
