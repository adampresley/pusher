package local

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/adampresley/pusher/pkg/rendering"
)

func RunLocalCommand(cmd LocalCommand) {
	commandText := cmd.Parse()
	// cwd, _ := os.Getwd()

	spinner := rendering.Spinner(fmt.Sprintf("Running: %s...", cmd.CommandDescription))

	// commandRunner := exec.Command("sh", "-s", "-", cwd, cmd.SerivceName, cmd.Host)
	commandRunner := exec.Command("sh", "-s")
	commandRunner.Stdin = strings.NewReader(commandText)

	if cmd.Debug {
		rendering.Print("Running local command: %s", commandText)
	}

	err := commandRunner.Run()

	if err != nil {
		spinner.Fail("Error running '" + cmd.CommandDescription + "'")
		rendering.BlankLine()

		rendering.Print("error: %s", err.Error())
		rendering.Print("command: %+v", cmd)
		rendering.Print("parsed command text: %s", commandText)

		os.Exit(1)
	}

	spinner.Success("Finished '" + cmd.CommandDescription + "'")
}
