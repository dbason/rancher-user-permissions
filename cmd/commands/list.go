package commands

import (
	"errors"
	"fmt"

	"github.com/dbason/rancher-user-permissions/pkg/enumerator"
	"github.com/spf13/cobra"
)

var (
	clusterId string
	projectId string
)

func BuildListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List roles provided by group memberships",
		RunE:  listPermissions,
	}
	command.Flags().StringVar(&clusterId, "clusterid", "", "cluster id to list roles for")
	command.Flags().StringVar(&projectId, "projectid", "", "project id to list roles for")
	return command
}

func listPermissions(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if len(args) != 1 {
		return errors.New("list accepts the userid as a single argument")
	}
	userExists, err := enumerator.UserExists(ctx, args[0])
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user ID not found")
	}
	userGroups, err := enumerator.GroupMemberships(ctx, args[0])
	if err != nil {
		return err
	}
	for _, userGroup := range userGroups {
		globalRoles, err := enumerator.GetGlobalRoles(ctx, userGroup.GroupPrincipalName)
		if err != nil {
			return err
		}
		if len(globalRoles) > 0 {
			fmt.Printf("Group %s provides the following global roles:\n", userGroup.DisplayName)
			fmt.Printf("%v\n", globalRoles)
		}
	}
	return nil
}
