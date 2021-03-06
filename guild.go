package gocord

import (
	"encoding/json"
	"net/http"

	"github.com/Soumil07/gocord/rest"
)

// defines guild related structs and helper functions

// Guild represents a Discord guild. NOTE: some guilds are unavailable at the ready event, and most
// guilds miss important properties, that are added at subsequent GuildCreate events
type Guild struct {
	ID                          string `json:"id"` // the ID of the guild
	Name                        string `json:"name"`
	Icon                        string `json:"icon,omitempty"`
	Splash                      string `json:"splash,omitempty"`
	OwnerID                     string `json:"owner_id"`
	Permissions                 int    `json:"permissions,omitempty"` // bitfield permissions
	Region                      string `json:"region"`
	AFKChannelID                string `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int    `json:"afk_timeout"`
	EmbedEnabled                bool   `json:"embed_enabled,omitempty"`
	EmbedChannelID              string `json:"embed_channel_id"`
	VerificationLevel           int    `json:"verification_level"`
	DefaultMessageNotifications int    `json:"default_message_notifications"`
	ExplicitContentFilter       int    `json:"explicit_content_filter"`

	Roles           []Role   `json:"roles"`
	Emojis          []Emoji  `json:"emojis"`
	Features        []string `json:"features"`
	MFALevel        int      `json:"mfa_level"`
	ApplicationID   string   `json:"application_id,omitempty"`
	WidgetEnabled   bool     `json:"widget_enabled"`
	WidgetChannelID string   `json:"widget_channel_id,omitempty"`
	SystemChannelID string   `json:"system_channel_id,omitempty"`

	JoinedAt    string                `json:"joined_at,omitempty"`
	Large       bool                  `json:"large,omitempty"`
	Unavailable bool                  `json:"unavailable,omitempty"`
	MemberCount int                   `json:"member_count,omitempty"`
	Members     []Member              `json:"members,omitempty"`
	Channels    []Channel             `json:"channels,omitempty"`
	Presences   []GuildMemberPresence `json:"presences,omitempty"`
}

func (g *Guild) String() string {
	return g.Name
}

type Role struct {
}

type Emoji struct {
}

type Member struct {
}

type GuildMemberPresence struct {
}

func (c *Cluster) BanMember(guildID, userID, reason string, deleteMessageDays int) (err error) {
	endpoint := rest.GuildBanMember(guildID, userID)

	body, err := json.Marshal(&struct {
		DeleteMessageDays int    `json:"delete-message-days"`
		Reason            string `json:"reason"`
	}{deleteMessageDays, reason})
	if err != nil {
		return
	}

	err = c.Rest.Do(http.MethodPut, endpoint, body, nil)
	return
}

func (c *Cluster) UnbanMember(guildID, userID string) (err error) {
	endpoint := rest.GuildBanMember(guildID, userID)
	err = c.Rest.Do(http.MethodDelete, endpoint, nil, nil)
	return
}
