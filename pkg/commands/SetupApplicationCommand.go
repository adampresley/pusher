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

var SetupApplicationCommand = sshutils.Step{
	Commands: []sshutils.Command{
		sshutils.NewCommand(
			`cd ~ && mkdir -p applications/{{.ServiceName}}`,
			"Preparing {{.ServiceName}}...",
		),
		sshutils.NewCommand(
			`cd ~/applications/{{.ServiceName}} && tee docker-compose.yml <<EOF
services:
  {{.ServiceName}}:
    image: {{.ServiceName}}:latest
    container_name: {{.ServiceName}}
    restart: unless-stopped
    ports:
      - 127.0.0.1:{{.Port}}:{{.Port}}
    env_file:
      - {{.EnvFile}}{{if len .Dependencies}}
    depends_on:{{range .Dependencies}}
      - {{.}}{{end}}{{end}}{{if len .Mounts}}
    volumes:{{range .Mounts}}
      - {{.}}{{end}}{{end}}
    networks:
      - applications
    labels:
      - traefik.enable=true
      - traefik.http.routers.{{.ServiceName}}.rule=Host("{{.Domain}}")
      - traefik.http.services.{{.ServiceName}}.loadbalancer.server.port={{.Port}}
      - traefik.http.routers.{{.ServiceName}}.tls=true
      - traefik.http.routers.{{.ServiceName}}.tls.certresolver=default
      - traefik.docker.network=applications

networks:
  applications:
    external: true
EOF`,
			"Preparing {{.ServiceName}}...",
		),
	},
	StartingMessage: "Setting up your application...",
	SuccessMessage:  "Application setup successfully.",
	ErrorMessage:    "There was a problem setting up your application: %s",
}
