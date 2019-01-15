package gocord

// Contains structs, definitions and helper methods related to permission bit-fields

const (
	PermissionsCreateInvite = 1 << iota
	PermissionsKickMembers
	PermissionsBanMembers
	PermissionsAdministrators
	PermissionsManageChannels
	PermissionsManageGuild
	PermissionsAddReactions
	PermissionsViewAuditLog
	PermissionsPrioritySpeaker
	_
	PermissionsViewChannel
	PermissionsSendMessages
	PermissionsSendTTSMessages
	PermissionsManageMessages
	PermissionsEmbedLinks
	PermissionsAttachFiles
	PermissionsReadMessageHistory
	PermissionsMentionEveryone
	PermissionsUseExternalEmojis
	_
	PermissionsConnect
	PermissionsSpeak
	PermissionsMuteMembers
	PermissionsDeafenMembers
	PermissionsMoveMembers
	PermissionsUseVAD
	PermissionsChangeNickname
	PermissionsManageNicknames
	PermissionsManageRoles
	PermissionsManageWebhooks
	PermissionsManageEmojis
)

func AddPermissions(original, added int) int {
	return original | added
}

func RemovePermissions(original, removed int) int {
	return original & (^removed)
}

// HasPermissions check if the supplied permission bitfield has a permission
func HasPermissions(original, perms int) bool {
	// check admin overwrites
	if (original & PermissionsAdministrators) == original {
		return true
	}
	return (original & perms) == perms
}
