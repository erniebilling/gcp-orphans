package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
)

type CommandFlags struct {
	gcpCredFile string
	gcpCreds    string
}

func CreateRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "gcp-orphans",
		Short: "GCP orphan discovery and cleanup",
		Long:  `A CLI to discover and cleanup orphaned GCP resources`,
	}
	root.AddCommand(createDiscoverOrphanedFirewallsCommand())

	return root
}

func AddCommandFlags(cmd *cobra.Command, commandFlags *CommandFlags) {
	cmd.Flags().StringVarP(&commandFlags.gcpCredFile, "gcpcredsfile", "g", "", "GCP service account credential filename")
	cmd.MarkFlagRequired("gcpcredsfile")
}

func ProcessOptions(commandFlags *CommandFlags) error {
	if commandFlags.gcpCreds == "" {
		data, err := ioutil.ReadFile(commandFlags.gcpCredFile)
		if err != nil {
			return err
		}
		commandFlags.gcpCreds = string(data)
	}
	return nil
}
