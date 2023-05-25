package structs

import "time"

type Icon struct {
	Host       string `json:"host"`
	Port       uint16 `json:"port"`
	Data       string `json:"data"`
	ObtainedAt time.Time 	 `json:"obtained_at"`
	ExpiresAt  time.Time 	 `json:"expires_at"`
}
