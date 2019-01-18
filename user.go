package gocord

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Soumil07/gocord/rest"
)

// For user/member related definitions and methods

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
}

// Tag returns a user's Discord tag (username + discriminator)
func (u *User) Tag() string {
	return fmt.Sprintf("%s#%s", u.Username, u.Discriminator)
}

// AvatarURL returns the user's CDN avatar url with the specified format and size
func (u *User) AvatarURL(format string, size int) string {
	if format == "" {
		if strings.HasPrefix(u.Avatar, "a_") {
			format = "gif"
		} else {
			format = "png"
		}
	}

	if u.Avatar == "" {
		return u.DefaultAvatarURL(size)
	}

	return fmt.Sprintf("%s/avatars/%s/%s.%s?size=%d", CdnUrl, u.ID, u.Avatar, format, size)
}

// DefaultAvatarURL returns a user's default Discord avatar url
func (u *User) DefaultAvatarURL(size int) string {
	discrim, _ := strconv.Atoi(u.Discriminator)
	return fmt.Sprintf("%s/embed/avatars/%d.png?size=%d", CdnUrl, discrim%5, size)
}

// FetchUser fetches a user given an ID
func (s *Shard) FetchUser(ID string) (u *User, err error) {
	endpoint := rest.User(ID)

	err = s.Rest.Do(http.MethodGet, endpoint, nil, &u)
	if err != nil {
		return
	}

	return
}
