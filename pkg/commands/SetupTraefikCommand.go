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

var SetupTraefikCommand = sshutils.Step{
	Commands: []sshutils.Command{
		sshutils.NewCommand(
			`cd ~ && mkdir -p traefik/ssl-certs`,
			"Installing Traefik...",
		),
		sshutils.NewCommand(
			`cd ~/traefik && tee docker-compose.yml <<EOF
services:
  traefik:
    image: traefik:v3.1
    container_name: traefik
    command: --api-insecure=false --providers.docker
    restart: unless-stopped
    ports:
      - 80:80
      - 443:443
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ~/traefik/traefik.yml:/etc/traefik/traefik.yml
      - ~/traefik/ssl-certs:/ssl-certs/
    networks:
      - web
      - applications

networks:
  web:
    external: true
  applications:
    external: true
EOF`,
			"Installing Traefik...",
		),
		sshutils.NewCommand(
			`cd ~/traefik && tee traefik.yml <<EOF
global:
  checkNewVersion: true
  sendAnonymousUsage: false

api:
  dashboard: false
  # Set insecure to false for production!
  insecure: false

entryPoints:
  web:
    address: :80
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https

  websecure:
    address: :443
    http:
      tls:
        certResolver: default

certificatesResolvers:
  default:
    acme:
      email: {{.Email}}
      storage: /ssl-certs/acme.json
      httpChallenge:
        entryPoint: web

providers:
  docker:
    exposedByDefault: false
EOF`,
			"Installing Traefik...",
		),
		sshutils.NewCommand(
			`cd traefik && sudo docker compose up -d`,
			"Starting Traefik...",
		),
	},
	StartingMessage: "Setting up Traefik...",
	SuccessMessage:  "Traefik setup successfully.",
	ErrorMessage:    "There was a problem setting up Traefik: %s",
}
