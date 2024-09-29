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
package project

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	PusherProjectFileName string = "pusher.yaml"
)

type PusherProject struct {
	CertEmail      string
	Dependencies   []string
	Domain         string
	EnvFile        string
	Host           string
	LastDeployDate string
	Mounts         Mounts
	Port           int
	ServiceName    string
	Version        int
}

func ProjectFileExists() bool {
	_, err := os.Stat(PusherProjectFileName)
	return err == nil
}

func (p *PusherProject) Load() error {
	var (
		err error
		f   *os.File
	)

	if f, err = os.Open(PusherProjectFileName); err != nil {
		return fmt.Errorf("There was an error trying to open the project file '%s': %s", PusherProjectFileName, err.Error())
	}

	decoder := yaml.NewDecoder(f)

	if err = decoder.Decode(p); err != nil {
		return fmt.Errorf("There was a problem decoding the project file '%s': %s", PusherProjectFileName, err.Error())
	}

	return nil
}

func (p *PusherProject) Save() error {
	var (
		err error
		out []byte
		f   *os.File
	)

	if out, err = yaml.Marshal(&p); err != nil {
		return fmt.Errorf("There was an error when converting project settings to YAML: %s", err.Error())
	}

	if f, err = os.Create(PusherProjectFileName); err != nil {
		return fmt.Errorf("There was an error attempting to create the project file '%s': %s", PusherProjectFileName, err.Error())
	}

	defer f.Close()

	if _, err = f.Write(out); err != nil {
		return fmt.Errorf("There was an error saving the project file '%s': %s", PusherProjectFileName, err.Error())
	}

	return nil
}

func (p *PusherProject) UpdateVersionAndDate() error {
	p.Version++
	p.LastDeployDate = time.Now().Format(time.RFC3339)

	return p.Save()
}
