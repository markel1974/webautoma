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

package terminal

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"github.com/markel1974/webautoma/src/shell/terminal/vt100"
	"io"
)

type EquipmentFactory struct {
}

func NewEquipmentFactory() *EquipmentFactory {
	return &EquipmentFactory{}
}

func (f *EquipmentFactory) Create(_ string, z io.Writer, debug bool) interfaces.ITerminal {
	return vt100.NewVt100(z, debug)
}
