package types

import "gorm.io/gorm"

type Proxy struct {
	gorm.Model

	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Status      int    `json:"status"`

	Source      string `json:"source"`
}

type ProxySource struct {
	Name string
	Url string
}

const (
	ProxyStatusFoo int = iota
	ProxyStatusBar
)

const (
	ErrInvalidIp = "invalid ip"
	ErrInvalidPort = "invalid port"
)
