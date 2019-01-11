package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/erniebilling/gcp-orphans/gcp"
	"github.com/spf13/cobra"
)

func DiscoverOrphanedFirewallsCreate(commandFlags *CommandFlags) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := ProcessOptions(commandFlags); err != nil {
			return err
		}
		gcpClient := gcp.GCPClient{}
		if err := json.Unmarshal([]byte(commandFlags.gcpCreds), &gcpClient.Creds); err != nil {
			return err
		}
		if err := gcpClient.Connect(gcpClient.Creds); err != nil {
			return err
		}
		firewallRules, err := gcpClient.GetFirewallRules()
		if err != nil {
			return err
		}
		vmInstances, err := gcpClient.GetVMInstances()
		if err != nil {
			return err
		}

		orphans := make(map[string]gcp.FirewallRule, len(*firewallRules))
		for _, rule := range *firewallRules {
			orphans[rule.Name] = rule
		}

		for _, rule := range *firewallRules {
			for _, instance := range *vmInstances {
				if len(rule.Tags) == 0 || !IsEmtpyIntersection(rule.Tags, instance.NetworkTags) {
					delete(orphans, rule.Name)
				}
			}
		}

		fmt.Printf("Found %v orphaned firewall rules", len(orphans))
		fmt.Println()
		for _, orphan := range orphans {
			fmt.Printf("%v: %v", orphan.Name, orphan.URL)
			fmt.Println()
		}
		return err
	}
}

func createDiscoverOrphanedFirewallsCommand() *cobra.Command {
	commandFlags := CommandFlags{}
	discoverCmd := &cobra.Command{
		Use:     "firewallrules",
		Example: "gcp-orphans firewallrules",
		Short:   "discover orphaned firewall rules",
		Long:    `discover orphaned INGRESS firewall rules - those that apply to no VM instances`,
		RunE:    DiscoverOrphanedFirewallsCreate(&commandFlags),
	}
	AddCommandFlags(discoverCmd, &commandFlags)

	return discoverCmd
}

func IsEmtpyIntersection(set1, set2 []string) bool {
	for _, val1 := range set1 {
		for _, val2 := range set2 {
			if val1 == val2 {
				return false
			}
		}
	}
	return true
}
