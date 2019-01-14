// Package embeds provide a helper API and structs and definitions related to Discord Embeds
package embeds

import (
	"strconv"
	"strings"
)

var colors = map[string]int{
	"default":            0x000000,
	"white":              0xFFFFFF,
	"aqua":               0x1ABC9C,
	"green":              0x2ECC71,
	"blue":               0x3498DB,
	"purple":             0x9B59B6,
	"vivid_pink":         0xE91E63,
	"gold":               0xF1C40F,
	"orange":             0xE67E22,
	"red":                0xE74C3C,
	"grey":               0x95A5A6,
	"navy":               0x34495E,
	"dark_aqua":          0x11806A,
	"dark_green":         0x1F8B4C,
	"dark_blue":          0x206694,
	"dark_purple":        0x71368A,
	"dark_vivid_pink":    0xAD1457,
	"dark_gold":          0xC27C0E,
	"dark_orange":        0xA84300,
	"dark_red":           0x992D22,
	"dark_grey":          0x979C9F,
	"darker_grey":        0x7F8C8D,
	"light_grey":         0xBCC0C0,
	"dark_navy":          0x2C3E50,
	"blurple":            0x7289DA,
	"greyple":            0x99AAB5,
	"dark_but_not_black": 0x2C2F33,
	"not_quite_black":    0x23272A,
}

type Embed struct {
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Timestamp   string         `json:"timestamp"`
	Color       int            `json:"color"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EmbedAuthor struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func New() *Embed {
	return &Embed{
		Type: "rich",
	}
}

func (e *Embed) SetTitle(t string) *Embed {
	if len(t) > 256 {
		t = t[0:255]
	}
	e.Title = t

	return e
}

func (e *Embed) SetDescription(d string) *Embed {
	if len(d) > 2048 {
		d = d[0:2047]
	}
	e.Description = d

	return e
}

func (e *Embed) SetColor(color interface{}) *Embed {
	switch color.(type) {
	case int, int32, int64:
		e.Color = color.(int)
	case string:
		if _, ok := colors[color.(string)]; ok {
			e.Color = colors[color.(string)]
		} else {
			hex, err := strconv.ParseInt(strings.Replace(color.(string), "#", "", 1), 16, 64)
			if err != nil {
				e.Color = int(hex)
			}
		}
	}

	return e
}
