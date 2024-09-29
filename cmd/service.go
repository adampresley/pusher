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
package cmd

import (
	"fmt"
	"os"

	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/adampresley/pusher/pkg/project"
	"github.com/adampresley/pusher/pkg/rendering"
	"github.com/adampresley/pusher/pkg/services"
	"github.com/adampresley/pusher/pkg/sshutils"
	"github.com/melbahja/goph"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Install a service on your server",
	Long: `Choose from a pre-configured list of services to install on your server,
      such as Postgres.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err                 error
			selectedServiceName string
			selectedService     services.ServiceItem
			sshClient           *goph.Client
			contextInfo         contextinfo.ContextInfo
		)

		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			rendering.Print("Debug enabled.")
		}

		rendering.Header("Deploy a Service")

		/*
		 * First load the project
		 */
		proj := &project.PusherProject{}

		if err = proj.Load(); err != nil {
			rendering.Error("There was a problem loading your project configuration file: %s", err.Error())
			os.Exit(1)
		}

		/*
		 * Prompt the user to choose from a set of services to install
		 */
		serviceList := services.GetServiceListNames()

		selectedServiceName, _ = pterm.DefaultInteractiveSelect.
			WithOptions(serviceList).
			WithDefaultText("Select a service").
			Show()

		selectedService = services.GetServiceByName(selectedServiceName)
		rendering.Paragraph(selectedService.Description)
		rendering.BlankLine()

		confirmation, _ := pterm.DefaultInteractiveConfirm.Show("Are you sure you wish to install this service?")

		if !confirmation {
			rendering.Warning("User cancelled. Aborting.")
			os.Exit(0)
		}

		/*
		 * Get an SSH client and start deploying.
		 */
		spinner := rendering.Spinner(fmt.Sprintf("Getting SSH client for host '%s'", proj.Host))

		if sshClient, contextInfo, err = sshutils.GetClientFromProject(proj); err != nil {
			spinner.Fail(fmt.Sprintf("Unable to get SSH client for host '%s': %s", proj.Host, err.Error()))
			os.Exit(1)
		}

		defer sshClient.Close()
		spinner.Success("Connection established.")

		/*
		 * Collect needed information then run.
		 */
		selectedService.Collector(&contextInfo)

		if err = selectedService.Step.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	serviceCmd.Flags().BoolP("debug", "d", false, "Enable debug output")
	rootCmd.AddCommand(serviceCmd)
}
