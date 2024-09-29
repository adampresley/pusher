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
	PostgresVersion string = "15.2"
)

var SetupServicePostgresCommand = sshutils.Step{
	Commands: []sshutils.Command{
		sshutils.NewCommand(
			`cd ~ && mkdir -p services/postgres/data`,
			"Installing PostgreSQL...",
		),
		sshutils.NewCommand(
			`cd ~/services/postgres && tee docker-compose.yml <<EOF
services:
  postgres:
    image: postgres:`+PostgresVersion+`
    container_name: postgres
    restart: unless-stopped
    ports:
      - 127.0.0.1:5432:5432
    environment:{{range $key, $value := .Env}}
      {{$key}}: "{{$value}}"{{end}}
    volumes:
      - ~/services/postgres/data:/var/lib/postgresql/data
    networks:
      - applications

networks:
  applications:
    external: true
EOF`,
			"Installing PostgreSQL...",
		),
		sshutils.NewCommand(
			`cd services/postgres && sudo docker compose up -d`,
			"Starting PostgreSQL...",
		),
	},
	StartingMessage: "Setting up PostgreSQL...",
	SuccessMessage:  "PostgreSQL setup successfully.",
	ErrorMessage:    "There was a problem setting up PostgreSQL: %s",
}
