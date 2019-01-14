package gocord

// For user/member related definitions and methods

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
}
