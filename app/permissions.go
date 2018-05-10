// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

const permissionsExportBatchSize = 100

func (a *App) ResetPermissionsSystem() *model.AppError {
	// Purge all roles from the database.
	if result := <-a.Srv.Store.Role().PermanentDeleteAll(); result.Err != nil {
		return result.Err
	}

	// Remove the "System" table entry that marks the advanced permissions migration as done.
	if result := <-a.Srv.Store.System().PermanentDeleteByName(ADVANCED_PERMISSIONS_MIGRATION_KEY); result.Err != nil {
		return result.Err
	}

	// Now that the permissions system has been reset, re-run the migration to reinitialise it.
	a.DoAdvancedPermissionsMigration()

	return nil
}

func (a *App) ExportPermissions() *model.AppError {

	rolesIterator := a.NewRoleIterator(permissionsExportBatchSize)

	for rolesIterator.HasNext {

		rolesBatch, err := rolesIterator.Next()
		if err != nil {
			return err
		}

		for _, role := range rolesBatch.([]*model.Role) {
			if !role.BuiltIn {
				roleExport, _ := json.Marshal(&struct {
					ID          string   `json:"id"`
					Name        string   `json:"name"`
					DisplayName string   `json:"display_name"`
					Description string   `json:"description"`
					Permissions []string `json:"permissions"`
				}{
					ID:          role.Id,
					Name:        role.Name,
					DisplayName: role.DisplayName,
					Description: role.Description,
					Permissions: role.Permissions,
				})
				fmt.Printf("%v\n", string(roleExport))
			}
		}

	}

	schemesIterator := a.NewSchemeIterator(permissionsExportBatchSize)

	for schemesIterator.HasNext {

		schemeBatch, err := schemesIterator.Next()
		if err != nil {
			return err
		}

		for _, scheme := range schemeBatch.([]*model.Scheme) {
			schemeExport, _ := json.Marshal(&struct {
				Name         string `json:"name"`
				Description  string `json:"description"`
				Scope        string `json:"scope"`
				TeamAdmin    string `json:"default_team_admin_role"`
				TeamUser     string `json:"default_team_user_role"`
				ChannelAdmin string `json:"default_channel_admin_role"`
				ChannelUser  string `json:"default_channel_user_role"`
			}{
				Name:         scheme.Name,
				Description:  scheme.Description,
				Scope:        scheme.Scope,
				TeamAdmin:    scheme.DefaultTeamAdminRole,
				TeamUser:     scheme.DefaultTeamUserRole,
				ChannelAdmin: scheme.DefaultChannelUserRole,
				ChannelUser:  scheme.DefaultChannelUserRole,
			})
			fmt.Printf("%v\n", string(schemeExport))
		}

	}

	return nil
}
