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

package shell

import (
	"github.com/markel1974/webautoma/src/shell/adaptiveticker"
	"github.com/markel1974/webautoma/src/shell/authenticator"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/context"
	"github.com/markel1974/webautoma/src/shell/terminal"
	"golang.org/x/term"
	"os"
)

func Create(autoSave bool, template *cli.Command) error {
	prompt := "% "
	factory := terminal.NewEquipmentFactory()
	ticker := adaptiveticker.NewAdaptiveTicker()
	auth := authenticator.NewSimpleAuthenticator()
	_, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	reader := os.Stdin
	writer := os.Stdout
	ctx := context.NewContext(ticker, reader, writer, auth, factory, template, prompt, autoSave)
	//ctx.SetEnterKey(13)
	ctx.Setup("VT100", false)
	ctx.Exec(true)

	return nil
}
