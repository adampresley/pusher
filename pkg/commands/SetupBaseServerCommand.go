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
package commands

import "github.com/adampresley/pusher/pkg/sshutils"

var SetupBaseServerCommand = sshutils.Step{
	Commands: []sshutils.Command{
		sshutils.NewCommand(
			"sudo apt update -y",
			"Updating packages list...",
		),
		sshutils.NewCommand(
			"sudo apt upgrade -y",
			"Upgrading OS packages...",
		),
		sshutils.NewCommand(
			"sudo apt install ca-certificates curl wget htop neovim git -y",
			"Installing additional packages...",
		),
		sshutils.NewCommand(
			"cd ~ && mkdir -p /applications/ && mkdir -p /services/",
			"Setting up directories...",
		),
	},
	StartingMessage: "Updating OS and installing base software components...",
	SuccessMessage:  "Server update and software installed successfully.",
	ErrorMessage:    "There was a problem updating the OS and installing software components: %s",
}
