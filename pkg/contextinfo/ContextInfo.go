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
package contextinfo

import (
	"os"
	"strings"
	"text/template"

	"github.com/adampresley/pusher/pkg/rendering"
)

type ContextInfo struct {
	Dependencies []string
	Domain       string
	Email        string
	Env          map[string]string
	EnvFile      string
	HostName     string
	IdentityFile string
	Mounts       []string
	Port         string
	ServiceName  string
	User         string
}

func (c ContextInfo) ExpandCommand(command string) string {
	var (
		err error
		t   *template.Template
	)

	result := &strings.Builder{}

	if t, err = template.New("context").Parse(command); err != nil {
		rendering.Error("Panic parsing command. Below is some context:")
		rendering.BlankLine()

		rendering.Print("err: %s", err.Error())
		rendering.Print("command: %s", command)
		rendering.Print("context: %+v", c)

		os.Exit(1)
	}

	if err = t.Execute(result, c); err != nil {
		rendering.Error("Panic executing command template. Below is some context:")
		rendering.BlankLine()

		rendering.Print("err: %s", err.Error())
		rendering.Print("command: %s", command)
		rendering.Print("context: %+v", c)

		os.Exit(1)
	}

	return result.String()
}
