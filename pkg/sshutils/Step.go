/*
Copyright © 2024 Adam Presley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package sshutils

import (
	"fmt"

	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/adampresley/pusher/pkg/rendering"
	"github.com/melbahja/goph"
	"github.com/pterm/pterm"
)

type Step struct {
	Commands        []Command
	StartingMessage string
	SuccessMessage  string
	ErrorMessage    string
}

func (s *Step) Run(sshClient *goph.Client, info contextinfo.ContextInfo, debug bool) error {
	var (
		err error
	)

	spinner := rendering.Spinner(s.StartingMessage)

	if err = s.runCommands(sshClient, info, spinner, debug); err != nil {
		spinner.Fail(fmt.Sprintf("%s: %s", s.ErrorMessage, err))
		return err
	}

	spinner.Success(s.SuccessMessage)
	return nil
}

func (s *Step) runCommands(sshClient *goph.Client, info contextinfo.ContextInfo, spinner *pterm.SpinnerPrinter, debug bool) error {
	for _, cmd := range s.Commands {
		var (
			err error
			b   []byte
		)

		spinner.UpdateText(cmd.Message)

		if b, err = s.runCommand(sshClient, info, cmd, debug); err != nil {
			if debug {
				rendering.Warning("DEBUG INFORMATION:")
				rendering.Paragraph(string(b))
			}

			return fmt.Errorf("There was an error running the command '%s': %s", cmd, err.Error())
		}

		if debug {
			rendering.Paragraph("DEBUG: %s", string(b))
		}
	}

	return nil
}

func (s *Step) runCommand(sshClient *goph.Client, info contextinfo.ContextInfo, command Command, debug bool) ([]byte, error) {
	cmd := info.ExpandCommand(command.Command)

	if debug {
		rendering.Print("COMMAND: %s", cmd)
	}

	return sshClient.Run(cmd)
}
