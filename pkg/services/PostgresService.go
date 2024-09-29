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
package services

import (
	"github.com/adampresley/pusher/pkg/commands"
	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/pterm/pterm"
)

var PostgresService = ServiceItem{
	ServiceName: "PostgreSQL",
	Description: `PostgreSQL is a powerful, open source object-relational
database system`,
	Step: &commands.SetupServicePostgresCommand,
	Collector: func(info *contextinfo.ContextInfo) {
		info.Env = map[string]string{}

		info.Env["POSTGRES_USER"], _ = pterm.DefaultInteractiveTextInput.
			WithDefaultText("root").
			Show("User name")

		info.Env["POSTGRES_PASSWORD"], _ = pterm.DefaultInteractiveTextInput.
			Show("Password")

		database, _ := pterm.DefaultInteractiveTextInput.
			Show("Database (optional)")

		if database != "" {
			info.Env["POSTGRES_DB"] = database
		}
	},
}
