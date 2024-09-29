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
package parsing_test

import (
	"strings"
	"testing"

	"github.com/adampresley/pusher/pkg/parsing"
	"github.com/stretchr/testify/assert"
)

func TestGetSSHConfigHosts(t *testing.T) {
	t.Run("returns a slice of hosts when successful", func(t *testing.T) {
		configFile := `Host *
   ServerAliveInterval 300

Host testing
   HostName 1.2.3.4
   User bob
   IdentityFile ~/.ssh/id_rsa`

		want := []string{
			"testing",
		}

		reader := strings.NewReader(configFile)
		got, err := parsing.GetSSHConfigHosts(reader)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
