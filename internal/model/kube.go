package model

type KubeBadges struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Badge string `json:"badge"`

	Key         string `json:"key"`
	Allowed     bool   `json:"allowed"`
	DisplayName string `json:"display_name"`
	AliasURL    string `json:"alias_url"`
}

type KubeBadgesConfig struct {
	BadgeBaseURL string `json:"badge_base_url"`
}
