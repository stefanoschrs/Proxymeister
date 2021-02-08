package types

import (
	"gorm.io/gorm"
	"time"
)

type Proxy struct {
	gorm.Model

	Ip   string
	Port int

	Status       int
	Latency      int64
	FailedChecks int

	Source string
}
type ProxyFE struct {
	Address   *string    `json:"address,omitempty"`
	Ip        *string    `json:"ip,omitempty"`
	Latency   *int64     `json:"latency,omitempty"`
	Port      *int       `json:"port,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type ProxySource struct {
	Name string
	Url  string
}

type ProxifiedGetOptions struct {
	TargetUrl string
	Tries     *uint
}

// Proxy Status
const (
	ProxyStatusFresh int = iota
	ProxyStatusActive
	ProxyStatusInactive
)

// Errors
const (
	ErrInvalidIp             = "invalid ip"
	ErrInvalidPort           = "invalid port"
	ErrEmptyResponse         = "empty repsponse"
	ErrValidationPassthrough = "proxy passthrough"
	ErrValidationBadSquid    = "squid misconfiguration"
)
