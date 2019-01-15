package gocord

import (
	"testing"
)

func TestPermissions(t *testing.T) {
	t.Run("has permissions", func(t *testing.T) {
		// 1024 = read messages + view_channel
		if !HasPermissions(PermissionsAdministrators, 1024) {
			t.Fail()
		}

		if !HasPermissions(268446768, PermissionsManageGuild) {
			t.Fail()
		}
	})

	t.Run("add permissions", func(t *testing.T) {
		if AddPermissions(PermissionsManageGuild, PermissionsManageRoles) != 268435488 {
			t.Fail()
		}
	})

	t.Run("remove permissions", func(t *testing.T) {
		if RemovePermissions(268435490, PermissionsKickMembers) != 268435488 {
			t.Fail()
		}
	})
}
