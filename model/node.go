package model

type Node struct {
	IP       string         `json:"ip"`
	Port     int            `json:"port"`
	State    string         `json:"state"`
	Address  string         `json:"address"`
	Metadata map[string]any `json:"extendInfo"`
}
