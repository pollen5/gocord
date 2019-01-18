package gocord

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"

	// used for image decoding
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Soumil07/gocord/rest"
)

func (c *Cluster) LeaveGuild(guildID string) (err error) {
	endpoint := rest.UserGuild("@me", guildID)

	err = c.Rest.Do(http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return
	}

	return
}

func (c *Cluster) ModifyUser(username, avatar string) (u *User, err error) {
	endpoint := rest.User("@me")
	body, err := json.Marshal(&struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}{username, avatar})

	if err != nil {
		return
	}

	err = c.Rest.Do(http.MethodPatch, endpoint, body, &u)
	return
}

func (c *Cluster) SetAvatar(avatar io.Reader) (u *User, err error) {
	_, ext, err := image.Decode(avatar)
	if err != nil {
		return
	}

	reader := bufio.NewReader(avatar)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	data := fmt.Sprintf("data:image/%s;base64;%s", strings.ToLower(ext), encoded)

	return c.ModifyUser("", data)
}

func (c *Cluster) SetUsername(username string) (*User, error) {
	return c.ModifyUser(username, "")
}
