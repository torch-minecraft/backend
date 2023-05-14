package structs

import "time"

type BedrockStatus struct {
	ServerGUID int64         `json:"server_guid"`
	Version    Version       `json:"version"`
	Edition    string        `json:"edition"`
	MOTD       *ParsedText   `json:"motd"`
	Players    Players       `json:"players"`
	ServerID   string        `json:"server_id"`
	Gamemode   string        `json:"gamemode"`
	GamemodeId int           `json:"gamemode_id"`
	Port       uint16        `json:"port"`
	PortIPv4   *int          `json:"port_ipv4"`
	PortIPv6   *int          `json:"port_ipv6"`
	Host       string        `json:"host"`
	ObtainedAt time.Time     `json:"obtained_at"`
	ExpiresAt  time.Time     `json:"expires_at"`
	Latency    time.Duration `json:"latency"`
}
