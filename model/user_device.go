package model

import "time"

// UserDevice is a user push device
type UserDevice struct {
	ID          string         `json:"id,omitempty"`
	DeviceType  UserDeviceType `json:"device_type,omitempty"`
	UserAgent   string         `json:"user_agent,omitempty"`
	Name        string         `json:"name,omitempty"`
	Token       string         `json:"token,omitempty"`
	Location    string         `json:"location,omitempty"`
	IP          string         `json:"ip,omitempty"`
	Detached    bool           `json:"detached,omitempty"`
	LatestLogin time.Time      `json:"latest_login,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	DetachedAt  time.Time      `json:"detached_at,omitempty"`
}

type UserDeviceType string

// write a documentation for the following code
const (
	UserDeviceTypeIOS     UserDeviceType = "ios"
	UserDeviceTypeAndroid UserDeviceType = "android"
	UserDeviceTypeWeb     UserDeviceType = "web"
	UserDeviceTypeOther   UserDeviceType = "other"
)
