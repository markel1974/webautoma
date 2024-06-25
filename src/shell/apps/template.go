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

package apps

import (
	"github.com/markel1974/webautoma/src/shell/apps/games"
	"github.com/markel1974/webautoma/src/shell/apps/history"
	"github.com/markel1974/webautoma/src/shell/apps/runtime"
	"github.com/markel1974/webautoma/src/shell/apps/stats"
	"github.com/markel1974/webautoma/src/shell/apps/tasks"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"io"
)

type template struct {
	ctx    interfaces.IContext
	writer io.Writer
}

func NewTemplate(ctx interfaces.IContext, writer io.Writer) *template {
	return &template{
		ctx:    ctx,
		writer: writer,
	}
}

func (t *template) Run(template *cli.Command) *cli.Command {
	root := t.createInternalCommands()
	if template != nil {
		//template childs only
		t.commandIterator(root, template)
	}

	return root
}

func (t *template) AddCommand(cmd *cli.Command, child *cli.Command) {
	if !cmd.CommonCommands {
		//cmd.AddCommand(t.createChangeDirectory(z))
		cmd.CommonCommands = true
	}
	_ = cmd.AddCommand(child)
}

func (t *template) CreateCommand() *cli.Command {
	cmd := cli.NewCommand()
	t.setupCommand(cmd)
	return cmd
}

func (t *template) createInternalCommands() *cli.Command {
	root := cli.NewCommand()

	t.AddCommand(root, tasks.Create(t))
	t.AddCommand(root, history.Create(t))
	t.AddCommand(root, stats.Create(t))
	t.AddCommand(root, runtime.Create(t))

	t.AddCommand(root, CreateExit(t))
	t.AddCommand(root, CreateActivate(t))
	t.AddCommand(root, CreateKill(t))
	t.AddCommand(root, CreateKillAll(t))
	t.AddCommand(root, CreatePs(t))
	t.AddCommand(root, CreateClear(t))
	t.AddCommand(root, CreateFg(t))
	t.AddCommand(root, games.Create(t))

	root.SetOut(t.writer)
	root.SetErr(t.writer)
	return root
}

func (t *template) cloneCommand(src *cli.Command) *cli.Command {
	dst := cli.NewCommand()

	dst.Use = src.Use
	dst.Aliases = src.Aliases
	dst.SuggestFor = src.SuggestFor
	dst.Short = src.Short
	dst.Long = src.Long
	dst.Example = src.Example
	dst.ValidArgs = src.ValidArgs
	dst.Args = src.Args
	dst.ArgAliases = src.ArgAliases
	dst.Deprecated = src.Deprecated
	dst.Hidden = src.Hidden
	dst.Annotations = src.Annotations
	dst.Version = src.Version
	dst.SilenceErrors = src.SilenceErrors
	dst.SilenceUsage = src.SilenceUsage
	dst.DisableFlagParsing = src.DisableFlagParsing
	dst.DisableAutoGenTag = src.DisableAutoGenTag
	dst.DisableFlagsInUseLine = src.DisableFlagsInUseLine
	dst.DisableSuggestions = src.DisableSuggestions
	dst.SuggestionsMinimumDistance = src.SuggestionsMinimumDistance
	dst.TraverseChildren = src.TraverseChildren
	dst.FParseErrWhitelist = src.FParseErrWhitelist
	dst.PersistentPreRun = src.PersistentPreRun
	dst.PersistentPreRunE = src.PersistentPreRunE
	dst.PreRun = src.PreRun
	dst.PreRunE = src.PreRunE
	dst.Activate = src.Activate
	dst.Background = src.Background
	dst.Run = src.Run
	dst.ReadEvent = src.ReadEvent
	dst.TimerEvent = src.TimerEvent
	dst.PaintEvent = src.PaintEvent
	dst.RunE = src.RunE
	dst.PostRun = src.PostRun
	dst.PostRunE = src.PostRunE
	dst.PersistentPostRun = src.PersistentPostRun
	dst.PersistentPostRunE = src.PersistentPostRunE
	dst.SetUsageFunc(src.GetUsageFunc())
	dst.SetUsageTemplate(src.GetUsageTemplate())
	dst.SetFlagErrorFunc(src.GetFlagErrorFunc())
	dst.SetHelpFunc(src.GetHelpFunc())
	dst.SetHelpCommand(dst.GetHelpCommand())
	dst.SetHelpTemplate(src.GetHelpTemplate())
	dst.SetVersionTemplate(src.GetVersionTemplate())
	if src.GetGlobalNormalizationFunc() != nil {
		dst.SetGlobalNormalizationFunc(src.GetGlobalNormalizationFunc())
	}

	t.setupCommand(dst)

	return dst
}

func (t *template) commandIterator(dst *cli.Command, src *cli.Command) {
	if src.HasSubCommands() {
		for _, srcChild := range src.Childs() {
			if srcChild != nil {
				dstChild := t.cloneCommand(srcChild)
				t.AddCommand(dst, dstChild)
				//dst.AddCommand(dstChild)
				t.commandIterator(dstChild, srcChild)
			}
		}
	}
}

func (t *template) setupCommand(cmd *cli.Command) {
	cmd.SetRootContext(t.ctx)
	cmd.SetOut(t.writer)
	cmd.SetErr(t.writer)
}
