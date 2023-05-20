package structs

type Icon struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	Data string `json:"data"`
}
