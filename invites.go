package gocord

import (
	"encoding/json"
	"net/http"

	"github.com/Soumil07/gocord/rest"
)

type Invite struct {
	Code          string   `json:"code"`
	Guild         *Guild   `json:"guild,omitempty"`
	Channel       *Channel `json:"channel,omitempty"`
	PresenceCount int      `json:"approximate_presence_count"`
	MemberCount   int      `json:"approximate_member_count"`
}

func (c *Cluster) FetchInvite(code string, withCounts bool) (i *Invite, err error) {
	endpoint := rest.Invite(code)
	body, err := json.Marshal(&struct {
		WithCounts bool `json:"with_counts"`
	}{withCounts})

	if err != nil {
		return
	}

	err = c.Rest.Do(http.MethodGet, endpoint, body, &i)
	return
}

func (c *Cluster) DeleteInvite(code string) (i *Invite, err error) {
	endpoint := rest.Invite(code)
	err = c.Rest.Do(http.MethodDelete, endpoint, nil, &i)
	return
}
