package model

import (
	"time"
)

// Invite is a representation of the invite model.
// Token field is required for proper working.
type Invite struct {
	ID          string    `json:"id"  bson:"_id"`
	Archived    bool      `json:"archived"  bson:"archived"`
	AppID       string    `json:"app_id"  bson:"app_id"`
	InviterID   string    `json:"inviter_id"  bson:"inviter_id"`
	InviterName string    `json:"inviter_name"  bson:"inviter_name"`
	Token       string    `json:"token"  bson:"token"`
	Email       string    `json:"email"  bson:"email"`
	Role        string    `json:"role"  bson:"role"`
	Tenant      string    `json:"tenant"  bson:"tenant"`
	TenantName  string    `json:"tenant_name"  bson:"tenant_name"`
	Group       string    `json:"group"  bson:"group"`
	CreatedBy   string    `json:"created_by"  bson:"created_by"`
	CreatedAt   time.Time `json:"created_at"  bson:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"  bson:"expires_at"`
}
