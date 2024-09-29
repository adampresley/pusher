package local

import (
	"os"
	"strings"
	"text/template"

	"github.com/adampresley/pusher/pkg/rendering"
)

type LocalCommand struct {
	Command            string
	CommandDescription string
	Debug              bool
	ServiceName        string
	Host               string
}

func (l LocalCommand) Parse() string {
	var (
		err error
		t   *template.Template
	)

	result := &strings.Builder{}

	if t, err = template.New("context").Parse(l.Command); err != nil {
		rendering.Error("Panic parsing local command. Below is some context:")
		rendering.BlankLine()

		rendering.Print("err: %s", err.Error())
		rendering.Print("command: %s", l.Command)
		rendering.Print("context: %+v", l)

		os.Exit(1)
	}

	if err = t.Execute(result, l); err != nil {
		rendering.Error("Panic executing local command template. Below is some context:")
		rendering.BlankLine()

		rendering.Print("err: %s", err.Error())
		rendering.Print("command: %s", l.Command)
		rendering.Print("context: %+v", l)

		os.Exit(1)
	}

	return result.String()
}
