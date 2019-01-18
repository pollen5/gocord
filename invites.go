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

func (s *Shard) FetchInvite(code string, withCounts bool) (i *Invite, err error) {
	endpoint := rest.Invite(code)
	body, err := json.Marshal(&struct {
		WithCounts bool `json:"with_counts"`
	}{withCounts})

	if err != nil {
		return
	}

	err = s.Rest.Do(http.MethodGet, endpoint, body, &i)
	return
}

func (s *Shard) DeleteInvite(code string) (i *Invite, err error) {
	endpoint := rest.Invite(code)
	err = s.Rest.Do(http.MethodDelete, endpoint, nil, &i)
	return
}
