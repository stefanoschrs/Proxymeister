package types

import "gorm.io/gorm"

type Proxy struct {
	gorm.Model

	Ip   string `json:"ip"`
	Port int    `json:"port"`

	Status       int   `json:"status"`
	Latency      int64 `json:"latency"`
	FailedChecks int   `json:"failedChecks"`

	Source string `json:"source"`
}

type ProxySource struct {
	Name string
	Url  string
}

const (
	ProxyStatusFresh int = iota
	ProxyStatusActive
	ProxyStatusInactive
)

const (
	ErrInvalidIp             = "invalid ip"
	ErrInvalidPort           = "invalid port"
	ErrValidationPassthrough = "proxy passthrough"
	ErrValidationBadSquid    = "squid misconfiguration"
)
