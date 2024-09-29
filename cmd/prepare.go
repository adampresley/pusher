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
	"io/fs"
	"os"

	"github.com/adampresley/pusher/pkg/commands"
	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/adampresley/pusher/pkg/parsing"
	"github.com/adampresley/pusher/pkg/project"
	"github.com/adampresley/pusher/pkg/rendering"
	"github.com/adampresley/pusher/pkg/sshutils"
	"github.com/adampresley/pusher/pkg/validation"
	"github.com/melbahja/goph"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare a server for deployment",
	Long: `Prepares a server for deploying your applications to. This command
will install the necessary software, such as Traefik and Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err            error
			f              fs.File
			availableHosts []string
			host           string
			certEmail      string
			contextInfo    contextinfo.ContextInfo
			sshClient      *goph.Client
		)

		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			rendering.Print("Debug enabled.")
		}

		rendering.Banner()
		rendering.Print("Let's setup your server. To get started, let's get some information.")
		rendering.Paragraph(`Pusher requires you to have an SSH config file with hosts
         configured. Each host entry must have a name, a HostName, a User, and an IdentityFile.`)
		rendering.BlankLine()

		/*
		 * If the project already exists, prompt the user about overwriting.
		 */
		if project.ProjectFileExists() {
			rendering.Warning(
				`A project file (%s) already exists. 
If you overwrite it, any settings
and previous deploys will be lost.`,
				project.PusherProjectFileName,
			)

			overwrite, _ := pterm.DefaultInteractiveConfirm.
				WithDefaultText("Overwrite project file?").
				Show()

			if !overwrite {
				rendering.Error("Aborting.")
				os.Exit(0)
			}
		}

		/*
		 * Get available hosts, then prompt the user to choose one.
		 */
		if f, err = parsing.OpenSSHConfigFile(parsing.DefaultSSHConfigFile); err != nil {
			rendering.Error(
				"There was a problem opening the SSH config file:\n  file: %s\n  error: %s\n",
				parsing.DefaultSSHConfigFile,
				err.Error(),
			)
			os.Exit(1)
		}

		defer f.Close()

		if availableHosts, err = parsing.GetSSHConfigHosts(f); err != nil {
			rendering.Error(
				"Unable to parse your SSH config file:\n  file: %s\n  error: %s\n",
				parsing.DefaultSSHConfigFile,
				err.Error(),
			)
			os.Exit(1)
		}

		host, _ = pterm.DefaultInteractiveSelect.
			WithOptions(availableHosts).
			WithDefaultText("Select a host").
			Show()

			/*
			 * Get an email to user for certificate generation
			 */
	enteremail:
		certEmailInput := pterm.DefaultInteractiveTextInput
		certEmailInput.DefaultText = "Enter an email for LetsEncrypt SSL certs"
		certEmail, _ = certEmailInput.Show()

		if !validation.IsValidEmail(certEmail) {
			rendering.Error("A valid email address is required for setting up Traefik and LetsEncrypt.")
			goto enteremail
		}

		/*
		 * Save a project file with our settings
		 */
		proj := project.PusherProject{
			CertEmail: certEmail,
			Host:      host,
		}

		if err = proj.Save(); err != nil {
			rendering.Error("%s - Aborting.", err.Error())
			os.Exit(1)
		}

		rendering.Success("Project file '%s' created.", project.PusherProjectFileName)

		/*
		 * Start setting up the server
		 */
		rendering.BlankLine()
		rendering.Header("Let's go! 🚀")

		spinner := rendering.Spinner("Connecting to remote host...")

		if sshClient, contextInfo, err = sshutils.GetClient(host); err != nil {
			rendering.Error("%s - Aborting.", err.Error())
			os.Exit(1)
		}

		contextInfo.Email = certEmail
		defer sshClient.Close()
		spinner.Success("Logged in successfully.")

		/*
		 * Start running through setup steps
		 */
		if err = commands.SetupBaseServerCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		if err = commands.SetupDockerCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		if err = commands.SetupTraefikCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		/*
		 * Done!
		 */
		rendering.BlankLine()
		rendering.Header("🥂 Your server is now setup!")
	},
}

func init() {
	prepareCmd.Flags().BoolP("debug", "d", false, "Enable debug output")
	rootCmd.AddCommand(prepareCmd)
}
