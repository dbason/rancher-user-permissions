package enumerator

import (
	"context"

	"github.com/dbason/rancher-user-permissions/pkg/util"

	managementv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

type UserGroup struct {
	DisplayName        string
	GroupPrincipalName string
}

func UserExists(ctx context.Context, userId string) (bool, error) {
	user := &managementv3.User{}
	err := util.K8sClient.Get(ctx, types.NamespacedName{
		Name: userId,
	}, user)
	if errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func GroupMemberships(ctx context.Context, userId string) ([]UserGroup, error) {
	var userGroups []UserGroup
	userAttributes := &managementv3.UserAttribute{}
	err := util.K8sClient.Get(ctx, types.NamespacedName{
		Name: userId,
	}, userAttributes)
	if err != nil {
		return userGroups, err
	}
	for _, v := range userAttributes.GroupPrincipals {
		if len(v.Items) > 0 {
			for _, principal := range v.Items {
				if principal.MemberOf {
					userGroups = append(userGroups, UserGroup{
						DisplayName:        principal.DisplayName,
						GroupPrincipalName: principal.Name,
					})
				}
			}
		}
	}
	return userGroups, nil
}

func GetGlobalRoles(ctx context.Context, groupPrincipalName string) ([]string, error) {
	var globalRoles []string
	globalRoleBindingList := &managementv3.GlobalRoleBindingList{}
	err := util.K8sClient.List(ctx, globalRoleBindingList)
	if err != nil {
		return globalRoles, err
	}
	for _, globalRoleBinding := range globalRoleBindingList.Items {
		if globalRoleBinding.GroupPrincipalName == groupPrincipalName {
			globalRoles = append(globalRoles, globalRoleBinding.GlobalRoleName)
		}
	}
	return globalRoles, nil
}
