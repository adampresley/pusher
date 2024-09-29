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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adampresley/pusher/pkg/commands"
	"github.com/adampresley/pusher/pkg/contextinfo"
	"github.com/adampresley/pusher/pkg/local"
	"github.com/adampresley/pusher/pkg/project"
	"github.com/adampresley/pusher/pkg/rendering"
	"github.com/adampresley/pusher/pkg/sshutils"
	"github.com/melbahja/goph"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your application",
	Long: `Create a Docker image of your application from a Dockerfile,
deploy it to your server, and increment the deploy version.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err                error
			port               string
			sshClient          *goph.Client
			contextInfo        contextinfo.ContextInfo
			changeDependencies bool
			dependencies       []string
			changeMounts       bool
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
		 * Gather data
		 */
	enterservicename:
		proj.ServiceName, _ = pterm.DefaultInteractiveTextInput.
			WithDefaultValue(proj.ServiceName).
			Show("Enter app name (directory and Docker friendly)")

		if proj.ServiceName == "" || strings.Contains(proj.ServiceName, " ") {
			rendering.Error("App name must not be empty of contain spaces.")
			goto enterservicename
		}

	enterport:
		port, _ = pterm.DefaultInteractiveTextInput.
			WithDefaultValue(strconv.Itoa(proj.Port)).
			Show("enter the port number your app binds to")

		if proj.Port, err = strconv.Atoi(port); err != nil {
			rendering.Error("Port must be a valid integer.")
			goto enterport
		}

		if proj.Port == 80 || proj.Port == 443 || proj.Port == 8080 {
			rendering.Error("Ports 80, 443, and 8080 are taken.")
			goto enterport
		}

		proj.Domain, _ = pterm.DefaultInteractiveTextInput.
			WithDefaultValue(proj.Domain).
			Show("Enter the domain (URL) to your app")

		if proj.EnvFile == "" {
			proj.EnvFile = ".env"
		}

		proj.EnvFile, _ = pterm.DefaultInteractiveTextInput.
			WithDefaultValue(proj.EnvFile).
			Show("Enter an env file containing your app settings")

		/*
		 * If we've already stored dependencies, ask the user if they
		 * want to keep the same list, or make a new one.
		 */
		if len(proj.Dependencies) > 0 {
			changeDependencies, _ = pterm.DefaultInteractiveConfirm.
				WithDefaultValue(false).
				Show("You already have dependencies defined for this project. Would you like to change them?")
		}

		if changeDependencies || len(proj.Dependencies) <= 0 {
		addmoredependencies:
			newDependency, _ := pterm.DefaultInteractiveTextInput.
				Show("Enter the name of a dependency (blank to finish)")

			if newDependency != "" {
				dependencies = append(dependencies, newDependency)
				goto addmoredependencies
			}
		}

		proj.Dependencies = dependencies

		/*
		 * Manage mounts
		 */
		if len(proj.Mounts) > 0 {
			changeMounts, _ = pterm.DefaultInteractiveConfirm.
				WithDefaultValue(false).
				Show("You already have mounts defined for this project. Would you like to manage them?")
		}

		if changeMounts || len(proj.Mounts) <= 0 {
			mountActions := []string{
				"Add mount",
				"Remove mount",
				"Done",
			}

		keepmanagingmounts:
			mountOptions := proj.Mounts.ToStrings()

			selectedMountAction, _ := pterm.DefaultInteractiveSelect.
				WithOptions(mountActions).
				WithDefaultText("What would you like to do?").
				Show()

			switch selectedMountAction {
			case "Add mount":
				newMountLocal, _ := pterm.DefaultInteractiveTextInput.
					Show("Enter local mount path")

				newMountRemote, _ := pterm.DefaultInteractiveTextInput.
					Show("Enter remote mount path")

				newMount := project.Mount{
					Local:  newMountLocal,
					Remote: newMountRemote,
				}

				proj.Mounts = append(proj.Mounts, newMount)

				rendering.Success("Mount added.")
				goto keepmanagingmounts

			case "Remove mount":
				selectedMountString, _ := pterm.DefaultInteractiveSelect.
					WithOptions(mountOptions).
					Show()

				newMounts := project.Mounts{}

				for _, m := range proj.Mounts {
					if m.String() != selectedMountString {
						newMounts = append(newMounts, m)
					}
				}

				rendering.Success("Mount removed.")
				goto keepmanagingmounts

			default:
				// We are done. Just keep going
			}
		}

		/*
		 * SAVE!
		 */
		if err = proj.Save(); err != nil {
			rendering.Error("There was a problem saving your project settings: %s", err.Error())
			os.Exit(1)
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
		 * Setup the app on the server
		 */
		contextInfo.ServiceName = proj.ServiceName
		contextInfo.Port = port
		contextInfo.Domain = proj.Domain
		contextInfo.EnvFile = proj.EnvFile
		contextInfo.Dependencies = proj.Dependencies
		contextInfo.Mounts = proj.Mounts.ToStrings()

		if debug {
			rendering.Print("context: %+v", contextInfo)
		}

		cwd, _ := os.Getwd()
		baseEnvFileName := filepath.Base(proj.EnvFile)
		envFileName := filepath.Join(cwd, baseEnvFileName)

		if debug {
			rendering.Print("local env: %s (base '%s')", envFileName, baseEnvFileName)
		}

		/*
		 * Upload the env file and docker-compose
		 */
		if err = commands.SetupApplicationCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		/*
		 * Upload env file
		 */
		uploadCmd := exec.Command("scp", envFileName, fmt.Sprintf("%s:~/applications/%s/%s", proj.Host, proj.ServiceName, baseEnvFileName))
		if err = uploadCmd.Run(); err != nil {
			rendering.Error("Unable to upload env file '%s': %s", envFileName, err.Error())
			os.Exit(1)
		}

		/*
		 * Create any missing mount folders on the server.
		 */
		for _, m := range proj.Mounts {
			createResult, err := sshClient.Run("mkdir -p " + m.Local)

			if err != nil {
				rendering.Error("Unable to create local mount folder on the server: %s", err.Error())
				rendering.Print("Folder: %s", m.Local)

				os.Exit(1)
			}

			if debug {
				rendering.Print("response: %s", string(createResult))
			}
		}

		/*
		 * Build docker image and copy it to the server.
		 */
		buildDockerImageCmd := local.LocalCommand{
			Command:            local.BuildDockerImageCommand,
			CommandDescription: "Build Docker Image",
			Debug:              debug,
			ServiceName:        proj.ServiceName,
			Host:               proj.Host,
		}

		local.RunLocalCommand(buildDockerImageCmd)

		copyDockerImageCmd := local.LocalCommand{
			Command:            local.UploadDockerImageCommand,
			CommandDescription: "Upload Docker Image",
			Debug:              debug,
			ServiceName:        proj.ServiceName,
			Host:               proj.Host,
		}

		local.RunLocalCommand(copyDockerImageCmd)

		if err = commands.LoadDockerApplicationCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		if err = commands.StartApplicationCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		if err = commands.CleanupApplicationCommand.Run(sshClient, contextInfo, debug); err != nil {
			os.Exit(1)
		}

		/*
		 * Update the project file version and date
		 */
		if err = proj.UpdateVersionAndDate(); err != nil {
			rendering.Error("Your application was deployed, but there was a problem updating the local project file: %s", err)
			os.Exit(1)
		}

		rendering.Header("🚀 Version %d deployed!", proj.Version)
	},
}

func init() {
	deployCmd.Flags().BoolP("debug", "d", false, "Enable debug output")
	rootCmd.AddCommand(deployCmd)
}
