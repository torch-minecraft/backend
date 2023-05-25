package structs

import "time"

type Srv struct {
	Target     string    `json:"target"`
	Port       uint16    `json:"port"`
	ObtainedAt time.Time `json:"obtained_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}
