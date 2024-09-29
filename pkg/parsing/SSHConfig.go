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
package parsing

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
)

var (
	DefaultSSHConfigFile string = filepath.Join(os.Getenv("HOME"), ".ssh", "config")
)

/*
GetSSHConfigHosts returns a string slice of all hosts names
found in an SSH config file.
*/
func GetSSHConfigHosts(reader io.Reader) ([]string, error) {
	var (
		err    error
		config *ssh_config.Config
		result []string
	)

	if config, err = ssh_config.Decode(reader); err != nil {
		return result, fmt.Errorf("There was an problem decoding the SSH config file: %s", err.Error())
	}

	m := map[string]struct{}{}

	for _, h := range config.Hosts {
		s := strings.TrimSpace(h.String())

		if s != "" {
			for _, p := range h.Patterns {
				if p.String() != "*" {
					if _, found := m[p.String()]; !found {
						m[p.String()] = struct{}{}
					}
				}
			}
		}
	}

	for k := range m {
		result = append(result, k)
	}

	sort.Strings(result)
	return result, nil
}

/*
GetSSHHost returns a host entry from an SSH config file.
*/
func GetSSHHost(reader io.Reader, hostKey string) (*ssh_config.Host, error) {
	var (
		err    error
		config *ssh_config.Config
		result *ssh_config.Host
	)

	if config, err = ssh_config.Decode(reader); err != nil {
		return result, fmt.Errorf("There was an problem decoding the SSH config file: %s", err.Error())
	}

	for _, h := range config.Hosts {
		s := strings.TrimSpace(h.String())

		if s != "" {
			for _, p := range h.Patterns {
				if p.String() == hostKey {
					return h, nil
				}
			}
		}
	}

	return result, fmt.Errorf("The host '%s' was not found in your SSH config", hostKey)
}

/*
OpenSSHConfigFile returns a file handle to an SSH config file
*/
func OpenSSHConfigFile(fileName string) (fs.File, error) {
	return os.Open(fileName)
}
