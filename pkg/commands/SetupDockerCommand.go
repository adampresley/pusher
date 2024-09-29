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

const (
	lazyDockerVersion string = "0.23.3"
)

var SetupDockerCommand = sshutils.Step{
	Commands: []sshutils.Command{
		sshutils.NewCommand(
			"sudo install -m 0755 -d /etc/apt/keyrings",
			"Setting up keyring...",
		),
		sshutils.NewCommand(
			"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
			"Setting up keyring...",
		),
		sshutils.NewCommand(
			"sudo chmod a+r /etc/apt/keyrings/docker.asc",
			"Setting up keyring...",
		),
		sshutils.NewCommand(
			`echo \
         "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
         $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
         sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`,
			"Setting up keyring...",
		),
		sshutils.NewCommand(
			"sudo apt update -y",
			"Updating package list...",
		),
		sshutils.NewCommand(
			"sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y",
			"Installing Docker...",
		),
		sshutils.NewCommand(
			"sudo usermod -aG docker {{.User}}",
			"Adding user to 'docker' group...",
		),
		sshutils.NewCommand(
			`sudo docker network create -d bridge applications || true && sudo docker network create -d bridge web || true`,
			"Setting up Docker network...",
		),
		sshutils.NewCommand(
			`cd ~ && wget https://github.com/jesseduffield/lazydocker/releases/download/v`+lazyDockerVersion+`/lazydocker_`+lazyDockerVersion+`_Linux_x86_64.tar.gz`,
			"Installing LazyDocker...",
		),
		sshutils.NewCommand(
			`cd ~ && mkdir -p ./lazydocker && tar xvf ./lazydocker_`+lazyDockerVersion+`_Linux_x86_64.tar.gz -C ./lazydocker && sudo ln -sf ~/lazydocker/lazydocker /usr/local/bin/lazydocker`,
			"Installing LazyDocker...",
		),
		sshutils.NewCommand(
			`rm ./lazydocker_`+lazyDockerVersion+`_Linux_x86_64.tar.gz`,
			"Cleaning up...",
		),
	},
	StartingMessage: "Setting up Docker...",
	SuccessMessage:  "Docker setup successfully.",
	ErrorMessage:    "There was a problem setting up Docker: %s",
}
