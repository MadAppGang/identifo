package model

import "time"

// UserDevice is a user push device
type UserDevice struct {
	ID          string
	DeviceType  UserDeviceType
	UserAgent   string
	Name        string
	Token       string
	Location    string
	IP          string
	LatestLogin time.Time
	CreatedAt   time.Time
}

type UserDeviceType string

// write a documentation for the following code
const (
	UserDeviceTypeIOS     UserDeviceType = "ios"
	UserDeviceTypeAndroid UserDeviceType = "android"
	UserDeviceTypeWeb     UserDeviceType = "web"
	UserDeviceTypeOther   UserDeviceType = "other"
)
