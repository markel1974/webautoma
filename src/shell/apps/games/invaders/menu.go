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
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"strings"
)

const (
	menuPad         = 10
	fgMenu          = interfaces.ColorRedDef
	bgMenu          = interfaces.ColorBlackDef
	fgMenuHighlight = interfaces.ColorWhiteDef
	bgMenuHighlight = interfaces.ColorCyanDef
	logoY           = 2
	logo            = `
                          _                     _
                         (_)                   | |              
 ___ _ __   __ _  ___ ___ _ _ ____   ____ _  __| | ___ _ __ ___ 
/ __| '_ \ / _  |/ __/ _ \ | '_ \ \ / /  | |/ _  |/ _ \ '__/ __|
\__ \ |_) | (_| | (_|  __/ | | | \ V / (_| | (_| |  __/ |  \__ \
|___/ .__/ \__,_|\___\___|_|_| |_|\_/ \__,_|\__,_|\___|_|  |___/
| |
|_|
`
)

const (
	FirstMenuItem     = 0
	Play          int = iota - 1
	Highscores
	Howto
	NumMenuItems
)

var (
	//menuItems      = map[int]string{Play: "PLAY", Highscores: "HIGHSCORES", Howto: "HOWTO"}
	menuItems      = map[int]string{Play: "PLAY"}
	logoLines      = strings.Split(logo, "\n")
	logoLineLength = len(logoLines[0])
	logoHeight     = len(logoLines)
)
