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
	"io/fs"
	"strconv"
	"strings"

	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/adampresley/pusher/pkg/parsing"
	"github.com/adampresley/pusher/pkg/project"
	"github.com/kevinburke/ssh_config"
	"github.com/melbahja/goph"
)

func getClient(hostKey string) (*goph.Client, contextinfo.ContextInfo, error) {
	var (
		err      error
		f        fs.File
		hostInfo *ssh_config.Host
		auth     goph.Auth
		result   *goph.Client

		sshInfo contextinfo.ContextInfo
	)

	if f, err = parsing.OpenSSHConfigFile(parsing.DefaultSSHConfigFile); err != nil {
		return result, sshInfo, err
	}

	if hostInfo, err = parsing.GetSSHHost(f, hostKey); err != nil {
		return result, sshInfo, err
	}

	sshInfo = contextinfo.ContextInfo{}

	for _, n := range hostInfo.Nodes {
		line := strings.TrimSpace(n.String())
		lowerLine := strings.ToLower(line)

		if !n.Pos().Invalid() && line != "" {
			if strings.HasPrefix(lowerLine, "hostname") {
				if sshInfo.HostName, err = getValue(line); err != nil {
					return result, sshInfo, err
				}
			}

			if strings.HasPrefix(lowerLine, "user") {
				if sshInfo.User, err = getValue(line); err != nil {
					return result, sshInfo, err
				}
			}

			if strings.HasPrefix(lowerLine, "identityfile") {
				if sshInfo.IdentityFile, err = getValue(line); err != nil {
					return result, sshInfo, err
				}

				sshInfo.IdentityFile = parsing.ExpandHomeDir(sshInfo.IdentityFile)
			}
		}
	}

	/*
	 * Validate
	 */
	if sshInfo.HostName == "" {
		return result, sshInfo, fmt.Errorf("'HostName' not found in your SSH config for '%s'", hostKey)
	}

	if sshInfo.User == "" {
		return result, sshInfo, fmt.Errorf("'User' not found in your SSH config for '%s'", hostKey)
	}

	if sshInfo.IdentityFile == "" {
		return result, sshInfo, fmt.Errorf("'IdentityFile' not found in your SSH config for '%s'", hostKey)
	}

	if auth, err = goph.Key(sshInfo.IdentityFile, ""); err != nil {
		return result, sshInfo, fmt.Errorf("There was an error parsing your SSH key '%s': %s", sshInfo.IdentityFile, err.Error())
	}

	if result, err = goph.New(sshInfo.User, sshInfo.HostName, auth); err != nil {
		return result, sshInfo, fmt.Errorf("There was a problem setting up an SSH client to '%s' (user '%s'): %s", sshInfo.HostName, sshInfo.User, err.Error())
	}

	return result, sshInfo, nil
}

func GetClient(hostKey string) (*goph.Client, contextinfo.ContextInfo, error) {
	return getClient(hostKey)
}

func GetClientFromProject(proj *project.PusherProject) (*goph.Client, contextinfo.ContextInfo, error) {
	client, info, err := getClient(proj.Host)
	info.Email = proj.CertEmail
	info.Port = strconv.Itoa(proj.Port)
	info.ServiceName = proj.ServiceName

	return client, info, err
}

func getValue(line string) (string, error) {
	split := strings.Split(line, " ")

	if len(split) != 2 {
		return "", fmt.Errorf("Invalid line in SSH config file: '%s'", line)
	}

	return split[1], nil
}
