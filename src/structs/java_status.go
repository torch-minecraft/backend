package structs

import "time"

type RawJavaStatus struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description interface{} `json:"description"`
	Favicon        string      `json:"favicon"`
	ModInfo     struct {
		List []struct {
			ModID   string `json:"modid"`
			Version string `json:"version"`
		} `json:"modList"`
		Type string `json:"type"`
	} `json:"modinfo"`
	ForgeData struct {
		Channels []struct {
			Required bool   `json:"required"`
			Res      string `json:"res"`
			Version  string `json:"version"`
		} `json:"channels"`
		FMLNetworkVersion int `json:"fmlNetworkVersion"`
		Mods              []struct {
			ModID   string `json:"modId"`
			Version string `json:"version"`
		} `json:"mods"`
	} `json:"forgeData"`
}

type Version struct {
	Name     *ParsedText `json:"name"`
	Protocol int         `json:"protocol"`
}

type Player struct {
	ID   string     `json:"id"`
	Name ParsedText `json:"name"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Player `json:"sample"`
}

type ModInfo struct {
	Type    string `json:"type"`
	ModList []Mod  `json:"modList"`
}

type Mod struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type SrvRecord struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type JavaStatus struct {
	Host        string        `json:"host"`
	Port        uint16        `json:"port"`
	Version     Version       `json:"version"`
	Players     Players       `json:"players"`
	Description *ParsedText   `json:"description"`
	Icon        string        `json:"icon"`
	ModInfo     *ModInfo      `json:"mod_info"`
	SrvRecord   *SrvRecord    `json:"used_srv"`
	Latency     time.Duration `json:"latency"`
	ObtainedAt  time.Time     `json:"obtained_at"`
	ExpiresAt   time.Time     `json:"expires_at"`
}

type OfflineServer struct {
	Offline bool   `json:"offline"`
	Host    string `json:"host"`
	Port    uint16 `json:"port"`
}
